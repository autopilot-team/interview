package core

import (
	"autopilot/backends/internal/types"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	texttemplate "text/template"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/unrolled/render"
	"github.com/wneessen/go-mail"
)

// Mailer is the interface that wraps the Mail methods
type Mailer interface {
	// BulkSend sends multiple emails using the same template but different data
	BulkSend(templateName string, messages []EmailMessage, opts *RenderOptions) error

	// Render renders both HTML and text versions of an email template
	Render(templateName string, data map[string]any, opts *RenderOptions) (html string, text string, err error)

	// Send sends an email using rendered templates
	Send(templateName string, msg EmailMessage, opts *RenderOptions) error

	// SetupPreviewRoutes configures preview routes if preview is enabled
	SetupPreviewRoutes(router chi.Router)

	// TemplateOptions returns the template options
	TemplateOptions() *MailTemplateOptions
}

// MailTemplateOptions contains configurable options for the mailer
type MailTemplateOptions struct {
	// Dir is the directory containing the email templates
	Dir string

	// ExtraFuncs is a list of extra template functions to add to the mailer
	ExtraFuncs []template.FuncMap

	// FS is the embedded filesystem containing the email templates
	FS FS

	// Layout is the layout template to use for rendering emails
	Layout string
}

// RenderOptions contains options for rendering a template
type RenderOptions struct {
	Layout string // Optional override for layout
	Locale string // Optional locale (defaults to "en")
}

// DefaultRenderOptions returns the default render options
func DefaultRenderOptions() *RenderOptions {
	return &RenderOptions{
		Locale: "en",
	}
}

// Mail handles email template rendering and sending
type Mail struct {
	i18nBundle   *I18nBundle
	logger       *slog.Logger
	preview      *MailerPreview
	smtpURL      string
	templateOpts *MailTemplateOptions
}

// MailerOptions contains options for creating a new mailer
type MailOptions struct {
	I18nBundle      *I18nBundle
	Logger          *slog.Logger
	Mode            types.Mode
	PreviewData     map[string]map[string]any
	TemplateOptions *MailTemplateOptions
	SMTPURL         string
}

// NewMail creates a new mail instance
func NewMail(opts MailOptions) (*Mail, error) {
	if opts.TemplateOptions == nil {
		return nil, fmt.Errorf("template options cannot be nil")
	}

	if opts.I18nBundle == nil {
		return nil, fmt.Errorf("i18n bundle cannot be nil")
	}

	if opts.Logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if opts.TemplateOptions.Layout == "" {
		return nil, fmt.Errorf("template layout is required")
	}

	var zeroFS embed.FS
	if opts.TemplateOptions.FS == zeroFS || opts.TemplateOptions.FS == nil {
		return nil, fmt.Errorf("template filesystem cannot be nil")
	}

	if opts.TemplateOptions.Dir == "" {
		return nil, fmt.Errorf("template directory cannot be empty")

	}

	if opts.SMTPURL == "" {
		return nil, fmt.Errorf("SMTP URL cannot be empty")
	}

	m := &Mail{
		i18nBundle:   opts.I18nBundle,
		logger:       opts.Logger,
		smtpURL:      opts.SMTPURL,
		templateOpts: opts.TemplateOptions,
	}

	// Initialize preview in debug mode
	if opts.Mode == types.DebugMode {
		m.preview = NewMailerPreview(m, opts.I18nBundle, opts.Logger)
		if opts.PreviewData != nil {
			m.preview.SetPreviewData(opts.PreviewData)
		}
	}

	return m, nil
}

// Preview returns the mailer preview instance if available
func (m *Mail) Preview() *MailerPreview {
	return m.preview
}

// SetupPreviewRoutes configures preview routes if preview is enabled
func (m *Mail) SetupPreviewRoutes(router chi.Router) {
	if m.preview != nil {
		m.preview.setupEmailPreviewRoutes(router)
	}
}

// TemplateOptions returns the template options
func (m *Mail) TemplateOptions() *MailTemplateOptions {
	return m.templateOpts
}

// Render renders both HTML and text versions of an email template
func (m *Mail) Render(templateName string, data map[string]any, opts *RenderOptions) (html string, text string, err error) {
	if data == nil {
		return "", "", fmt.Errorf("data map cannot be nil")
	}

	// Use default options if none provided
	options := DefaultRenderOptions()
	if opts != nil {
		options = opts
	}

	localizer := NewLocalizer(m.i18nBundle, options.Locale)

	// Create base template functions
	baseFuncs := template.FuncMap{
		"t": func(messageID string, args ...any) (string, error) {
			if len(data) > 0 {
				for k, v := range data {
					args = append(args, k, v)
				}
			}

			return localizer.T()(messageID, args...)
		},
	}

	// Combine base functions with any extra functions
	funcMaps := []template.FuncMap{baseFuncs}
	funcMaps = append(funcMaps, m.templateOpts.ExtraFuncs...)

	// Merge all function maps into a single map
	allFuncs := template.FuncMap{}
	for _, fm := range funcMaps {
		for name, fn := range fm {
			allFuncs[name] = fn
		}
	}

	// Use override layout if provided, otherwise use default
	layout := m.templateOpts.Layout
	if options.Layout != "" {
		layout = options.Layout
	}

	// Create renderer for HTML
	htmlRenderer := render.New(render.Options{
		Directory:  m.templateOpts.Dir,
		Extensions: []string{".html"},
		FileSystem: func() render.FileSystem {
			if embedFS, ok := m.templateOpts.FS.(embed.FS); ok {
				return &render.EmbedFileSystem{
					FS: embedFS,
				}
			}

			return &render.LocalFileSystem{}
		}(),
		Funcs:  []template.FuncMap{allFuncs},
		Layout: layout,
	})

	// Render HTML version
	var htmlBuf strings.Builder
	if err := htmlRenderer.HTML(&htmlBuf, http.StatusOK, templateName, data); err != nil {
		return "", "", err
	}

	// Render text version using text/template
	var textBuf strings.Builder
	textTemplatePath := filepath.Join(m.templateOpts.Dir, templateName+".txt")

	// Read the template content
	textContent, err := m.templateOpts.FS.ReadFile(textTemplatePath)
	if err != nil {
		return htmlBuf.String(), "", nil // Return empty text version if template doesn't exist
	}

	// Create text template with layout support
	textTmpl := texttemplate.New(templateName).Funcs(texttemplate.FuncMap(allFuncs))

	// Add the layout template if specified
	if layout != "" {
		layoutPath := filepath.Join(m.templateOpts.Dir, layout+".txt")
		layoutContent, err := m.templateOpts.FS.ReadFile(layoutPath)
		if err == nil {
			// Parse layout template first
			textTmpl, err = textTmpl.Parse(string(layoutContent))
			if err != nil {
				return "", "", err
			}
		}
	}

	// Parse the main template
	textTmpl, err = textTmpl.Parse(string(textContent))
	if err != nil {
		return "", "", err
	}

	// Execute text template
	if err := textTmpl.Execute(&textBuf, data); err != nil {
		return "", "", err
	}

	return htmlBuf.String(), textBuf.String(), nil
}

// EmailMessage represents an email to be sent
type EmailMessage struct {
	From         string
	To           []string
	Subject      string
	Data         map[string]any
	MessageID    string // Optional message ID
	EnvelopeFrom string // Optional envelope from address
	Attachments  []*mail.File
}

// createSMTPClient creates a new SMTP client using the configured SMTP URL
func (m *Mail) createSMTPClient() (*mail.Client, error) {
	// Parse SMTP URL
	u, err := url.Parse(m.smtpURL)
	if err != nil {
		return nil, err
	}

	// Extract credentials from URL
	password, _ := u.User.Password()
	username := u.User.Username()

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, err
	}

	tlsPolicy := mail.TLSMandatory
	if u.Scheme == "smtp" {
		tlsPolicy = mail.NoTLS
	}

	mailOpts := []mail.Option{
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTLSPortPolicy(tlsPolicy),
		mail.WithUsername(username),
		mail.WithPassword(password),
	}

	return mail.NewClient(u.Hostname(), mailOpts...)
}

// Send sends an email using rendered templates
func (m *Mail) Send(templateName string, msg EmailMessage, opts *RenderOptions) error {
	// Render email templates
	html, text, err := m.Render(templateName, msg.Data, opts)
	if err != nil {
		return err
	}

	// Create mail message
	message := mail.NewMsg()
	if msg.EnvelopeFrom != "" {
		if err := message.EnvelopeFrom(msg.EnvelopeFrom); err != nil {
			return err
		}
	}
	if err := message.From(msg.From); err != nil {
		return err
	}

	// Add all recipients
	for _, to := range msg.To {
		if err := message.To(to); err != nil {
			return err
		}
	}

	message.Subject(msg.Subject)
	message.SetBodyString(mail.TypeTextHTML, html)
	message.AddAlternativeString(mail.TypeTextPlain, text)

	// Set attachments if any
	if len(msg.Attachments) > 0 {
		message.SetAttachments(msg.Attachments)
	}

	// Set message ID if provided
	if msg.MessageID != "" {
		message.SetMessageIDWithValue(msg.MessageID)
	} else {
		message.SetMessageID()
	}

	// Set date and bulk headers
	message.SetDate()
	message.SetBulk()

	client, err := m.createSMTPClient()
	if err != nil {
		return err
	}

	// Send the email
	if err := client.DialAndSend(message); err != nil {
		return err
	}

	m.logger.Info("Successfully sent email",
		"template", templateName,
		"from", msg.From,
		"to", msg.To,
		"subject", msg.Subject,
		"attachments", len(msg.Attachments),
	)

	return nil
}

// BulkSend sends multiple emails using the same template but different data
func (m *Mail) BulkSend(templateName string, messages []EmailMessage, opts *RenderOptions) error {
	if len(messages) == 0 {
		return fmt.Errorf("no messages to send")
	}

	client, err := m.createSMTPClient()
	if err != nil {
		return err
	}

	// Prepare all messages
	var mailMsgs []*mail.Msg
	totalAttachments := 0
	for _, msg := range messages {
		// Render email templates for each message
		html, text, err := m.Render(templateName, msg.Data, opts)
		if err != nil {
			return err
		}

		// Create mail message
		message := mail.NewMsg()
		if msg.EnvelopeFrom != "" {
			if err := message.EnvelopeFrom(msg.EnvelopeFrom); err != nil {
				return err
			}
		}

		if err := message.From(msg.From); err != nil {
			return err
		}

		// Add all recipients
		for _, to := range msg.To {
			if err := message.To(to); err != nil {
				return err
			}
		}

		message.Subject(msg.Subject)
		message.SetBodyString(mail.TypeTextHTML, html)
		message.AddAlternativeString(mail.TypeTextPlain, text)

		// Set attachments if any
		if len(msg.Attachments) > 0 {
			message.SetAttachments(msg.Attachments)
		}
		totalAttachments += len(msg.Attachments)

		// Set message ID if provided
		if msg.MessageID != "" {
			message.SetMessageIDWithValue(msg.MessageID)
		} else {
			message.SetMessageID()
		}

		// Set date and bulk headers
		message.SetDate()
		message.SetBulk()

		mailMsgs = append(mailMsgs, message)
	}

	// Send all messages in bulk
	if err := client.DialAndSend(mailMsgs...); err != nil {
		return err
	}

	m.logger.Info("Successfully sent bulk emails",
		"template", templateName,
		"count", len(messages),
		"total_attachments", totalAttachments,
	)

	return nil
}

// SetPreview sets the preview handler for the mailer
func (m *Mail) SetPreview(preview *MailerPreview) {
	m.preview = preview
}

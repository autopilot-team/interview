package core

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"log/slog"

	"github.com/go-chi/chi/v5"
)

const (
	mailerPreviewPath = "/mailer/preview"
	defaultLocale     = "en"
)

// MailerPreview handles email template preview functionality
type MailerPreview struct {
	mailer     Mailer
	i18nBundle *I18nBundle
	logger     *slog.Logger
	data       map[string]map[string]interface{}
}

// NewMailerPreview creates a new MailerPreview instance
func NewMailerPreview(mailer Mailer, i18nBundle *I18nBundle, logger *slog.Logger) *MailerPreview {
	return &MailerPreview{
		mailer:     mailer,
		i18nBundle: i18nBundle,
		logger:     logger,
		data:       make(map[string]map[string]interface{}),
	}
}

// SetPreviewData sets the preview data for templates
func (mp *MailerPreview) SetPreviewData(data map[string]map[string]interface{}) {
	mp.data = data
}

// GetPreviewData returns preview data for a template
func (mp *MailerPreview) GetPreviewData(templateName string) map[string]interface{} {
	if data, ok := mp.data[templateName]; ok {
		return data
	}
	return map[string]interface{}{}
}

// setupEmailPreviewRoutes configures routes for email template previews
func (mp *MailerPreview) setupEmailPreviewRoutes(router chi.Router) {
	router.Route(mailerPreviewPath, func(r chi.Router) {
		r.Get("/", mp.handleListEmailTemplates)
		r.Get("/{template}", mp.handlePreviewEmailTemplate)
		r.Post("/{template}/send", mp.handleSendTestEmail)
	})
}

// TemplateManager handles template-related operations
type TemplateManager struct {
	fs  fs.FS
	dir string
}

func newTemplateManager(fs fs.FS, dir string) *TemplateManager {
	return &TemplateManager{
		fs:  fs,
		dir: dir,
	}
}

// getTemplates returns a list of available email templates
func (tm *TemplateManager) getTemplates() ([]string, error) {
	templates := []string{}
	err := fs.WalkDir(tm.fs, tm.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && !strings.HasPrefix(path, tm.dir+"/layouts") {
			baseName := strings.TrimPrefix(path, tm.dir+"/")
			baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))

			if !contains(templates, baseName) {
				templates = append(templates, baseName)
			}
		}
		return nil
	})
	return templates, err
}

// PreviewRenderer handles the rendering of preview pages
type PreviewRenderer struct {
	cssProvider *CSSProvider
	jsProvider  *JSProvider
	templateMgr *TemplateManager
}

func newPreviewRenderer(templateMgr *TemplateManager) *PreviewRenderer {
	return &PreviewRenderer{
		cssProvider: newCSSProvider(),
		jsProvider:  newJSProvider(),
		templateMgr: templateMgr,
	}
}

// handleListEmailTemplates lists all available email templates
func (mp *MailerPreview) handleListEmailTemplates(w http.ResponseWriter, r *http.Request) {
	if mp.mailer == nil {
		http.Error(w, "Mailer not configured", http.StatusNotFound)
		return
	}

	templateMgr := newTemplateManager(mp.mailer.TemplateOptions().FS, mp.mailer.TemplateOptions().Dir)
	templates, err := templateMgr.getTemplates()
	if err != nil {
		http.Error(w, "Failed to list templates", http.StatusInternalServerError)
		return
	}

	locales := mp.i18nBundle.Locales()
	selectedTemplate := r.URL.Query().Get("template")
	selectedLocale := r.URL.Query().Get("locale")
	if selectedLocale == "" {
		selectedLocale = defaultLocale
	}

	if selectedTemplate == "" && len(templates) > 0 {
		redirectURL := fmt.Sprintf("%s?template=%s&locale=%s", mailerPreviewPath, templates[0], selectedLocale)
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	renderer := newPreviewRenderer(templateMgr)
	renderer.renderPreviewPage(w, r, templates, locales, selectedTemplate, selectedLocale)
}

// handlePreviewEmailTemplate handles previewing a specific email template
func (mp *MailerPreview) handlePreviewEmailTemplate(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "template")
	if templateName == "" {
		http.Error(w, "Template name is required", http.StatusBadRequest)
		return
	}

	format := r.URL.Query().Get("format")
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = defaultLocale
	}

	data := mp.GetPreviewData(templateName)
	opts := &RenderOptions{Locale: locale}

	html, text, err := mp.mailer.Render(templateName, data, opts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
		return
	}

	var content string
	if format == "text" {
		content = text
		if content == "" {
			content = fmt.Sprintf("No text template found for %s\nCreate a file named %s.txt in your templates directory.", templateName, templateName)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		content = html
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}

	fmt.Fprint(w, content)
}

// handleSendTestEmail handles sending a test email
func (mp *MailerPreview) handleSendTestEmail(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email  string `json:"email"`
		Locale string `json:"locale"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	templateName := chi.URLParam(r, "template")
	data := mp.GetPreviewData(templateName)
	opts := &RenderOptions{Locale: request.Locale}

	msg := EmailMessage{
		To:      []string{request.Email},
		Data:    data,
		Subject: fmt.Sprintf("Preview: %s", templateName),
		From:    "noreply@example.com", // TODO: Make configurable
	}

	if err := mp.mailer.Send(templateName, msg, opts); err != nil {
		mp.logger.Error("failed to send email", "error", err)
		http.Error(w, fmt.Sprintf("Failed to send email: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// CSSProvider handles CSS styling for the preview page
type CSSProvider struct{}

func newCSSProvider() *CSSProvider {
	return &CSSProvider{}
}

func (cp *CSSProvider) getCSS() string {
	return `
		:root {
			--font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
			--font-size-sm: 0.875rem;
			--font-size-base: 1rem;
			--font-weight-medium: 500;
			--spacing-1: 0.25rem;
			--spacing-2: 0.5rem;
			--spacing-3: 0.75rem;
			--spacing-4: 1rem;
			--spacing-8: 2rem;
			--radius: 0.375rem;
			--transition: all 0.2s ease-in-out;
			--primary: #0070f3;
			--primary-hover: #0061df;
			--gray-50: #f9fafb;
			--gray-100: #f3f4f6;
			--gray-200: #e5e7eb;
			--gray-300: #d1d5db;
			--gray-400: #9ca3af;
			--gray-700: #374151;
			--gray-800: #1f2937;
			--sidebar-width: 240px;
		}

		* {
			box-sizing: border-box;
			margin: 0;
			padding: 0;
		}

		body {
			font-family: var(--font-family);
			font-size: var(--font-size-base);
			line-height: 1.5;
			height: 100vh;
			display: flex;
			flex-direction: column;
			overflow: hidden;
		}

		header {
			padding: var(--spacing-4);
			background: white;
			border-bottom: 1px solid var(--gray-200);
		}

		h1 {
			font-size: 1.25rem;
			font-weight: var(--font-weight-medium);
			color: var(--gray-800);
		}

		main {
			display: flex;
			flex: 1;
			overflow: hidden;
			gap: 1px;
			background: var(--gray-200);
		}

		.sidebar {
			width: var(--sidebar-width);
			background: white;
			overflow-y: auto;
			display: flex;
			flex-direction: column;
		}

		.template-list {
			list-style: none;
			padding: var(--spacing-3);
			margin: 0;
		}

		.template-item:not(:last-child) {
			margin-bottom: var(--spacing-1);
		}

		.template-link {
			display: block;
			padding: var(--spacing-2) var(--spacing-3);
			color: var(--gray-700);
			text-decoration: none;
			border-radius: var(--radius);
			font-size: var(--font-size-sm);
			font-weight: var(--font-weight-medium);
			transition: var(--transition);
		}

		.template-link:hover {
			background: var(--gray-50);
			color: var(--gray-800);
		}

		.template-link.active {
			background: var(--primary);
			color: white;
		}

		.preview-container {
			flex: 1;
			display: flex;
			flex-direction: column;
			background: white;
		}

		.preview-toolbar {
			display: flex;
			flex-wrap: wrap;
			gap: var(--spacing-3);
			padding: var(--spacing-4);
			background: white;
			border-bottom: 1px solid var(--gray-200);
		}

		.preview-controls {
			display: flex;
			gap: var(--spacing-2);
			align-items: center;
		}

		.preview-button {
			padding: var(--spacing-2) var(--spacing-4);
			border: 1px solid var(--gray-300);
			border-radius: var(--radius);
			background: white;
			color: var(--gray-700);
			font-size: var(--font-size-sm);
			cursor: pointer;
			transition: var(--transition);
		}

		.preview-button:hover {
			border-color: var(--primary);
			color: var(--primary);
		}

		.preview-button.active {
			background: var(--primary);
			color: white;
			border-color: var(--primary);
		}

		.preview-frame-container {
			flex: 1;
			display: flex;
			overflow: hidden;
		}

		.preview-frame {
			flex: 1;
			border: none;
		}

		.preview-frame.mobile {
			max-width: 375px;
			margin: 20px auto;
			box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.1),
					0 2px 8px -2px rgba(0, 0, 0, 0.1),
					0 8px 24px -4px rgba(0, 0, 0, 0.15);
			border-radius: var(--radius);
			background: white;
		}

		.preview-frame.tablet {
			max-width: 768px;
			margin: 20px auto;
			box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.1),
					0 2px 8px -2px rgba(0, 0, 0, 0.1),
					0 8px 24px -4px rgba(0, 0, 0, 0.15);
			border-radius: var(--radius);
			background: white;
		}

		.email-actions {
			margin-left: auto;
		}

		.send-test-button {
			background: var(--primary);
			color: white;
			border: none;
			padding: var(--spacing-2) var(--spacing-4);
			border-radius: var(--radius);
			cursor: pointer;
			font-size: var(--font-size-sm);
			transition: var(--transition);
		}

		.send-test-button:hover {
			background: var(--primary-hover);
			transform: translateY(-1px);
		}

		.preview-tabs {
			display: flex;
			gap: var(--spacing-1);
			padding: 0 var(--spacing-4);
			background: var(--gray-50);
			border-bottom: 1px solid var(--gray-200);
		}

		.preview-tab {
			padding: var(--spacing-3) var(--spacing-4);
			color: var(--gray-600);
			border: none;
			background: none;
			cursor: pointer;
			font-size: var(--font-size-sm);
			border-bottom: 2px solid transparent;
		}

		.preview-tab:hover {
			color: var(--gray-800);
		}

		.preview-tab.active {
			color: var(--primary);
			border-bottom-color: var(--primary);
		}

		select {
			padding: var(--spacing-2) var(--spacing-8) var(--spacing-2) var(--spacing-3);
			border: 1px solid var(--gray-300);
			border-radius: var(--radius);
			min-width: 140px;
			font-size: var(--font-size-sm);
			color: var(--gray-700);
			appearance: none;
			transition: var(--transition);
		}

		select:hover {
			border-color: var(--gray-400);
		}

		select:focus {
			outline: none;
			border-color: var(--primary);
			box-shadow: 0 0 0 1px var(--primary);
		}

		@media (max-width: 768px) {
			main {
				flex-direction: column;
			}

			.sidebar {
				width: 100%;
				max-height: 200px;
			}

			.preview-container {
				height: calc(100vh - 200px);
			}
		}

		.format-toggle {
			display: flex;
			gap: var(--spacing-2);
			margin-bottom: var(--spacing-4);
		}

		.format-toggle button {
			padding: var(--spacing-2) var(--spacing-4);
			border: 1px solid var(--gray-300);
			border-radius: var(--radius-sm);
			background: white;
			color: var(--gray-600);
			cursor: pointer;
			transition: var(--transition);
		}

		.format-toggle button.active {
			background: var(--primary);
			border-color: var(--primary);
			color: white;
		}

		.format-toggle button:hover:not(.active) {
			border-color: var(--primary);
			color: var(--primary);
		}
	`
}

// JSProvider handles JavaScript functionality for the preview page
type JSProvider struct{}

func newJSProvider() *JSProvider {
	return &JSProvider{}
}

func (jp *JSProvider) getJS() string {
	return `
        function updatePreview(template, locale) {
            const links = document.querySelectorAll('.template-link');
            const container = document.querySelector('.preview-container');

            links.forEach(link => {
                if (link.dataset.template === template) {
                    link.classList.add('active');
                } else {
                    link.classList.remove('active');
                }
            });

            const url = new URL(window.location);
            url.searchParams.set('template', template);
            url.searchParams.set('locale', locale);
            if (!url.searchParams.has('format')) {
                url.searchParams.set('format', 'html');
            }
            window.history.pushState({}, '', url);

            updatePreviewFrame(template, locale);
        }

        function updatePreviewFrame(template, locale) {
            const container = document.querySelector('.preview-frame-container');
            const url = new URL(window.location);
            const viewMode = url.searchParams.get('view') || 'desktop';
            const format = url.searchParams.get('format') || 'html';

            const activeButton = document.querySelector('.preview-button.active');
            if (!activeButton || activeButton.dataset.mode !== viewMode) {
                document.querySelectorAll('.preview-button').forEach(btn => btn.classList.remove('active'));
                document.querySelector('.preview-button[data-mode="' + viewMode + '"]').classList.add('active');
            }

            const activeTab = document.querySelector('.preview-tab.active');
            if (!activeTab || activeTab.dataset.format !== format) {
                document.querySelectorAll('.preview-tab').forEach(tab => tab.classList.remove('active'));
                document.querySelector('.preview-tab[data-format="' + format + '"]').classList.add('active');
            }

            const existing = container.querySelector('.preview-frame');
            if (existing) existing.remove();

            const newFrame = document.createElement('iframe');
            newFrame.id = 'preview-frame';
            newFrame.className = 'preview-frame ' + viewMode;

            const previewUrl = new URL('/mailer/preview/' + template, window.location.origin);
            previewUrl.searchParams.set('locale', locale);
            previewUrl.searchParams.set('format', format);
            if (format === 'text') {
                previewUrl.searchParams.set('plain', 'true');
            } else {
                previewUrl.searchParams.set('html', 'true');
            }
            previewUrl.searchParams.set('view', viewMode);

            newFrame.src = previewUrl.toString();
            container.appendChild(newFrame);
        }

        function handleLocaleChange(select) {
            const template = new URL(window.location).searchParams.get('template');
            if (template) {
                updatePreview(template, select.value);
            }
        }

        function setViewMode(button, mode) {
            document.querySelectorAll('.preview-button').forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');

            const url = new URL(window.location);
            url.searchParams.set('view', mode);
            window.history.pushState({}, '', url);

            updatePreviewFrame(
                new URL(window.location).searchParams.get('template'),
                document.getElementById('locale-select').value
            );
        }

        function setFormat(button, format) {
            document.querySelectorAll('.preview-tab').forEach(tab => tab.classList.remove('active'));
            button.classList.add('active');

            const url = new URL(window.location);
            url.searchParams.set('format', format);
            window.history.pushState({}, '', url);

            updatePreviewFrame(
                new URL(window.location).searchParams.get('template'),
                document.getElementById('locale-select').value
            );
        }

        function initState() {
            const url = new URL(window.location);

            const viewMode = url.searchParams.get('view');
            if (!viewMode) {
                url.searchParams.set('view', 'desktop');
                window.history.replaceState({}, '', url);
            }
            const button = document.querySelector('.preview-button[data-mode="' + (viewMode || 'desktop') + '"]');
            if (button) {
                document.querySelectorAll('.preview-button').forEach(btn => btn.classList.remove('active'));
                button.classList.add('active');
            }

            const format = url.searchParams.get('format');
            if (!format) {
                url.searchParams.set('format', 'html');
                window.history.replaceState({}, '', url);
            }
            const tab = document.querySelector('.preview-tab[data-format="' + (format || 'html') + '"]');
            if (tab) {
                document.querySelectorAll('.preview-tab').forEach(t => t.classList.remove('active'));
                tab.classList.add('active');
            }
        }

        window.addEventListener('popstate', () => {
            const url = new URL(window.location);
            const template = url.searchParams.get('template');
            const locale = url.searchParams.get('locale') || '%s';
            const viewMode = url.searchParams.get('view');
            const format = url.searchParams.get('format');

            if (!viewMode) {
                url.searchParams.set('view', 'desktop');
            }
            if (!format) {
                url.searchParams.set('format', 'html');
            }
            window.history.replaceState({}, '', url);

            const button = document.querySelector('.preview-button[data-mode="' + (viewMode || 'desktop') + '"]');
            if (button) {
                document.querySelectorAll('.preview-button').forEach(btn => btn.classList.remove('active'));
                button.classList.add('active');
            }

            const tab = document.querySelector('.preview-tab[data-format="' + (format || 'html') + '"]');
            if (tab) {
                document.querySelectorAll('.preview-tab').forEach(t => t.classList.remove('active'));
                tab.classList.add('active');
            }

            if (template) {
                updatePreview(template, locale);
            }
        });

        document.addEventListener('DOMContentLoaded', initState);

        function sendTestEmail() {
            const template = new URL(window.location).searchParams.get('template');
            const locale = document.getElementById('locale-select').value;
            const email = prompt('Enter email address for test:');

            if (email) {
                fetch('/mailer/preview/' + template + '/send', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        email: email,
                        locale: locale
                    })
                })
                .then(response => {
                    if (response.ok) {
                        alert('Test email sent successfully!');
                    } else {
                        alert('Failed to send test email. Please try again.');
                    }
                })
                .catch(error => {
                    alert('Error sending test email: ' + error);
                });
            }
        }
    `
}

// renderPreviewPage renders the complete preview page
func (pr *PreviewRenderer) renderPreviewPage(w http.ResponseWriter, r *http.Request, templates []string, locales []string, selectedTemplate string, selectedLocale string) {
	format := r.URL.Query().Get("format")
	if format == "" {
		if r.URL.Query().Get("plain") == "true" {
			format = "text"
		} else if r.URL.Query().Get("html") == "true" {
			format = "html"
		} else {
			format = "html"
		}
	}

	viewMode := r.URL.Query().Get("view")
	if viewMode == "" {
		viewMode = "desktop"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Mailer Preview</title>
    <style>%s</style>
</head>
<body>
    <header>
        <h1>Mailer Preview</h1>
    </header>
    <main>
        %s
        <div class="preview-container">
            %s
            %s
        </div>
    </main>
    <script>%s</script>
</body>
</html>`,
		pr.cssProvider.getCSS(),
		pr.renderSidebar(templates, selectedTemplate),
		pr.renderPreviewToolbar(locales, selectedLocale, viewMode, format),
		pr.renderPreviewFrame(selectedTemplate, selectedLocale, viewMode, format),
		pr.jsProvider.getJS(),
	)
}

// renderSidebar renders the sidebar HTML
func (pr *PreviewRenderer) renderSidebar(templates []string, selectedTemplate string) string {
	var options strings.Builder
	options.WriteString(`<div class="sidebar">
		<ul class="template-list">`)
	for _, template := range templates {
		selected := ""
		if template == selectedTemplate {
			selected = " active"
		}
		options.WriteString(fmt.Sprintf(`
			<li class="template-item">
				<a href="#"
				   class="template-link%s"
				   data-template="%s"
				   onclick="updatePreview('%s', document.getElementById('locale-select').value); return false;">
					%s
				</a>
			</li>`,
			selected,
			template,
			template,
			template,
		))
	}
	options.WriteString(`
		</ul>
	</div>`)
	return options.String()
}

// renderPreviewToolbar renders the preview toolbar HTML
func (pr *PreviewRenderer) renderPreviewToolbar(locales []string, selectedLocale string, viewMode string, format string) string {
	var options strings.Builder
	options.WriteString(`
        <div class="preview-toolbar">
            <div class="preview-controls">
                <select id="locale-select" onchange="handleLocaleChange(this)">`)

	for _, locale := range locales {
		selected := ""
		if locale == selectedLocale {
			selected = " selected"
		}
		options.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`,
			locale,
			selected,
			locale,
		))
	}

	desktopActive := ""
	if viewMode == "desktop" {
		desktopActive = " active"
	}
	tabletActive := ""
	if viewMode == "tablet" {
		tabletActive = " active"
	}
	mobileActive := ""
	if viewMode == "mobile" {
		mobileActive = " active"
	}
	htmlActive := ""
	if format == "html" {
		htmlActive = " active"
	}
	textActive := ""
	if format == "text" {
		textActive = " active"
	}

	options.WriteString(fmt.Sprintf(`
                    </select>
                    <button class="preview-button%s" data-mode="desktop" onclick="setViewMode(this, 'desktop')">Desktop</button>
                    <button class="preview-button%s" data-mode="tablet" onclick="setViewMode(this, 'tablet')">Tablet</button>
                    <button class="preview-button%s" data-mode="mobile" onclick="setViewMode(this, 'mobile')">Mobile</button>
                </div>
                <div class="email-actions">
                    <button class="send-test-button" onclick="sendTestEmail()">Send Test Email</button>
                </div>
            </div>
            <div class="preview-tabs">
                <button class="preview-tab%s" data-format="html" onclick="setFormat(this, 'html')">HTML</button>
                <button class="preview-tab%s" data-format="text" onclick="setFormat(this, 'text')">Plain Text</button>
            </div>
            <div class="preview-frame-container">`,
		desktopActive,
		tabletActive,
		mobileActive,
		htmlActive,
		textActive,
	))

	return options.String()
}

// renderPreviewFrame renders the preview frame HTML
func (pr *PreviewRenderer) renderPreviewFrame(selectedTemplate string, selectedLocale string, viewMode string, format string) string {
	if selectedTemplate != "" {
		formatParam := "html"
		if format == "text" {
			formatParam = "plain"
		}
		return fmt.Sprintf(`
                <iframe id="preview-frame" class="preview-frame %s" src="/mailer/preview/%s?locale=%s&format=%s&%s=true"></iframe>`,
			viewMode,
			selectedTemplate,
			selectedLocale,
			format,
			formatParam,
		)
	}
	return `
                <div class="no-preview">
                    Select a template from the sidebar to preview
                </div>`
}

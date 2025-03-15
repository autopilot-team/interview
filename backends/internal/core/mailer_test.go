package core

import (
	"autopilot/backends/internal/types"
	"embed"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wneessen/go-mail"
)

var (
	//go:embed all:testdata
	mailerFS embed.FS
	smtpUrl  = "smtp://localhost:1025"
)

func TestNewMail(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	i18nBundle, err := NewI18nBundle(mailerFS, "testdata/locales")
	require.NoError(t, err)

	tests := []struct {
		name    string
		opts    *MailTemplateOptions
		wantErr bool
	}{
		{
			name: "should initialize with valid configuration",
			opts: &MailTemplateOptions{
				FS:  mailerFS,
				Dir: "testdata/templates",
				ExtraFuncs: []template.FuncMap{
					{
						"fail": func(msg string) (string, error) {
							return "", fmt.Errorf("%s", msg)
						},
						"upper": strings.ToUpper,
					},
				},
				Layout: "layouts/test",
			},
			wantErr: false,
		},
		{
			name:    "should reject nil options",
			opts:    nil,
			wantErr: true,
		},
		{
			name: "should reject missing filesystem",
			opts: &MailTemplateOptions{
				Dir:    "testdata/templates",
				Layout: "layouts/test",
			},
			wantErr: true,
		},
		{
			name: "should reject missing path",
			opts: &MailTemplateOptions{
				FS:     mailerFS,
				Layout: "layouts/test",
			},
			wantErr: true,
		},
		{
			name: "should reject missing layout",
			opts: &MailTemplateOptions{
				FS:  mailerFS,
				Dir: "testdata/templates",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailer, err := NewMail(MailOptions{
				I18nBundle:      i18nBundle,
				Logger:          logger,
				Mode:            types.DebugMode,
				SmtpUrl:         smtpUrl,
				TemplateOptions: tt.opts,
			})
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, mailer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, mailer)
				assert.Equal(t, tt.opts.Dir, mailer.TemplateOptions().Dir)
				assert.Equal(t, tt.opts.Layout, mailer.TemplateOptions().Layout)
				assert.Equal(t, tt.opts.FS, mailer.TemplateOptions().FS)
				assert.Len(t, mailer.TemplateOptions().ExtraFuncs, len(tt.opts.ExtraFuncs))
				assert.Equal(t, i18nBundle, mailer.i18nBundle)
			}
		})
	}
}

func TestMailer_Render(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	i18nBundle, err := NewI18nBundle(mailerFS, "testdata/locales")
	require.NoError(t, err)

	mailer, err := NewMail(MailOptions{
		I18nBundle: i18nBundle,
		Logger:     logger,
		Mode:       types.DebugMode,
		SmtpUrl:    smtpUrl,
		TemplateOptions: &MailTemplateOptions{
			FS:  mailerFS,
			Dir: "testdata/templates",
			ExtraFuncs: []template.FuncMap{
				{
					"fail": func(msg string) (string, error) {
						return "", fmt.Errorf("%s", msg)
					},
					"upper": strings.ToUpper,
				},
			},
			Layout: "layouts/test",
		},
	})
	require.NoError(t, err)

	tests := []struct {
		name         string
		templateName string
		data         map[string]any
		opts         *RenderOptions
		wantHTML     string
		wantText     string
		wantErr      bool
	}{
		{
			name:         "basic template",
			templateName: "basic",
			data: map[string]any{
				"Name": "John",
			},
			wantHTML: "<p>Hello, John!</p>",
			wantText: "Hello, John!",
			wantErr:  false,
		},
		{
			name:         "with i18n",
			templateName: "welcome",
			data: map[string]any{
				"AppName": "TestApp",
				"Name":    "John",
			},
			wantHTML: "<!DOCTYPE html>\n<html>\n<body>\nWelcome John!\nThank you for joining TestApp.\n\n</body>\n</html>\n",
			wantText: "Welcome John!\nThank you for joining TestApp.\n",
			wantErr:  false,
		},
		{
			name:         "with custom locale",
			templateName: "welcome",
			data: map[string]any{
				"AppName": "TestApp",
				"Name":    "John",
			},
			opts: &RenderOptions{
				Locale: "zh-CN",
			},
			wantHTML: "<!DOCTYPE html>\n<html>\n<body>\nJohn，欢迎您！\n感谢您加入 TestApp。\n\n</body>\n</html>\n",
			wantText: "John，欢迎您！\n感谢您加入 TestApp。\n",
			wantErr:  false,
		},
		{
			name:         "with pluralization",
			templateName: "notifications",
			data: map[string]any{
				"NotificationCount": 5,
			},
			wantHTML: "You have 5 notifications",
			wantText: "You have 5 notifications",
			wantErr:  false,
		},
		{
			name:         "with single notification",
			templateName: "notifications",
			data: map[string]any{
				"NotificationCount": 1,
			},
			wantHTML: "You have 1 notification",
			wantText: "You have 1 notification",
			wantErr:  false,
		},
		{
			name:         "with custom layout",
			templateName: "basic",
			data: map[string]any{
				"Name": "John",
			},
			opts: &RenderOptions{
				Layout: "layouts/custom",
			},
			wantHTML: "<custom><p>Hello, John!</p>\n</custom>\n",
			wantText: "Hello, John!",
			wantErr:  false,
		},
		{
			name:         "with extra template functions",
			templateName: "with-functions",
			data: map[string]any{
				"Text": "hello",
			},
			wantHTML: "HELLO",
			wantText: "HELLO",
			wantErr:  false,
		},
		{
			name:         "template not found",
			templateName: "nonexistent",
			data:         map[string]any{},
			wantErr:      true,
		},
		{
			name:         "nil data map",
			templateName: "basic",
			data:         nil,
			wantErr:      true,
		},
		{
			name:         "missing required template data",
			templateName: "basic",
			data:         map[string]any{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, text, err := mailer.Render(tt.templateName, tt.data, tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			if tt.wantHTML != "" {
				assert.Contains(t, html, tt.wantHTML)
			}

			if tt.wantText != "" {
				assert.Contains(t, text, tt.wantText)
			}
		})
	}
}

func TestMailer_Send(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	i18nBundle, err := NewI18nBundle(mailerFS, "testdata/locales")
	require.NoError(t, err)

	tests := []struct {
		name         string
		templateName string
		msg          EmailMessage
		opts         *RenderOptions
		wantErr      bool
	}{
		{
			name:         "basic email",
			templateName: "basic",
			msg: EmailMessage{
				From:    "test@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "Test Email",
				Data: map[string]any{
					"Name": "John",
				},
			},
			wantErr: false,
		},
		{
			name:         "email with attachment",
			templateName: "basic",
			msg: EmailMessage{
				From:    "test@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "Test Email with Attachment",
				Data: map[string]any{
					"Name": "John",
				},
				Attachments: []*mail.File{
					{
						Name:        "test.txt",
						ContentType: mail.TypeTextPlain,
						Header:      make(map[string][]string),
						Writer: func(w io.Writer) (int64, error) {
							return io.Copy(w, strings.NewReader("test content"))
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:         "missing from address",
			templateName: "basic",
			msg: EmailMessage{
				To:      []string{"recipient@example.com"},
				Subject: "Test Email",
				Data: map[string]any{
					"Name": "John",
				},
			},
			wantErr: true,
		},
		{
			name:         "missing to address",
			templateName: "basic",
			msg: EmailMessage{
				From:    "test@example.com",
				Subject: "Test Email",
				Data: map[string]any{
					"Name": "John",
				},
			},
			wantErr: true,
		},
		{
			name:         "invalid template",
			templateName: "nonexistent",
			msg: EmailMessage{
				From:    "test@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "Test Email",
				Data:    map[string]any{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailer, err := NewMail(MailOptions{
				I18nBundle: i18nBundle,
				Logger:     logger,
				Mode:       types.DebugMode,
				SmtpUrl:    smtpUrl,
				TemplateOptions: &MailTemplateOptions{
					FS:  mailerFS,
					Dir: "testdata/templates",
					ExtraFuncs: []template.FuncMap{
						{
							"fail": func(msg string) (string, error) {
								return "", fmt.Errorf("%s", msg)
							},
							"upper": strings.ToUpper,
						},
					},
					Layout: "layouts/test",
				},
			})
			require.NoError(t, err)

			err = mailer.Send(tt.templateName, tt.msg, tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMailer_ConcurrentRender(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	i18nBundle, err := NewI18nBundle(mailerFS, "testdata/locales")
	require.NoError(t, err)

	mailer, err := NewMail(MailOptions{
		I18nBundle: i18nBundle,
		Logger:     logger,
		Mode:       types.DebugMode,
		SmtpUrl:    smtpUrl,
		TemplateOptions: &MailTemplateOptions{
			FS:     mailerFS,
			Dir:    "testdata/templates",
			Layout: "layouts/test",
			ExtraFuncs: []template.FuncMap{
				{
					"fail": func(msg string) (string, error) {
						return "", fmt.Errorf("%s", msg)
					},
					"upper": strings.ToUpper,
				},
			},
		},
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := mailer.Render("basic", map[string]any{
				"Name": "John",
			}, nil)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()
}

func TestMailer_BulkSend(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	i18nBundle, err := NewI18nBundle(mailerFS, "testdata/locales")
	require.NoError(t, err)

	tests := []struct {
		name         string
		templateName string
		messages     []EmailMessage
		opts         *RenderOptions
		wantErr      bool
	}{
		{
			name:         "successful bulk send with attachments",
			templateName: "basic",
			messages: []EmailMessage{
				{
					From:    "test@example.com",
					To:      []string{"recipient1@example.com"},
					Subject: "Test Email 1",
					Data: map[string]any{
						"Name": "John",
					},
					Attachments: []*mail.File{
						{
							Name:        "test1.txt",
							ContentType: mail.TypeTextPlain,
							Header:      make(map[string][]string),
							Writer: func(w io.Writer) (int64, error) {
								return io.Copy(w, strings.NewReader("test content 1"))
							},
						},
					},
				},
				{
					From:    "test@example.com",
					To:      []string{"recipient2@example.com"},
					Subject: "Test Email 2",
					Data: map[string]any{
						"Name": "Jane",
					},
					Attachments: []*mail.File{
						{
							Name:        "test2.txt",
							ContentType: mail.TypeTextPlain,
							Header:      make(map[string][]string),
							Writer: func(w io.Writer) (int64, error) {
								return io.Copy(w, strings.NewReader("test content 2"))
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:         "empty message list",
			templateName: "basic",
			messages:     []EmailMessage{},
			wantErr:      true,
		},
		{
			name:         "invalid template",
			templateName: "nonexistent",
			messages: []EmailMessage{
				{
					From:    "test@example.com",
					To:      []string{"recipient@example.com"},
					Subject: "Test Email",
					Data:    map[string]any{},
				},
			},
			wantErr: true,
		},
		{
			name:         "missing from address",
			templateName: "basic",
			messages: []EmailMessage{
				{
					To:      []string{"recipient@example.com"},
					Subject: "Test Email",
					Data: map[string]any{
						"Name": "John",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailer, err := NewMail(MailOptions{
				I18nBundle: i18nBundle,
				Logger:     logger,
				Mode:       types.DebugMode,
				SmtpUrl:    smtpUrl,
				TemplateOptions: &MailTemplateOptions{
					FS:  mailerFS,
					Dir: "testdata/templates",
					ExtraFuncs: []template.FuncMap{
						{
							"fail": func(msg string) (string, error) {
								return "", fmt.Errorf("%s", msg)
							},
							"upper": strings.ToUpper,
						},
					},
					Layout: "layouts/test",
				},
			})
			require.NoError(t, err)

			err = mailer.BulkSend(tt.templateName, tt.messages, tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

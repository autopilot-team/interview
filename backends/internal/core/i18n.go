package core

import (
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// I18nBundle wraps the i18n.Bundle and adds functionality to track available locales.
// It provides internationalization support for the application by managing translation
// files and locale-specific messages.
type I18nBundle struct {
	*i18n.Bundle
	locales []string
}

// NewI18nBundle creates a new I18nBundle instance by loading translation files from
// an embedded filesystem. It initializes the bundle with English as the default language
// and loads all JSON translation files from the specified path.
//
// Parameters:
//   - localesFS: An embedded filesystem containing the locale files
//   - localesDir: The directory within the filesystem where locale files are stored
//
// Returns:
//   - *I18nBundle: A new bundle instance with loaded translations
//   - error: An error if loading fails
func NewI18nBundle(localesFS FS, localesDir string) (*I18nBundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	entries, err := localesFS.ReadDir(localesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read locales directory: %w", err)
	}

	locales := []string{}
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			if _, err := bundle.LoadMessageFileFS(localesFS, localesDir+"/"+entry.Name()); err != nil {
				return nil, fmt.Errorf("failed to parse translation file: %w", err)
			}

			locales = append(locales, strings.ReplaceAll(entry.Name(), ".json", ""))
		}
	}

	return &I18nBundle{bundle, locales}, nil
}

// Locales returns a slice of available locale codes that have been loaded into the bundle.
// Each locale code corresponds to a JSON translation file that was successfully loaded.
func (b *I18nBundle) Locales() []string {
	return b.locales
}

// Localizer provides locale-specific translation functionality.
// It wraps the underlying i18n.Localizer and maintains references to
// the bundle and active locales for the current context.
type Localizer struct {
	bundle    *I18nBundle
	locales   []string
	localizer *i18n.Localizer
}

// NewLocalizer creates a new Localizer instance for handling translations in a specific locale context.
// It initializes the localizer with the provided bundle and locale preferences.
//
// Parameters:
//   - bundle: The I18nBundle containing all available translations
//   - locales: Variable number of locale codes in order of preference (e.g., "en", "zh-CN")
//
// Returns:
//   - *Localizer: A new localizer instance configured for the specified locales
func NewLocalizer(bundle *I18nBundle, locales ...string) *Localizer {
	localizer := i18n.NewLocalizer(bundle.Bundle, locales...)
	return &Localizer{bundle, locales, localizer}
}

// T returns a translation function that can be used to localize messages.
// The returned function accepts a message ID and optional arguments for template data.
// It supports two styles of argument passing:
//  1. A single map argument containing all template data
//  2. Key-value pairs as consecutive arguments
//
// Parameters:
//   - messageID: The identifier of the message to translate
//   - args: Variable arguments that can be either a single map or key-value pairs
//
// Returns:
//   - func(string, ...interface{}) (string, error): A function that performs the actual translation
//   - The returned function returns the translated string and any error that occurred
func (l *Localizer) T() func(messageID string, args ...interface{}) (string, error) {
	return func(messageID string, args ...interface{}) (string, error) {
		// If we have exactly one argument and it's a map, use it directly as template data
		if len(args) == 1 {
			if templateData, ok := args[0].(map[string]interface{}); ok {
				lc := &i18n.LocalizeConfig{
					MessageID:    messageID,
					TemplateData: templateData,
				}

				return l.localizer.Localize(lc)
			}
		}

		// Handle the original key-value pair style arguments
		kv := map[string]interface{}{}
		key := ""
		for _, arg := range args {
			if key == "" {
				key, _ = arg.(string)
				if key == "" {
					return "", fmt.Errorf("expected string key but got %#v", arg)
				}
			} else {
				kv[key] = arg
				key = ""
			}
		}

		lc := &i18n.LocalizeConfig{
			Funcs: template.FuncMap{
				"t": func(messageID string, args ...interface{}) (string, error) {
					lc := &i18n.LocalizeConfig{
						MessageID:    messageID,
						TemplateData: kv,
					}

					return l.localizer.Localize(lc)
				},
			},
			MessageID:    messageID,
			PluralCount:  kv["PluralCount"],
			TemplateData: kv,
		}

		return l.localizer.Localize(lc)
	}
}

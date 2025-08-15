package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type localeDiff struct {
	BasePath      string
	ReferenceLang string
	ReferenceFile string
	TargetLang    string
	TargetFile    string
	MissingKeys   []string
	ExtraKeys     []string
}

func NewSyncLocaleCmd(ctx context.Context) *cobra.Command {
	var base string
	var names []string
	var ignore []string
	var refLang string

	cmd := &cobra.Command{
		Use:   "sync-locales",
		Short: "Check and synchronize locale files against en base",
		RunE: func(cmd *cobra.Command, args []string) error {
			var localesPaths []string

			err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					return nil
				}
				trimmed := strings.TrimPrefix(path, base)
				if slices.Contains(ignore, d.Name()) || strings.HasPrefix(trimmed, ".") {
					return fs.SkipDir
				}
				if slices.Contains(names, d.Name()) {
					localesPaths = append(localesPaths, path)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("error finding locales in project: %w", err)
			}

			// Collect all diffs
			var allDiffs []localeDiff
			for _, path := range localesPaths {
				diffs := checkLocalesDir(path, refLang)
				allDiffs = append(allDiffs, diffs...)
			}

			// Styles
			fileStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("75"))

			missingStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("196"))

			keyStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("245"))

			successStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("82")).
				Bold(true)

			if len(allDiffs) == 0 {
				fmt.Println(successStyle.Render("✓ All locale files are synchronized with en base!"))
				return nil
			}

			// Print diffs
			for _, diff := range allDiffs {
				fmt.Println(fileStyle.Render(fmt.Sprintf("--- %s", diff.ReferenceFile)))
				fmt.Println(fileStyle.Render(fmt.Sprintf("+++ %s", diff.TargetFile)))
				fmt.Println(keyStyle.Render(fmt.Sprintf("@@ -%d keys, +%d keys @@", len(diff.MissingKeys), len(diff.ExtraKeys))))

				for _, key := range diff.MissingKeys {
					fmt.Println(missingStyle.Render(fmt.Sprintf("- %s", key)))
				}

				extraStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("82"))
				for _, key := range diff.ExtraKeys {
					fmt.Println(extraStyle.Render(fmt.Sprintf("+ %s", key)))
				}
				fmt.Println()
			}

			return fmt.Errorf("locale files are not in sync with en base")
		},
	}

	cmd.Flags().StringVarP(&base, "base", "b", ".", "Base directory to search for locales")
	cmd.Flags().StringVarP(&refLang, "ref", "r", "en", "reference lang to show diffs from")
	cmd.Flags().StringArrayVarP(&names, "names", "n", []string{"locales"}, "locale data directory name")
	cmd.Flags().StringArrayVarP(&ignore, "ignore", "i", []string{"build", "dist", "node_modules", "vendor"}, "ignore paths")

	return cmd
}

func checkLocalesDir(dir, refLang string) []localeDiff {
	var diffs []localeDiff
	entries, err := os.ReadDir(dir)
	if err != nil {
		return diffs
	}

	var localeFiles []string
	var localeDirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			localeDirs = append(localeDirs, entry.Name())
			continue
		}
		if path.Ext(entry.Name()) == ".json" {
			localeFiles = append(localeFiles, filepath.Join(dir, entry.Name()))
		}
	}

	// Handle flat structure (en.json, zh.json, etc.)
	if len(localeFiles) != 0 {
		var refFile string
		var otherFiles []string

		for _, file := range localeFiles {
			base := strings.TrimSuffix(path.Base(file), path.Ext(file))
			if base == refLang {
				refFile = file
			} else {
				otherFiles = append(otherFiles, file)
			}
		}

		if refFile != "" && len(otherFiles) > 0 {
			diffs = append(diffs, checkFiles(dir, refLang, refFile, otherFiles, false)...)
		}
	}

	// Handle directory structure (en/, zh/, etc.)
	if len(localeDirs) != 0 {
		if !slices.Contains(localeDirs, refLang) {
			return diffs
		}

		// Get all files in en directory
		refFiles, err := os.ReadDir(filepath.Join(dir, refLang))
		if err != nil {
			return diffs
		}

		// For each file in the ref directory, check corresponding files in other locale directories
		for _, refFile := range refFiles {
			if !refFile.IsDir() && path.Ext(refFile.Name()) == ".json" {
				refPath := filepath.Join(dir, refLang, refFile.Name())

				var targetFiles []string
				for _, locale := range localeDirs {
					if locale == refLang {
						continue
					}
					targetFile := filepath.Join(dir, locale, refFile.Name())
					if _, err := os.Stat(targetFile); err == nil {
						targetFiles = append(targetFiles, targetFile)
					}
				}

				if len(targetFiles) > 0 {
					diffs = append(diffs, checkFiles(dir, refLang, refPath, targetFiles, true)...)
				}
			}
		}
	}

	return diffs
}

func checkFiles(base, refLang string, refFile string, targetFiles []string, subdirs bool) []localeDiff {
	var diffs []localeDiff
	buf, err := os.ReadFile(refFile)
	if err != nil {
		return diffs
	}

	var refMap map[string]any
	if err := json.Unmarshal(buf, &refMap); err != nil {
		return diffs
	}

	refDefs := make(map[string]string)
	if err := flattenJSON(refMap, "", refDefs); err != nil {
		return diffs
	}

	refKeys := make(map[string]struct{})
	for k := range refDefs {
		refKeys[k] = struct{}{}
	}

	// Check each target file against ref
	for _, targetFile := range targetFiles {
		buf, err := os.ReadFile(targetFile)
		if err != nil {
			continue
		}

		var targetMap map[string]any
		if err := json.Unmarshal(buf, &targetMap); err != nil {
			continue
		}

		targetDefs := make(map[string]string)
		if err := flattenJSON(targetMap, "", targetDefs); err != nil {
			continue
		}

		targetLang := strings.TrimSuffix(path.Base(targetFile), path.Ext(targetFile))
		if subdirs {
			targetLang = filepath.Base(filepath.Dir(targetFile))
		}

		// Find missing keys (in ref but not in target)
		var missing []string
		for k := range refKeys {
			if _, ok := targetDefs[k]; !ok {
				missing = append(missing, k)
			}
		}
		slices.Sort(missing)

		// Find extra keys (in target but not in ref)
		var extra []string
		for k := range targetDefs {
			if _, ok := refKeys[k]; !ok {
				extra = append(extra, k)
			}
		}
		slices.Sort(extra)

		if len(missing) > 0 || len(extra) > 0 {
			// Use absolute paths for clickable links
			absRefFile, _ := filepath.Abs(refFile)
			absTargetFile, _ := filepath.Abs(targetFile)
			diff := localeDiff{
				BasePath:      base,
				ReferenceLang: refLang,
				ReferenceFile: absRefFile,
				TargetLang:    targetLang,
				TargetFile:    absTargetFile,
				MissingKeys:   missing,
				ExtraKeys:     extra,
			}
			diffs = append(diffs, diff)
		}
	}

	return diffs
}

func flattenJSON(data map[string]any, prefix string, result map[string]string) error {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]any:
			if err := flattenJSON(v, fullKey, result); err != nil {
				return err
			}
		case string:
			result[fullKey] = v
		default:
			return fmt.Errorf("unknown json value of %s (type %T): %v", fullKey, v, v)
		}
	}

	return nil
}

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

func NewSyncLocaleCmd(ctx context.Context) *cobra.Command {
	var base string
	var names []string
	var ignore []string

	cmd := &cobra.Command{
		Use:   "sync-locales",
		Short: "Check and synchronize locale information",
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

			var errs []error
			for _, path := range localesPaths {
				errs = append(errs, checkLocalesDir(path)...)
			}
			for _, err := range errs {
				fmt.Println(err.Error())
			}
			if len(errs) != 0 {
				return fmt.Errorf("not in sync")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&base, "base", "b", ".", "Base directory to search for locales")
	cmd.Flags().StringArrayVarP(&names, "names", "n", []string{"locales"}, "locale data directory name")
	cmd.Flags().StringArrayVarP(&ignore, "ignore", "i", []string{"build", "dist", "node_modules", "vendor"}, "ignore paths")

	return cmd
}

func checkLocalesDir(dir string) (errs []error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []error{fmt.Errorf("could not load directory %q: %w", dir, err)}
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

	if len(localeFiles) != 0 {
		slices.Sort(localeFiles)
		errs = append(errs, checkLangFiles(dir, localeFiles, false)...)
	}
	if len(localeDirs) != 0 {
		nameMap := make(map[string]struct{})
		for _, locale := range localeDirs {
			files, err := os.ReadDir(filepath.Join(dir, locale))
			if err != nil {
				return []error{fmt.Errorf("could not load directory %q: %w", filepath.Join(dir, locale), err)}
			}
			for _, name := range files {
				nameMap[name.Name()] = struct{}{}
			}
		}

		names := make([]string, 0, len(nameMap))
		for k := range nameMap {
			names = append(names, k)
		}
		slices.Sort(names)

		for _, name := range names {
			var files []string
			for _, l := range localeDirs {
				files = append(files, filepath.Join(dir, l, name))
			}
			errs = append(errs, checkLangFiles(dir, files, true)...)
		}
	}
	return errs
}

func checkLangFiles(base string, files []string, subdirs bool) (errs []error) {
	fileDefs := make([]map[string]string, len(files))
	allKeys := make(map[string]struct{})
	for i, f := range files {
		buf, err := os.ReadFile(f)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not open %q: %w", f, err))
			continue
		}

		var localeMap map[string]any
		if err := json.Unmarshal(buf, &localeMap); err != nil {
			errs = append(errs, fmt.Errorf("could not decode %q: %w", f, err))
			continue
		}

		defs := make(map[string]string)
		if err := flattenJSON(localeMap, "", defs); err != nil {
			errs = append(errs, fmt.Errorf("could not decode %q: %w", f, err))
		}
		fileDefs[i] = defs
		for k := range defs {
			allKeys[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	langs := make([]string, 0, len(files))
	var langLength int
	for _, f := range files {
		lang := strings.TrimSuffix(path.Base(f), path.Ext(f))
		if subdirs {
			lang = filepath.Base(filepath.Dir(f))
		}
		langs = append(langs, lang)
		langLength = max(langLength, len(lang))
	}

	for i, lang := range langs {
		if fileDefs[i] == nil {
			continue
		}
		for _, k := range keys {
			if _, ok := fileDefs[i][k]; !ok {
				if subdirs {
					f := files[i]
					name := strings.TrimSuffix(filepath.Base(f), path.Ext(f))
					errs = append(errs, fmt.Errorf("%s: missing in %*s/%s: %s", base, langLength, lang, name, k))
				} else {
					errs = append(errs, fmt.Errorf("%s: missing in %*s: %s", base, langLength, lang, k))
				}
			}
		}
	}
	if len(errs) != 0 {
		errs = append(errs, errors.New("")) // empty line
	}

	return errs
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
			return fmt.Errorf("unknown json value of  %s (type %T): %v", fullKey, v, v)
		}
	}
	return nil
}

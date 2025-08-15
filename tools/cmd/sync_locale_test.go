package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestFlattenJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected map[string]string
	}{
		{
			name: "simple object",
			input: map[string]any{
				"hello": "world",
				"foo":   "bar",
			},
			expected: map[string]string{
				"hello": "world",
				"foo":   "bar",
			},
		},
		{
			name: "nested object",
			input: map[string]any{
				"user": map[string]any{
					"name":  "John",
					"email": "john@example.com",
				},
				"app": map[string]any{
					"title": "My App",
				},
			},
			expected: map[string]string{
				"user.name":  "John",
				"user.email": "john@example.com",
				"app.title":  "My App",
			},
		},
		{
			name: "deeply nested",
			input: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": "value",
					},
				},
			},
			expected: map[string]string{
				"level1.level2.level3": "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]string)
			err := flattenJSON(tt.input, "", result)
			if err != nil {
				t.Errorf("flattenJSON() error = %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("flattenJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFlattenJSONError(t *testing.T) {
	input := map[string]any{
		"invalid": 123,
	}
	result := make(map[string]string)
	err := flattenJSON(input, "", result)
	if err == nil {
		t.Error("flattenJSON() should return error for non-string values")
	}
}

func TestCheckFilesAgainstEn(t *testing.T) {
	tmpDir := t.TempDir()

	enContent := `{
		"hello": "Hello",
		"user": {
			"name": "Name",
			"email": "Email"
		}
	}`

	zhContent := `{
		"hello": "你好",
		"user": {
			"name": "姓名"
		}
	}`

	frContent := `{
		"hello": "Bonjour",
		"user": {
			"name": "Nom",
			"email": "Email"
		}
	}`

	enFile := filepath.Join(tmpDir, "en.json")
	zhFile := filepath.Join(tmpDir, "zh.json")
	frFile := filepath.Join(tmpDir, "fr.json")

	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(zhFile, []byte(zhContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(frFile, []byte(frContent), 0644); err != nil {
		t.Fatal(err)
	}

	targetFiles := []string{zhFile, frFile}
	diffs := checkFiles(tmpDir, "en", enFile, targetFiles, false)

	if len(diffs) != 1 {
		t.Errorf("Expected 1 diff, got %d", len(diffs))
	}

	zhDiff := diffs[0]
	if zhDiff.TargetLang != "zh" {
		t.Errorf("Expected target lang 'zh', got %s", zhDiff.TargetLang)
	}

	expectedMissing := []string{"user.email"}
	if !reflect.DeepEqual(zhDiff.MissingKeys, expectedMissing) {
		t.Errorf("Expected missing keys %v, got %v", expectedMissing, zhDiff.MissingKeys)
	}

	if len(zhDiff.ExtraKeys) != 0 {
		t.Errorf("Expected no extra keys, got %v", zhDiff.ExtraKeys)
	}
}

func TestCheckFilesAgainstEnExtraKeys(t *testing.T) {
	tmpDir := t.TempDir()

	enContent := `{
		"hello": "Hello"
	}`

	zhContent := `{
		"hello": "你好",
		"extra": "额外的",
		"another": "另一个"
	}`

	enFile := filepath.Join(tmpDir, "en.json")
	zhFile := filepath.Join(tmpDir, "zh.json")

	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(zhFile, []byte(zhContent), 0644); err != nil {
		t.Fatal(err)
	}

	targetFiles := []string{zhFile}
	diffs := checkFiles(tmpDir, "en", enFile, targetFiles, false)

	if len(diffs) != 1 {
		t.Errorf("Expected 1 diff, got %d", len(diffs))
	}

	diff := diffs[0]
	if diff.TargetLang != "zh" {
		t.Errorf("Expected target lang 'zh', got %s", diff.TargetLang)
	}

	if len(diff.MissingKeys) != 0 {
		t.Errorf("Expected no missing keys, got %v", diff.MissingKeys)
	}

	expectedExtra := []string{"another", "extra"}
	if !reflect.DeepEqual(diff.ExtraKeys, expectedExtra) {
		t.Errorf("Expected extra keys %v, got %v", expectedExtra, diff.ExtraKeys)
	}
}

func TestCheckFilesAgainstEnSubdirs(t *testing.T) {
	tmpDir := t.TempDir()
	enDir := filepath.Join(tmpDir, "en")
	zhDir := filepath.Join(tmpDir, "zh")

	if err := os.MkdirAll(enDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(zhDir, 0755); err != nil {
		t.Fatal(err)
	}

	enContent := `{
		"common": {
			"save": "Save",
			"cancel": "Cancel"
		}
	}`

	zhContent := `{
		"common": {
			"save": "保存"
		}
	}`

	enFile := filepath.Join(enDir, "common.json")
	zhFile := filepath.Join(zhDir, "common.json")
	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(zhFile, []byte(zhContent), 0644); err != nil {
		t.Fatal(err)
	}

	targetFiles := []string{zhFile}
	diffs := checkFiles(tmpDir, "en", enFile, targetFiles, true)
	if len(diffs) != 1 {
		t.Errorf("Expected 1 diff, got %d", len(diffs))
	}

	diff := diffs[0]
	if diff.TargetLang != "zh" {
		t.Errorf("Expected target lang 'zh', got %s", diff.TargetLang)
	}

	expectedMissing := []string{"common.cancel"}
	if !reflect.DeepEqual(diff.MissingKeys, expectedMissing) {
		t.Errorf("Expected missing keys %v, got %v", expectedMissing, diff.MissingKeys)
	}

	if len(diff.ExtraKeys) != 0 {
		t.Errorf("Expected no extra keys, got %v", diff.ExtraKeys)
	}
}

func TestCheckLocalesDirFlat(t *testing.T) {
	tmpDir := t.TempDir()
	enContent := `{"hello": "Hello"}`
	zhContent := `{}`
	enFile := filepath.Join(tmpDir, "en.json")
	zhFile := filepath.Join(tmpDir, "zh.json")

	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(zhFile, []byte(zhContent), 0644); err != nil {
		t.Fatal(err)
	}

	diffs := checkLocalesDir(tmpDir, "en")
	if len(diffs) != 1 {
		t.Errorf("Expected 1 diff, got %d", len(diffs))
	}

	if diffs[0].TargetLang != "zh" {
		t.Errorf("Expected target lang 'zh', got %s", diffs[0].TargetLang)
	}

	expectedMissing := []string{"hello"}
	if !reflect.DeepEqual(diffs[0].MissingKeys, expectedMissing) {
		t.Errorf("Expected missing keys %v, got %v", expectedMissing, diffs[0].MissingKeys)
	}

	if len(diffs[0].ExtraKeys) != 0 {
		t.Errorf("Expected no extra keys, got %v", diffs[0].ExtraKeys)
	}
}

func TestCheckLocalesDirStructured(t *testing.T) {
	tmpDir := t.TempDir()
	enDir := filepath.Join(tmpDir, "en")
	zhDir := filepath.Join(tmpDir, "zh")
	if err := os.MkdirAll(enDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(zhDir, 0755); err != nil {
		t.Fatal(err)
	}

	enContent := `{"title": "Title"}`
	zhContent := `{}`
	enFile := filepath.Join(enDir, "app.json")
	zhFile := filepath.Join(zhDir, "app.json")
	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(zhFile, []byte(zhContent), 0644); err != nil {
		t.Fatal(err)
	}

	diffs := checkLocalesDir(tmpDir, "en")
	if len(diffs) != 1 {
		t.Errorf("Expected 1 diff, got %d", len(diffs))
	}

	if diffs[0].TargetLang != "zh" {
		t.Errorf("Expected target lang 'zh', got %s", diffs[0].TargetLang)
	}
}

func TestCheckLocalesDirNoEnglish(t *testing.T) {
	tmpDir := t.TempDir()
	zhDir := filepath.Join(tmpDir, "zh")
	frDir := filepath.Join(tmpDir, "fr")
	if err := os.MkdirAll(zhDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(frDir, 0755); err != nil {
		t.Fatal(err)
	}

	zhContent := `{"hello": "你好"}`
	frContent := `{"hello": "Bonjour"}`
	zhFile := filepath.Join(zhDir, "app.json")
	frFile := filepath.Join(frDir, "app.json")
	if err := os.WriteFile(zhFile, []byte(zhContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(frFile, []byte(frContent), 0644); err != nil {
		t.Fatal(err)
	}

	diffs := checkLocalesDir(tmpDir, "en")
	if len(diffs) != 0 {
		t.Errorf("Expected 0 diffs when no English base, got %d", len(diffs))
	}
}

func TestCheckLocalesDirSynchronized(t *testing.T) {
	tmpDir := t.TempDir()
	enContent := `{"hello": "Hello", "goodbye": "Goodbye"}`
	zhContent := `{"hello": "你好", "goodbye": "再见"}`
	enFile := filepath.Join(tmpDir, "en.json")
	zhFile := filepath.Join(tmpDir, "zh.json")
	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(zhFile, []byte(zhContent), 0644); err != nil {
		t.Fatal(err)
	}

	diffs := checkLocalesDir(tmpDir, "en")
	if len(diffs) != 0 {
		t.Errorf("Expected 0 diffs for synchronized files, got %d", len(diffs))
	}
}

func TestCheckLocalesDirInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	enContent := `{"hello": "Hello"}`
	zhContent := `{invalid json}`
	enFile := filepath.Join(tmpDir, "en.json")
	zhFile := filepath.Join(tmpDir, "zh.json")
	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(zhFile, []byte(zhContent), 0644); err != nil {
		t.Fatal(err)
	}

	diffs := checkLocalesDir(tmpDir, "en")
	if len(diffs) != 0 {
		t.Errorf("Expected 0 diffs for invalid JSON (should be skipped), got %d", len(diffs))
	}
}

func TestCheckLocalesDirMissingTargetFile(t *testing.T) {
	tmpDir := t.TempDir()
	enDir := filepath.Join(tmpDir, "en")
	zhDir := filepath.Join(tmpDir, "zh")
	if err := os.MkdirAll(enDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(zhDir, 0755); err != nil {
		t.Fatal(err)
	}

	enContent := `{"hello": "Hello"}`
	enFile := filepath.Join(enDir, "app.json")
	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatal(err)
	}

	diffs := checkLocalesDir(tmpDir, "en")
	if len(diffs) != 0 {
		t.Errorf("Expected 0 diffs when target file doesn't exist, got %d", len(diffs))
	}
}

func TestCheckLocalesDirMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	enDir := filepath.Join(tmpDir, "en")
	zhDir := filepath.Join(tmpDir, "zh")
	frDir := filepath.Join(tmpDir, "fr")

	if err := os.MkdirAll(enDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(zhDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(frDir, 0755); err != nil {
		t.Fatal(err)
	}

	enCommon := `{"save": "Save", "cancel": "Cancel"}`
	enAuth := `{"login": "Login", "logout": "Logout"}`

	zhCommon := `{"save": "保存"}`
	zhAuth := `{"login": "登录", "logout": "登出"}`

	frCommon := `{"save": "Enregistrer", "cancel": "Annuler"}`
	frAuth := `{"login": "Connexion"}`

	if err := os.WriteFile(filepath.Join(enDir, "common.json"), []byte(enCommon), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(enDir, "auth.json"), []byte(enAuth), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(zhDir, "common.json"), []byte(zhCommon), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(zhDir, "auth.json"), []byte(zhAuth), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(frDir, "common.json"), []byte(frCommon), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(frDir, "auth.json"), []byte(frAuth), 0644); err != nil {
		t.Fatal(err)
	}

	diffs := checkLocalesDir(tmpDir, "en")

	if len(diffs) != 2 {
		t.Errorf("Expected 2 diffs, got %d", len(diffs))
	}

	var zhCommonDiff, frAuthDiff *localeDiff
	for i := range diffs {
		if diffs[i].TargetLang == "zh" && strings.HasSuffix(diffs[i].TargetFile, "zh/common.json") {
			zhCommonDiff = &diffs[i]
		}
		if diffs[i].TargetLang == "fr" && strings.HasSuffix(diffs[i].TargetFile, "fr/auth.json") {
			frAuthDiff = &diffs[i]
		}
	}

	if zhCommonDiff == nil {
		t.Error("Expected zh common.json diff not found")
	} else {
		expectedMissing := []string{"cancel"}
		if !reflect.DeepEqual(zhCommonDiff.MissingKeys, expectedMissing) {
			t.Errorf("Expected zh common missing keys %v, got %v", expectedMissing, zhCommonDiff.MissingKeys)
		}
		if len(zhCommonDiff.ExtraKeys) != 0 {
			t.Errorf("Expected no extra keys, got %v", zhCommonDiff.ExtraKeys)
		}
	}

	if frAuthDiff == nil {
		t.Error("Expected fr auth.json diff not found")
	} else {
		expectedMissing := []string{"logout"}
		if !reflect.DeepEqual(frAuthDiff.MissingKeys, expectedMissing) {
			t.Errorf("Expected fr auth missing keys %v, got %v", expectedMissing, frAuthDiff.MissingKeys)
		}
		if len(frAuthDiff.ExtraKeys) != 0 {
			t.Errorf("Expected no extra keys, got %v", frAuthDiff.ExtraKeys)
		}
	}
}

func TestCheckFilesAgainstEnBothMissingAndExtra(t *testing.T) {
	tmpDir := t.TempDir()

	enContent := `{
		"common": {
			"save": "Save",
			"cancel": "Cancel",
			"submit": "Submit"
		},
		"auth": {
			"login": "Login"
		}
	}`

	zhContent := `{
		"common": {
			"save": "保存",
			"delete": "删除"
		},
		"auth": {
			"login": "登录",
			"register": "注册"
		},
		"extra": {
			"test": "测试"
		}
	}`

	enFile := filepath.Join(tmpDir, "en.json")
	zhFile := filepath.Join(tmpDir, "zh.json")

	if err := os.WriteFile(enFile, []byte(enContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(zhFile, []byte(zhContent), 0644); err != nil {
		t.Fatal(err)
	}

	targetFiles := []string{zhFile}
	diffs := checkFiles(tmpDir, "en", enFile, targetFiles, false)

	if len(diffs) != 1 {
		t.Errorf("Expected 1 diff, got %d", len(diffs))
	}

	diff := diffs[0]
	if diff.TargetLang != "zh" {
		t.Errorf("Expected target lang 'zh', got %s", diff.TargetLang)
	}

	expectedMissing := []string{"common.cancel", "common.submit"}
	if !reflect.DeepEqual(diff.MissingKeys, expectedMissing) {
		t.Errorf("Expected missing keys %v, got %v", expectedMissing, diff.MissingKeys)
	}

	expectedExtra := []string{"auth.register", "common.delete", "extra.test"}
	if !reflect.DeepEqual(diff.ExtraKeys, expectedExtra) {
		t.Errorf("Expected extra keys %v, got %v", expectedExtra, diff.ExtraKeys)
	}
}

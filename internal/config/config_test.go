package config

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestShouldIgnore(t *testing.T) {
	cfg := defaultConfig()

	testCases := []struct {
		path string
		want bool
	}{
		{"foo_test.go", true},                                     // **/*_test.go
		{"internal/order/order_test.go", true},                    // **/*_test.go
		{filepath.Join("test", "utils", "helper.go"), true},       // test/**
		{filepath.Join("service", "mocks", "user_mock.go"), true}, // **/mocks/**
		{filepath.Join("service", "mockups", "data.go"), false},   // should not ignore
		{filepath.Join("cmd", "main.go"), false},                  // should not ignore
	}

	for _, tc := range testCases {
		got := cfg.ShouldIgnore(tc.path)
		if got != tc.want {
			t.Errorf("ShouldIgnore(%s) = %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	// Save and restore current directory
	startDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(startDir)

	// Create temporary directory for tests
	tmpDir := t.TempDir()

	t.Run("config not found in start dir", func(t *testing.T) {
		// Check that there's no config in the start directory
		if _, err := os.Stat(filepath.Join(startDir, ".unused-interface-methods.yml")); err == nil {
			t.Fatal("config file exists in start dir, test cannot proceed")
		}

		// Should get default config
		cfg, err := LoadConfig("")
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}
		want := defaultConfig()
		if !reflect.DeepEqual(cfg, want) {
			t.Errorf("LoadConfig() = %v, want %v", cfg, want)
		}
	})

	t.Run("config found in temp dir", func(t *testing.T) {
		// Change to temporary directory
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}

		// Create config in temporary directory
		content := []byte(`ignore:
  - "vendor/**"
  - "**/*.pb.go"`)
		if err := os.WriteFile(".unused-interface-methods.yml", content, 0644); err != nil {
			t.Fatal(err)
		}

		// Check that file exists
		if _, err := os.Stat(".unused-interface-methods.yml"); err != nil {
			t.Fatal("config file was not created")
		}

		// Load config
		cfg, err := LoadConfig("")
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		want := &Config{
			Ignore: []string{
				"vendor/**",
				"**/*.pb.go",
			},
		}
		if !reflect.DeepEqual(cfg, want) {
			t.Errorf("LoadConfig() = %v, want %v", cfg, want)
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		// Change to temporary directory
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}

		// Create invalid config
		content := []byte(`ignore: [`)
		if err := os.WriteFile(".unused-interface-methods.yml", content, 0644); err != nil {
			t.Fatal(err)
		}

		_, err := LoadConfig("")
		if err == nil {
			t.Error("LoadConfig() error = nil, want error for invalid yaml")
		}
	})

	t.Run("explicit config path", func(t *testing.T) {
		// Create config in non-standard location
		content := []byte(`ignore:
  - "custom/**"`)
		customPath := filepath.Join(tmpDir, "custom.yml")
		if err := os.WriteFile(customPath, content, 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadConfig(customPath)
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		want := &Config{
			Ignore: []string{
				"custom/**",
			},
		}
		if !reflect.DeepEqual(cfg, want) {
			t.Errorf("LoadConfig() = %v, want %v", cfg, want)
		}
	})

	t.Run("permission denied", func(t *testing.T) {
		// Skip this test on Windows as permission handling is different
		if runtime.GOOS == "windows" {
			t.Skip("Skipping permission test on Windows")
		}

		// Change to temporary directory
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}

		// Create subdirectory without read permissions
		noAccessDir := filepath.Join(tmpDir, "noaccess")
		if err := os.Mkdir(noAccessDir, 0700); err != nil {
			t.Fatal(err)
		}

		// Create config
		content := []byte(`ignore:
  - "vendor/**"`)
		configPath := filepath.Join(noAccessDir, ".unused-interface-methods.yml")
		if err := os.WriteFile(configPath, content, 0644); err != nil {
			t.Fatal(err)
		}

		// Remove read permissions from directory
		if err := os.Chmod(noAccessDir, 0000); err != nil {
			t.Fatal(err)
		}

		// Try to load config
		_, err := LoadConfig(configPath)
		if err == nil {
			t.Error("LoadConfig() error = nil, want error for permission denied")
		}

		// Restore permissions for cleanup
		os.Chmod(noAccessDir, 0700)
	})
}

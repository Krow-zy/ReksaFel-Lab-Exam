package config

import (
	"path/filepath"
	"testing"
)

func TestSaveAndLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	original := &Config{
		Port:        9090,
		SecretKey:   "this-is-a-super-long-testing-secret-key-32-chars",
		IPWhitelist: []string{"127.0.0.1", "10.0.0.1"},
	}

	// 1. Test Save
	err := SaveConfig(configPath, original)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// 2. Test Load
	loaded, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// 3. Compare values
	if loaded.Port != original.Port {
		t.Errorf("expected port %d, got %d", original.Port, loaded.Port)
	}

	if loaded.SecretKey != original.SecretKey {
		t.Errorf("expected secret key '%s', got '%s'", original.SecretKey, loaded.SecretKey)
	}

	if len(loaded.IPWhitelist) != len(original.IPWhitelist) {
		t.Fatalf("expected whitelist length %d, got %d", len(original.IPWhitelist), len(loaded.IPWhitelist))
	}

	for i, ip := range loaded.IPWhitelist {
		if ip != original.IPWhitelist[i] {
			t.Errorf("expected whitelist IP at index %d to be '%s', got '%s'", i, original.IPWhitelist[i], ip)
		}
	}
}

func TestLoadConfig_NonExistent(t *testing.T) {
	_, err := LoadConfig("non-existent-config-file.json")
	if err == nil {
		t.Error("expected error loading non-existent config file, got nil")
	}
}

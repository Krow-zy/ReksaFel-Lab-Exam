package config

import (
	"bytes"
	"os"
	"testing"
)

func TestBackupAndRestoreConfig(t *testing.T) {
	// Create temporary source config
	srcFile, err := os.CreateTemp("", "config_src_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp source config: %v", err)
	}
	defer os.Remove(srcFile.Name())

	content := []byte(`{"auth_key": "testkey123", "admin_url": "http://127.0.0.1:8080"}`)
	if _, err := srcFile.Write(content); err != nil {
		t.Fatalf("Failed to write source config: %v", err)
	}
	srcFile.Close()

	// Create temporary backup path
	backupFile, err := os.CreateTemp("", "config_bak_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp backup path: %v", err)
	}
	backupPath := backupFile.Name()
	backupFile.Close()
	os.Remove(backupPath) // Delete so BackupConfig can create it

	defer os.Remove(backupPath)

	// Test Backup
	if err := BackupConfig(srcFile.Name(), backupPath); err != nil {
		t.Fatalf("BackupConfig failed: %v", err)
	}

	// Verify backup content
	bakData, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("Failed to read backup file: %v", err)
	}
	if !bytes.Equal(bakData, content) {
		t.Error("Backup content does not match source content")
	}

	// Corrupt source
	if err := os.WriteFile(srcFile.Name(), []byte("corrupted"), 0600); err != nil {
		t.Fatalf("Failed to write bad data: %v", err)
	}

	// Test Restore
	if err := RestoreBackup(backupPath, srcFile.Name()); err != nil {
		t.Fatalf("RestoreBackup failed: %v", err)
	}

	// Verify restored content
	restData, err := os.ReadFile(srcFile.Name())
	if err != nil {
		t.Fatalf("Failed to read restored file: %v", err)
	}
	if !bytes.Equal(restData, content) {
		t.Error("Restored content does not match backup content")
	}
}

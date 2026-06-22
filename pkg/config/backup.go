package config

import (
	"fmt"
	"io"
	"os"
)

// BackupConfig copies the source configuration file to a backup path.
func BackupConfig(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer srcFile.Close()

	// Ensure destination file is created with restricted user permissions
	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("create backup: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("copy data: %w", err)
	}

	return nil
}

// RestoreBackup restores configuration from a backup file.
func RestoreBackup(backupPath, configPath string) error {
	return BackupConfig(backupPath, configPath)
}

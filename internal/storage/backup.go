package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Backup struct {
	sourcePath      string
	backupDirectory string
}

func NewBackup(sourcePath, backupDirectory string) *Backup {
	return &Backup{
		sourcePath:      sourcePath,
		backupDirectory: backupDirectory,
	}
}

func (b *Backup) CreateBackup() (string, error) {
	if _, err := os.Stat(b.sourcePath); err != nil {
		return "", fmt.Errorf("source file not found: %w", err)
	}

	if err := os.MkdirAll(b.backupDirectory, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Base(b.sourcePath)
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]

	backupPath := filepath.Join(b.backupDirectory, fmt.Sprintf("%s.%s%s", nameWithoutExt, timestamp, ext))

	sourceData, err := os.ReadFile(b.sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to read source file: %w", err)
	}

	if err := os.WriteFile(backupPath, sourceData, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup file: %w", err)
	}

	return backupPath, nil
}

func (b *Backup) RestoreFromBackup(backupPath string) error {
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	if err := os.WriteFile(b.sourcePath, backupData, 0644); err != nil {
		return fmt.Errorf("failed to restore file: %w", err)
	}

	return nil
}

func (b *Backup) ListBackups() ([]string, error) {
	entries, err := os.ReadDir(b.backupDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []string
	filename := filepath.Base(b.sourcePath)
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]

	for _, entry := range entries {
		if !entry.IsDir() && entry.Name()[:len(nameWithoutExt)] == nameWithoutExt {
			backups = append(backups, filepath.Join(b.backupDirectory, entry.Name()))
		}
	}

	return backups, nil
}

func (b *Backup) DeleteOldBackups(maxBackups int) error {
	backups, err := b.ListBackups()
	if err != nil {
		return err
	}

	if len(backups) <= maxBackups {
		return nil
	}

	// Sort and delete oldest
	for i := 0; i < len(backups)-maxBackups; i++ {
		if err := os.Remove(backups[i]); err != nil {
			return fmt.Errorf("failed to delete old backup: %w", err)
		}
	}

	return nil
}

func (b *Backup) GetLatestBackup() (string, error) {
	backups, err := b.ListBackups()
	if err != nil {
		return "", err
	}

	if len(backups) == 0 {
		return "", fmt.Errorf("no backups found")
	}

	// Last backup is typically the latest due to timestamp naming
	return backups[len(backups)-1], nil
}

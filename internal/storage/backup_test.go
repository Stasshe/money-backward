package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewBackup(t *testing.T) {
	b := NewBackup("/path/to/file", "/path/to/backups")
	if b.sourcePath != "/path/to/file" {
		t.Errorf("expected sourcePath '/path/to/file', got %s", b.sourcePath)
	}
	if b.backupDirectory != "/path/to/backups" {
		t.Errorf("expected backupDirectory '/path/to/backups', got %s", b.backupDirectory)
	}
}

func TestCreateBackup(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "test.json")
	backupDir := filepath.Join(tmpDir, "backups")

	testContent := []byte(`{"test": "data"}`)
	if err := os.WriteFile(sourceFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	b := NewBackup(sourceFile, backupDir)
	backupPath, err := b.CreateBackup()
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	if _, err := os.Stat(backupPath); err != nil {
		t.Fatalf("backup file not created: %v", err)
	}

	content, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("failed to read backup: %v", err)
	}

	if string(content) != string(testContent) {
		t.Error("backup content does not match source")
	}
}

func TestRestoreFromBackup(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "test.json")
	backupDir := filepath.Join(tmpDir, "backups")

	originalContent := []byte(`{"version": 1}`)
	modifiedContent := []byte(`{"version": 2}`)

	if err := os.WriteFile(sourceFile, originalContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	b := NewBackup(sourceFile, backupDir)
	backupPath, err := b.CreateBackup()
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Modify the source file
	if err := os.WriteFile(sourceFile, modifiedContent, 0644); err != nil {
		t.Fatalf("failed to modify source file: %v", err)
	}

	// Restore from backup
	if err := b.RestoreFromBackup(backupPath); err != nil {
		t.Fatalf("RestoreFromBackup failed: %v", err)
	}

	// Verify content is restored
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}

	if string(content) != string(originalContent) {
		t.Error("restored content does not match original backup")
	}
}

func TestListBackups(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "test.json")
	backupDir := filepath.Join(tmpDir, "backups")

	testContent := []byte(`{"test": "data"}`)
	if err := os.WriteFile(sourceFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	b := NewBackup(sourceFile, backupDir)

	// Create multiple backups
	for i := 0; i < 3; i++ {
		_, err := b.CreateBackup()
		if err != nil {
			t.Fatalf("CreateBackup failed: %v", err)
		}
	}

	backups, err := b.ListBackups()
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}

	if len(backups) < 3 {
		t.Errorf("expected at least 3 backups, got %d", len(backups))
	}
}

func TestDeleteOldBackups(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "test.json")
	backupDir := filepath.Join(tmpDir, "backups")

	testContent := []byte(`{"test": "data"}`)
	if err := os.WriteFile(sourceFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	b := NewBackup(sourceFile, backupDir)

	// Create 5 backups
	for i := 0; i < 5; i++ {
		_, err := b.CreateBackup()
		if err != nil {
			t.Fatalf("CreateBackup failed: %v", err)
		}
	}

	backupsBefore, _ := b.ListBackups()
	if len(backupsBefore) < 5 {
		t.Fatalf("expected at least 5 backups before cleanup")
	}

	// Delete old backups, keeping only 2
	if err := b.DeleteOldBackups(2); err != nil {
		t.Fatalf("DeleteOldBackups failed: %v", err)
	}

	backupsAfter, _ := b.ListBackups()
	if len(backupsAfter) > 2 {
		t.Errorf("expected 2 or fewer backups after cleanup, got %d", len(backupsAfter))
	}
}

func TestGetLatestBackup(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "test.json")
	backupDir := filepath.Join(tmpDir, "backups")

	testContent := []byte(`{"test": "data"}`)
	if err := os.WriteFile(sourceFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	b := NewBackup(sourceFile, backupDir)

	// Create backup
	firstBackup, _ := b.CreateBackup()

	latestBackup, err := b.GetLatestBackup()
	if err != nil {
		t.Fatalf("GetLatestBackup failed: %v", err)
	}

	if latestBackup != firstBackup {
		t.Errorf("expected latest backup to be %s, got %s", firstBackup, latestBackup)
	}
}

func TestListBackupsNonexistentDir(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "test.json")
	backupDir := filepath.Join(tmpDir, "nonexistent")

	testContent := []byte(`{"test": "data"}`)
	if err := os.WriteFile(sourceFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	b := NewBackup(sourceFile, backupDir)
	backups, err := b.ListBackups()
	if err != nil {
		t.Fatalf("ListBackups should not error on nonexistent dir: %v", err)
	}

	if len(backups) != 0 {
		t.Errorf("expected 0 backups, got %d", len(backups))
	}
}

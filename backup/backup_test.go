package backup

import (
	"testing"

	"github.com/Verthandii/palworld-go/config"
)

func TestBackup(t *testing.T) {
	backup := New(&config.Config{
		GamePath:   "./tmp/src",
		BackupPath: "./tmp/dest",
	})
	backup.backup()
}

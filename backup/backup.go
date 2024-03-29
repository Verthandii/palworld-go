package backup

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/logger"
)

var log = logger.NewLogger("Backup")

type Backup struct {
	c *config.Config
}

func New(c *config.Config) *Backup {
	return &Backup{
		c: c,
	}
}

func (b *Backup) Schedule(ctx context.Context) {
	if b.c.BackupInterval <= 0 {
		return
	}

	log.Printf("启动成功，定时备份服务器数据\n")

	duration := time.Duration(b.c.BackupInterval) * time.Second
	ticker := time.NewTicker(duration)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.backup()
			ticker.Reset(duration)
		}
	}
}

func (b *Backup) backup() {
	nowAt := time.Now()
	currentDate := nowAt.Format("2006-01-02")
	currentTime := nowAt.Format("15-04-05")
	backupDir := filepath.Join(b.c.BackupPath, currentDate, currentTime)
	if err := makeDir(backupDir); err != nil {
		log.Printf("备份文件夹创建失败【%v】\n", err)
		return
	}

	src := filepath.Join(b.c.GamePath, "Pal", "Saved")
	dst := filepath.Join(backupDir, "Pal", "Saved")

	if err := copyDir(src, dst); err != nil {
		log.Printf("备份失败【%v】\n", err)
	} else {
		log.Printf("成功备份至【%s】\n", backupDir)
	}
}

func makeDir(backupDir string) error {
	_, err := os.Stat(backupDir)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(backupDir, 0777); err != nil {
		return err
	}
	return nil
}

func copyDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return &os.PathError{Op: "copy", Path: src, Err: os.ErrInvalid}
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return &os.PathError{Op: "copy", Path: src, Err: os.ErrInvalid}
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

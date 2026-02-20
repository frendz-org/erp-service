package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const dateFormat = "2006-01-02"

type dailyFileWriter struct {
	dir       string
	prefix    string
	compress  bool
	maxAge    int
	retainAll bool

	mu    sync.Mutex
	file  atomic.Pointer[os.File]
	today atomic.Value
}

func newDailyFileWriter(dir, prefix string, compress bool, maxAge int, retainAll bool) (*dailyFileWriter, error) {
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	w := &dailyFileWriter{
		dir:       dir,
		prefix:    prefix,
		compress:  compress,
		maxAge:    maxAge,
		retainAll: retainAll,
	}

	if err := w.openForToday(); err != nil {
		return nil, err
	}

	go w.cleanOldFiles()

	return w, nil
}

func (w *dailyFileWriter) Write(p []byte) (n int, err error) {
	today := time.Now().Format(dateFormat)
	if today != w.today.Load().(string) {
		w.mu.Lock()
		if today != w.today.Load().(string) {
			if err := w.rotate(today); err != nil {
				w.mu.Unlock()
				return 0, err
			}
		}
		w.mu.Unlock()
	}

	f := w.file.Load()
	return f.Write(p)
}

func (w *dailyFileWriter) Sync() error {
	if f := w.file.Load(); f != nil {
		return f.Sync()
	}
	return nil
}

func (w *dailyFileWriter) openForToday() error {
	today := time.Now().Format(dateFormat)
	filename := filepath.Join(w.dir, fmt.Sprintf("%s-%s.log", w.prefix, today))

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	w.file.Store(f)
	w.today.Store(today)
	return nil
}

func (w *dailyFileWriter) rotate(newDay string) error {
	oldFile := w.file.Load()
	oldPath := ""
	if oldFile != nil {
		oldPath = oldFile.Name()
	}

	filename := filepath.Join(w.dir, fmt.Sprintf("%s-%s.log", w.prefix, newDay))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	w.file.Store(f)
	w.today.Store(newDay)

	if oldFile != nil {
		oldFile.Close()
	}

	if oldPath != "" {
		go w.archiveFile(oldPath)
	}

	go w.cleanOldFiles()

	return nil
}

func (w *dailyFileWriter) archiveFile(path string) {
	base := filepath.Base(path)
	dateStr := strings.TrimPrefix(base, w.prefix+"-")
	dateStr = strings.TrimSuffix(dateStr, ".log")

	fileDate, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return
	}

	monthDir := filepath.Join(w.dir, fileDate.Format("2006-01"))
	if err := os.MkdirAll(monthDir, 0750); err != nil {
		return
	}

	if w.compress {
		gzPath := filepath.Join(monthDir, base+".gz")
		if compressFileTo(path, gzPath) {
			os.Remove(path)
		}
	} else {
		dest := filepath.Join(monthDir, base)
		os.Rename(path, dest)
	}
}

func compressFileTo(srcPath, dstPath string) bool {
	src, err := os.Open(srcPath)
	if err != nil {
		return false
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return false
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	gz.Name = filepath.Base(srcPath)
	gz.ModTime = time.Now()

	if _, err := io.Copy(gz, src); err != nil {
		gz.Close()
		os.Remove(dstPath)
		return false
	}

	if err := gz.Close(); err != nil {
		os.Remove(dstPath)
		return false
	}

	return true
}

func (w *dailyFileWriter) cleanOldFiles() {
	if w.retainAll || w.maxAge <= 0 {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -w.maxAge)

	monthDirs, err := filepath.Glob(filepath.Join(w.dir, "[0-9][0-9][0-9][0-9]-[0-9][0-9]"))
	if err != nil {
		return
	}

	for _, monthDir := range monthDirs {
		w.cleanDirectory(monthDir, cutoff)
		w.removeEmptyDir(monthDir)
	}

	w.cleanDirectory(w.dir, cutoff)
}

func (w *dailyFileWriter) cleanDirectory(dir string, cutoff time.Time) {
	pattern := filepath.Join(dir, w.prefix+"-*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	for _, path := range matches {
		base := filepath.Base(path)
		dateStr := strings.TrimPrefix(base, w.prefix+"-")
		dateStr = strings.TrimSuffix(dateStr, ".gz")
		dateStr = strings.TrimSuffix(dateStr, ".log")

		fileDate, err := time.Parse(dateFormat, dateStr)
		if err != nil {
			continue
		}

		if fileDate.Before(cutoff) {
			os.Remove(path)
		}
	}
}

func (w *dailyFileWriter) removeEmptyDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	if len(entries) == 0 {
		os.Remove(dir)
	}
}

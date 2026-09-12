package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sidx1/sure/internal/config"
)

func cachePath(name string) string {
	return filepath.Join(config.CacheDir(), "docs", fmt.Sprintf("%s.txt", name))
}

func Get(command string) (string, bool) {
	path := cachePath(command)
	info, err := os.Stat(path)
	if err != nil {
		return "", false
	}
	
	if time.Since(info.ModTime()) > 7*24*time.Hour {
		_ = os.Remove(path)
		return "", false
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	
	return string(data), true
}

func Set(command, doc string) error {
	dir := filepath.Join(config.CacheDir(), "docs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := cachePath(command)
	return os.WriteFile(path, []byte(doc), 0644)
}

func Clear() error {
	dir := filepath.Join(config.CacheDir(), "docs")
	return os.RemoveAll(dir)
}

func Size() (int, error) {
	dir := filepath.Join(config.CacheDir(), "docs")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	
	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}
	return count, nil
}

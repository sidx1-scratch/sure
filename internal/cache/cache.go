package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sidx1/sure/internal/config"
)

// Storage defines the pluggable doc cache block interface.
// Developers and forks can replace this with Redis, SQLite, Memcached, etc.
type Storage interface {
	Get(command string) (string, bool)
	Set(command, doc string) error
	Clear() error
	Size() (int, error)
}

type FileStorage struct {
	mu sync.RWMutex
}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func cachePath(name string) string {
	return filepath.Join(config.CacheDir(), "docs", fmt.Sprintf("%s.txt", name))
}

func (s *FileStorage) Get(command string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

func (s *FileStorage) Set(command, doc string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := filepath.Join(config.CacheDir(), "docs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := cachePath(command)
	return os.WriteFile(path, []byte(doc), 0644)
}

func (s *FileStorage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := filepath.Join(config.CacheDir(), "docs")
	return os.RemoveAll(dir)
}

func (s *FileStorage) Size() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

var globalStorage Storage = NewFileStorage()

func SetGlobalStorage(s Storage) {
	if s != nil {
		globalStorage = s
	}
}

func Get(command string) (string, bool) {
	return globalStorage.Get(command)
}

func Set(command, doc string) error {
	return globalStorage.Set(command, doc)
}

func Clear() error {
	return globalStorage.Clear()
}

func Size() (int, error) {
	return globalStorage.Size()
}

package discovery

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sidx1/sure/internal/cache"
)

type CommandInfo struct {
	Name     string
	Path     string
	HelpText string
}

// Discoverer defines the pluggable command documentation & system discovery block.
type Discoverer interface {
	DiscoverCommands() ([]CommandInfo, error)
	GetCommandDoc(name string) string
}

type SystemDiscoverer struct{}

func NewSystemDiscoverer() *SystemDiscoverer {
	return &SystemDiscoverer{}
}

func (d *SystemDiscoverer) DiscoverCommands() ([]CommandInfo, error) {
	return DiscoverCommands()
}

func (d *SystemDiscoverer) GetCommandDoc(name string) string {
	return GetCommandDoc(name)
}

var globalDiscoverer Discoverer = NewSystemDiscoverer()

func SetGlobalDiscoverer(d Discoverer) {
	if d != nil {
		globalDiscoverer = d
	}
}

func DiscoverCommands() ([]CommandInfo, error) {
	pathEnv := os.Getenv("PATH")
	dirs := strings.Split(pathEnv, string(os.PathListSeparator))

	var commands []CommandInfo
	seen := make(map[string]bool)

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if !seen[name] {
				seen[name] = true
				commands = append(commands, CommandInfo{
					Name: name,
					Path: filepath.Join(dir, name),
				})
			}
		}
	}

	return commands, nil
}

func GetCommandHelp(name string) (string, error) {
	cmd := exec.Command(name, "--help")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	_ = cmd.Run()
	if out.Len() > 0 {
		return out.String(), nil
	}

	cmd = exec.Command(name, "-h")
	var out2 bytes.Buffer
	cmd.Stdout = &out2
	cmd.Stderr = &out2
	_ = cmd.Run()
	if out2.Len() > 0 {
		return out2.String(), nil
	}

	return "", fmt.Errorf("no help found for %s", name)
}

func GetManPage(name string) (string, error) {
	cmd := exec.Command("man", name)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	_ = cmd.Run()
	if out.Len() > 0 {
		return out.String(), nil
	}
	return "", fmt.Errorf("no man page found for %s", name)
}

func GetCommandDoc(name string) string {
	doc, found := cache.Get(name)
	if found {
		return doc
	}

	doc, err := GetCommandHelp(name)
	if err != nil {
		doc, err = GetManPage(name)
		if err != nil {
			doc = ""
		}
	}

	if len(doc) > 2000 {
		doc = doc[:2000] + "... (truncated)"
	}

	_ = cache.Set(name, doc)
	return doc
}

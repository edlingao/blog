package layout

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type ManifestEntry struct {
	File    string   `json:"file"`
	CSS     []string `json:"css,omitempty"`
	IsEntry bool     `json:"isEntry,omitempty"`
}

type Manifest map[string]ManifestEntry

var (
	manifest     Manifest
	manifestOnce sync.Once
	manifestPath string
)

func SetManifestPath(path string) {
	manifestPath = path
}

func LoadManifest() error {
	var err error
	manifestOnce.Do(func() {
		if manifestPath == "" {
			manifestPath = filepath.Join("web", "static", ".vite", "manifest.json")
		}
		data, readErr := os.ReadFile(manifestPath)
		if readErr != nil {
			err = readErr
			return
		}
		err = json.Unmarshal(data, &manifest)
	})
	return err
}

func GetAssetPath(entry string) string {
	if manifest == nil {
		return ""
	}
	if e, ok := manifest[entry]; ok {
		return "/static/" + e.File
	}
	return ""
}

func GetCSSPaths(entry string) []string {
	if manifest == nil {
		return nil
	}
	if e, ok := manifest[entry]; ok {
		paths := make([]string, len(e.CSS))
		for i, css := range e.CSS {
			paths[i] = "/static/" + css
		}
		return paths
	}
	return nil
}

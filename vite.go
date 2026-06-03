package goravelinertia

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// viteChunk is one entry of a Vite build manifest.json.
type viteChunk struct {
	File    string   `json:"file"`
	Src     string   `json:"src"`
	IsEntry bool     `json:"isEntry"`
	CSS     []string `json:"css"`
}

// Vite resolves asset tags for Inertia's root template, mirroring the Laravel
// @vite directive: it serves assets from the running dev server when one is
// active, otherwise from the hashed files listed in the build manifest.
type Vite struct {
	publicPath string // filesystem dir served at "/", e.g. "public"
	buildDir   string // build output dir under publicPath, e.g. "build"
	hotFile    string // path to the dev-server hot file, e.g. "public/hot"
	devURL     string // explicit dev-server URL; forces dev mode when set

	mu       sync.RWMutex
	manifest map[string]viteChunk
	loaded   bool
}

// NewVite builds a Vite helper. Empty arguments fall back to Laravel-compatible
// defaults: publicPath "public", buildDir "build", hotFile "public/hot".
func NewVite(publicPath, buildDir, hotFile, devURL string) *Vite {
	if publicPath == "" {
		publicPath = "public"
	}
	if buildDir == "" {
		buildDir = "build"
	}
	if hotFile == "" {
		hotFile = filepath.Join(publicPath, "hot")
	}

	return &Vite{
		publicPath: publicPath,
		buildDir:   buildDir,
		hotFile:    hotFile,
		devURL:     devURL,
	}
}

// TemplateFunc returns the function registered as {{ vite "entry" ... }}.
func (v *Vite) TemplateFunc() func(entries ...string) template.HTML {
	return v.Render
}

// Render produces the <script>/<link> tags for the given entry points.
func (v *Vite) Render(entries ...string) template.HTML {
	if url := v.devServerURL(); url != "" {
		return v.devTags(url, entries)
	}

	return v.prodTags(entries)
}

// devServerURL returns the dev-server base URL when running in dev mode: the hot
// file takes precedence (written by the dev server), then the configured devURL.
func (v *Vite) devServerURL() string {
	if data, err := os.ReadFile(v.hotFile); err == nil {
		if url := strings.TrimSpace(string(data)); url != "" {
			return strings.TrimRight(url, "/")
		}
	}

	return strings.TrimRight(v.devURL, "/")
}

func (v *Vite) devTags(url string, entries []string) template.HTML {
	var b strings.Builder

	b.WriteString(moduleScript(url + "/@vite/client"))
	for _, entry := range entries {
		b.WriteString(moduleScript(url + "/" + strings.TrimLeft(entry, "/")))
	}

	return template.HTML(b.String())
}

func (v *Vite) prodTags(entries []string) template.HTML {
	manifest, err := v.loadManifest()
	if err != nil {
		return template.HTML(fmt.Sprintf("<!-- vite: %s -->", template.HTMLEscapeString(err.Error())))
	}

	var b strings.Builder
	for _, entry := range entries {
		chunk, ok := manifest[strings.TrimLeft(entry, "/")]
		if !ok {
			continue
		}

		for _, css := range chunk.CSS {
			b.WriteString(styleLink(v.asset(css)))
		}
		b.WriteString(moduleScript(v.asset(chunk.File)))
	}

	return template.HTML(b.String())
}

// asset returns the public URL for a built file, e.g. "/build/assets/app-h4sh.js".
func (v *Vite) asset(file string) string {
	return "/" + urlJoin(v.buildDir, file)
}

// loadManifest reads and caches the build manifest, trying the Vite 5+ location
// (.vite/manifest.json) before the legacy top-level manifest.json.
func (v *Vite) loadManifest() (map[string]viteChunk, error) {
	v.mu.RLock()
	if v.loaded {
		defer v.mu.RUnlock()
		if v.manifest == nil {
			return nil, fmt.Errorf("manifest not found under %s", filepath.Join(v.publicPath, v.buildDir))
		}
		return v.manifest, nil
	}
	v.mu.RUnlock()

	v.mu.Lock()
	defer v.mu.Unlock()

	if v.loaded {
		return v.manifest, nil
	}
	v.loaded = true

	candidates := []string{
		filepath.Join(v.publicPath, v.buildDir, ".vite", "manifest.json"),
		filepath.Join(v.publicPath, v.buildDir, "manifest.json"),
	}

	for _, candidate := range candidates {
		data, err := os.ReadFile(candidate)
		if err != nil {
			continue
		}

		var manifest map[string]viteChunk
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("parse manifest %s: %w", candidate, err)
		}

		v.manifest = manifest
		return v.manifest, nil
	}

	return nil, fmt.Errorf("manifest not found under %s", filepath.Join(v.publicPath, v.buildDir))
}

func moduleScript(src string) string {
	return fmt.Sprintf("<script type=\"module\" src=\"%s\"></script>", template.HTMLEscapeString(src))
}

func styleLink(href string) string {
	return fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\" />", template.HTMLEscapeString(href))
}

// urlJoin joins URL segments with "/" regardless of OS separator.
func urlJoin(parts ...string) string {
	return strings.Join(parts, "/")
}

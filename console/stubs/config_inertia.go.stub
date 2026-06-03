package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("inertia", map[string]any{
		"root_view":   "resources/inertia/app.gohtml",
		"version":     config.Env("INERTIA_VERSION", ""),
		"ssr":         config.Env("INERTIA_SSR", false),
		"ssr_url":     config.Env("INERTIA_SSR_URL", "http://127.0.0.1:13714/render"),
		"ssr_timeout": config.Env("INERTIA_SSR_TIMEOUT", 5),
		"flash_keys":  []string{"success", "error", "warning", "info", "message"},
		"vite": map[string]any{
			"public_path": "public",
			"build_dir":   "build",
			"hot_file":    "public/hot",
			// Set the dev server URL while developing (e.g. http://localhost:5173)
			// to load assets from Vite; leave empty in production to use the manifest.
			"dev_url": config.Env("VITE_DEV_URL", ""),
		},
	})
}

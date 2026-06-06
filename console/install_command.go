// Package console provides the artisan commands shipped with goravel-inertia.
package console

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

//go:embed stubs/*
var stubs embed.FS

// InstallCommand scaffolds the Inertia.js (Vue 3) frontend, demo pages, root
// template and config into a Goravel application, mirroring an artisan installer.
type InstallCommand struct{}

// NewInstallCommand builds the installer command.
func NewInstallCommand() *InstallCommand {
	return &InstallCommand{}
}

// fileMap maps an embedded stub to its destination path in the application.
// routes/web.go is handled separately (it needs the module name templated in).
var fileMap = map[string]string{
	"stubs/config_inertia.go.stub": "config/inertia.go",
	"stubs/app.gohtml.stub":        "resources/inertia/app.gohtml",
	"stubs/app.ts.stub":            "resources/js/app.ts",
	"stubs/ssr.ts.stub":            "resources/js/ssr.ts",
	"stubs/Layout.vue.stub":        "resources/js/Layout.vue",
	"stubs/Logo.vue.stub":          "resources/js/components/Logo.vue",
	"stubs/goravel-inertia.png":    "resources/js/assets/goravel-inertia.png",
	"stubs/favicon.png":            "public/favicon.png",
	"stubs/Home.vue.stub":          "resources/js/Pages/Home.vue",
	"stubs/Feed.vue.stub":          "resources/js/Pages/Feed.vue",
	"stubs/Contact.vue.stub":       "resources/js/Pages/Contact.vue",
	"stubs/About.vue.stub":         "resources/js/Pages/About.vue",
	"stubs/global.d.ts.stub":       "resources/js/types/global.d.ts",

	"stubs/home_controller.go.stub":    "app/http/controllers/home_controller.go",
	"stubs/feed_controller.go.stub":    "app/http/controllers/feed_controller.go",
	"stubs/contact_controller.go.stub": "app/http/controllers/contact_controller.go",
	"stubs/about_controller.go.stub":   "app/http/controllers/about_controller.go",

	"stubs/vite.config.ts.stub": "vite.config.ts",
	"stubs/tsconfig.json.stub":  "tsconfig.json",
	"stubs/package.json.stub":   "package.json",
}

const (
	webRoutesPath     = "routes/web.go"
	welcomePath       = "resources/views/welcome.tmpl"
	modulePlaceholder = "__MODULE__"
)

// Signature is the unique command name.
func (r *InstallCommand) Signature() string {
	return "inertia:install"
}

// Description is shown in the artisan command list.
func (r *InstallCommand) Description() string {
	return "Scaffold the Inertia.js (Vue 3) frontend, demo pages, root template and config"
}

// Extend declares the command category and flags.
func (r *InstallCommand) Extend() command.Extend {
	return command.Extend{
		Category: "inertia",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:  "force",
				Usage: "Overwrite files that already exist",
			},
		},
	}
}

// Handle scaffolds the stub files, wires routes/web.go to serve the Inertia demo
// (replacing the default welcome view), then prints the remaining manual step.
func (r *InstallCommand) Handle(ctx console.Context) error {
	force := ctx.OptionBool("force")

	created, skipped, err := scaffold(force)
	for _, dst := range created {
		ctx.Info(fmt.Sprintf("  created: %s", dst))
	}
	for _, dst := range skipped {
		ctx.Warning(fmt.Sprintf("  skipped (exists): %s", dst))
	}
	if err != nil {
		return err
	}

	// routes/web.go — needs the consumer's module path for the controllers import.
	module, modErr := moduleName()
	if modErr != nil {
		ctx.Warning(fmt.Sprintf("  skipped %s: %v", webRoutesPath, modErr))
	} else {
		status, werr := installWebRoutes(module, force)
		if werr != nil {
			return werr
		}
		ctx.Info("  " + status)
	}

	// Drop the default Goravel welcome view; "/" now renders the Inertia Home page.
	if removed, rerr := removeWelcome(); rerr != nil {
		return rerr
	} else if removed {
		ctx.Info("  removed: " + welcomePath)
	}

	ctx.NewLine()
	ctx.Success(fmt.Sprintf("Inertia scaffolding complete (%d created, %d skipped).", len(created), len(skipped)))
	r.printNextSteps(ctx)

	return nil
}

// scaffold writes the stub files into the current working directory relative to
// their mapped destinations, skipping existing files unless force is set. It
// returns the created and skipped destination paths. Extracted from Handle so the
// file logic is testable without a console context.
func scaffold(force bool) (created, skipped []string, err error) {
	// Stable order regardless of map iteration.
	srcs := make([]string, 0, len(fileMap))
	for src := range fileMap {
		srcs = append(srcs, src)
	}
	sort.Strings(srcs)

	for _, src := range srcs {
		dst := fileMap[src]

		if _, statErr := os.Stat(dst); statErr == nil && !force {
			skipped = append(skipped, dst)
			continue
		}

		content, readErr := stubs.ReadFile(src)
		if readErr != nil {
			return created, skipped, fmt.Errorf("read stub %s: %w", src, readErr)
		}

		if mkErr := os.MkdirAll(filepath.Dir(dst), 0o755); mkErr != nil {
			return created, skipped, fmt.Errorf("create dir for %s: %w", dst, mkErr)
		}

		if writeErr := os.WriteFile(dst, content, 0o644); writeErr != nil {
			return created, skipped, fmt.Errorf("write %s: %w", dst, writeErr)
		}

		created = append(created, dst)
	}

	return created, skipped, nil
}

// moduleName reads the module path from the application's go.mod.
func moduleName() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("read go.mod: %w", err)
	}

	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if rest, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(rest), nil
		}
	}

	return "", fmt.Errorf("module directive not found in go.mod")
}

// installWebRoutes writes routes/web.go from the stub with the module path filled
// in. It replaces the file when it's missing or still the default Goravel skeleton
// (detected by the welcome.tmpl reference) or when force is set, backing up any
// existing file. A customized web.go is left untouched and the generated version
// is written alongside for manual merging.
func installWebRoutes(module string, force bool) (string, error) {
	stub, err := stubs.ReadFile("stubs/web.go.stub")
	if err != nil {
		return "", fmt.Errorf("read web.go stub: %w", err)
	}
	out := []byte(strings.ReplaceAll(string(stub), modulePlaceholder, module))

	existing, statErr := os.ReadFile(webRoutesPath)
	if statErr != nil {
		if mkErr := os.MkdirAll(filepath.Dir(webRoutesPath), 0o755); mkErr != nil {
			return "", mkErr
		}
		if wErr := os.WriteFile(webRoutesPath, out, 0o644); wErr != nil {
			return "", wErr
		}
		return "created: " + webRoutesPath, nil
	}

	isDefault := strings.Contains(string(existing), "welcome.tmpl")
	if isDefault || force {
		if bErr := os.WriteFile(webRoutesPath+".bak", existing, 0o644); bErr != nil {
			return "", bErr
		}
		if wErr := os.WriteFile(webRoutesPath, out, 0o644); wErr != nil {
			return "", wErr
		}
		return "replaced: " + webRoutesPath + " (backup at " + webRoutesPath + ".bak)", nil
	}

	alt := "routes/web.inertia.go.txt"
	if wErr := os.WriteFile(alt, out, 0o644); wErr != nil {
		return "", wErr
	}
	return "skipped " + webRoutesPath + " (customized); wrote " + alt + " — merge manually", nil
}

// removeWelcome deletes the default Goravel welcome view if present.
func removeWelcome() (bool, error) {
	if _, err := os.Stat(welcomePath); err != nil {
		return false, nil
	}
	if err := os.Remove(welcomePath); err != nil {
		return false, fmt.Errorf("remove %s: %w", welcomePath, err)
	}
	return true, nil
}

// printNextSteps lists what's left after scaffolding. The service provider is
// registered by `package:install`; the route + middleware wiring is done in web.go.
func (r *InstallCommand) printNextSteps(ctx console.Context) {
	ctx.NewLine()
	ctx.Comment("Next steps:")
	ctx.Line("  1. Make sure the service provider is registered (done automatically by")
	ctx.Line("     `./artisan package:install github.com/eddyjj92/goravel-inertia`). If you")
	ctx.Line("     added the package manually, register it in bootstrap/providers.go:")
	ctx.Line(`       goravelinertia "github.com/eddyjj92/goravel-inertia"`)
	ctx.Line("       ...")
	ctx.Line("       &goravelinertia.ServiceProvider{},")
	ctx.NewLine()
	ctx.Line("  2. Install JS deps and start the dev server:")
	ctx.Line("       npm install")
	ctx.Line("       npm run dev   # writes public/hot, then in another shell:")
	ctx.Line("       go run .")
}

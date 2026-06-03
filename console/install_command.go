// Package console provides the artisan commands shipped with goravel-inertia.
package console

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

//go:embed stubs/*
var stubs embed.FS

// InstallCommand scaffolds the Inertia.js (Vue 3) frontend, root template and
// config into a Goravel application, mirroring `php artisan inertia:install`.
type InstallCommand struct{}

// NewInstallCommand builds the installer command.
func NewInstallCommand() *InstallCommand {
	return &InstallCommand{}
}

// fileMap maps an embedded stub to its destination path in the application.
var fileMap = map[string]string{
	"stubs/config_inertia.go.stub": "config/inertia.go",
	"stubs/app.gohtml.stub":        "resources/inertia/app.gohtml",
	"stubs/app.ts.stub":            "resources/js/app.ts",
	"stubs/ssr.ts.stub":            "resources/js/ssr.ts",
	"stubs/Layout.vue.stub":        "resources/js/Layout.vue",
	"stubs/Dashboard.vue.stub":     "resources/js/Pages/Dashboard.vue",
	"stubs/About.vue.stub":         "resources/js/Pages/About.vue",
	"stubs/global.d.ts.stub":       "resources/js/types/global.d.ts",
	"stubs/vite.config.ts.stub":    "vite.config.ts",
	"stubs/tsconfig.json.stub":     "tsconfig.json",
	"stubs/package.json.stub":      "package.json",
}

// Signature is the unique command name.
func (r *InstallCommand) Signature() string {
	return "inertia:install"
}

// Description is shown in the artisan command list.
func (r *InstallCommand) Description() string {
	return "Scaffold the Inertia.js (Vue 3) frontend, root template and config"
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

// Handle writes the stub files into the application, skipping existing files
// unless --force is given, then prints the remaining manual wiring steps.
func (r *InstallCommand) Handle(ctx console.Context) error {
	created, skipped, err := scaffold(ctx.OptionBool("force"))

	for _, dst := range created {
		ctx.Info(fmt.Sprintf("  created: %s", dst))
	}
	for _, dst := range skipped {
		ctx.Warning(fmt.Sprintf("  skipped (exists): %s", dst))
	}
	if err != nil {
		return err
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

// printNextSteps lists the wiring the installer cannot safely automate (editing
// the user's providers list and route file).
func (r *InstallCommand) printNextSteps(ctx console.Context) {
	ctx.NewLine()
	ctx.Comment("Next steps:")
	ctx.Line("  1. Register the service provider in bootstrap/providers.go:")
	ctx.Line(`       inertiaproviders "github.com/eddyjj92/goravel-inertia/providers"`)
	ctx.Line("       ...")
	ctx.Line("       &inertiaproviders.InertiaServiceProvider{},")
	ctx.NewLine()
	ctx.Line("  2. Add the global middleware in routes/web.go (session must come first):")
	ctx.Line(`       sessionmiddleware "github.com/goravel/framework/session/middleware"`)
	ctx.Line(`       inertiamiddleware "github.com/eddyjj92/goravel-inertia/middleware"`)
	ctx.Line("       ...")
	ctx.Line("       facades.Route().GlobalMiddleware(")
	ctx.Line("           sessionmiddleware.StartSession(),")
	ctx.Line("           inertiamiddleware.Inertia(),")
	ctx.Line("       )")
	ctx.NewLine()
	ctx.Line("  3. Install JS deps and start the dev server:")
	ctx.Line("       npm install")
	ctx.Line("       npm run dev   # then: VITE_DEV_URL=http://localhost:5173 go run .")
}

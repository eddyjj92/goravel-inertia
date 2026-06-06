# Roadmap — goravel-inertia

Path to a stable **v1.0.0**. Tracks what's shipped and what's next.

> Detailed historical log of the v0.1.0 build lives in the development notes

---

## Current state

**Latest published:** `v0.2.1` · branches `master` == `develop`.

The Go adapter is feature-complete for a single frontend stack (Vue 3):

| Area | Status |
|------|--------|
| Render (HTML initial + X-Inertia JSON), version check (409) | ✅ |
| Per-request v3 props (Defer/Optional/Always/Merge/DeepMerge/Prepend/Scroll/Once/Prop) | ✅ |
| Shared props — hybrid: `share()` middleware (per-request) + facade `Share`/`ShareFunc` | ✅ |
| Session flash + validation errors bridged into `props.flash` / `props.errors` | ✅ |
| Inertia-aware redirects (303/302) + external `Location` | ✅ |
| Vite integration — HMR dev (`public/hot`) + hashed prod build | ✅ |
| Asset versioning from manifest hash | ✅ |
| SSR + **automatic CSR fallback** when SSR is unreachable | ✅ |
| `inertia:install` artisan command (Vue 3 demo scaffold) | ✅ |
| `HandleInertiaRequests` publishable middleware (Laravel-style) | ✅ |
| `package:install` setup (auto-registers ServiceProvider) | ✅ |
| Tests — core coverage ~89% | ✅ |
| Mocks (`mocks/Inertia.go`) for consumer controller tests | ✅ |
| README / THIRD_PARTY_NOTICES / LICENSE | ✅ |

### Key architectural fact

The Go side (adapter, manager, props, middleware, SSR fallback) is **fully
agnostic to the JS framework** — the Inertia protocol is framework-independent.
Everything frontend-specific lives in `console/stubs/` and is selected by
`inertia:install`. **Adding a new stack = new stubs + installer flag, no Go core
changes.**

---

## Next — v0.3.0: React support

Goal: `inertia:install --stack=react` scaffolds a full React 18 + Inertia demo,
on par with the current Vue 3 scaffold. Vue stays the default.

### Scope: frontend scaffolding + installer only (zero Go core risk)

#### 1. Installer — stack selection
- [ ] `--stack` flag on `inertia:install` (`vue` default, `react`). Validate value.
- [ ] (Optional) interactive prompt when `--stack` omitted.
- [ ] `fileMap` becomes stack-aware: resolve stub set by chosen stack.

#### 2. Stub reorganization
- [ ] Move Vue stubs to `console/stubs/vue/`.
- [ ] Add `console/stubs/react/`.
- [ ] Keep **shared** stubs at `console/stubs/` (stack-independent): `web.go.stub`,
      `*_controller.go.stub`, `config_inertia.go.stub`,
      `handle_inertia_requests.go.stub`, `app.gohtml.stub` (entry path differs →
      may need per-stack), `favicon.png`, brand images.
- [ ] Update `//go:embed` directives for the new tree.

#### 3. React stubs (mirror of the Vue set)
- [ ] `app.tsx` — `createInertiaApp` + `createRoot` (`react-dom/client`),
      `resolvePageComponent` over `./Pages/**/*.tsx`, persistent layout, progress.
- [ ] `ssr.tsx` — `createServer` (`@inertiajs/react/server`) +
      `ReactDOMServer.renderToString`, resolve mirroring `app.tsx`.
- [ ] `Layout.tsx`, `Logo.tsx`.
- [ ] `Pages/{Home,Feed,Contact,About}.tsx` — same demo features:
      Deferred (`<Deferred>`), Merge ("load more"), flash banner, form with
      `useForm` + `props.errors`, active nav link.
- [ ] `global.d.ts` — React `PageProps` augmentation (`@inertiajs/core`).
- [ ] `package.json` — `react`, `react-dom`, `@inertiajs/react`,
      `@vitejs/plugin-react`, `@types/react`, `@types/react-dom`; drop Vue deps.
- [ ] `vite.config.ts` — `@vitejs/plugin-react`, input `resources/js/app.tsx`,
      keep `goravelHot` plugin + `/build` base + dev origin (reuse as-is).
- [ ] `tsconfig.json` — `"jsx": "react-jsx"`, React lib types.
- [ ] `app.gohtml` — entry `{{ vite "resources/js/app.tsx" }}` (root div identical).

#### 4. Tests
- [ ] `install_command_test.go` parametrized per stack: each scaffolds its file
      set in a clean dir; rerun skips; `--force` overwrites.
- [ ] E2E (manual or scripted): React scaffold **compiles** (`tsc`, `vite build`)
      and renders both initial + X-Inertia paths.

#### 5. Docs
- [ ] README: stack table + `--stack` usage; show both Vue and React snippets.
- [ ] `INERTIA.md` (scaffolded): note the chosen stack.

**Gate:** both stacks scaffold from a clean Goravel app, compile, and render
initial + partial reload. Go core unchanged → existing tests stay green.

---

## Path to v1.0.0

| Version | Theme |
|---------|-------|
| **v0.3.0** | React support (above). |
| **v0.4.0** | Stack polish + (optional) Svelte; shared stub abstraction proven by ≥2 stacks. |
| **v0.5.0** | Hardening: automated `setup/` tests (currently manual), CHANGELOG, API surface review. |
| **v1.0.0** | **Stable API commitment.** Multi-stack (Vue + React), full coverage, docs, semver freeze. |

### v1.0.0 exit criteria
- [ ] ≥2 first-class stacks (Vue + React) scaffolding cleanly.
- [ ] `contracts.Inertia` reviewed and frozen (no planned breaking changes).
- [ ] Automated tests for `setup/` (package-install path).
- [ ] `CHANGELOG.md` maintained.
- [ ] Core coverage held >85%.
- [ ] Docs cover every public method and both stacks.

### Parallel track (external, does not block 1.0.0)
- Official Goravel package candidacy (`goravel/inertia`): publish GitHub
  Discussion, await maintainer decision. If accepted → coordinated module-path
  migration (major-version event). Tracked separately; depends on a third party,
  so it is **not** a 1.0.0 gate.

---

## Conventions

- Work on a feature/fix branch → merge to `develop` → fast-forward to `master`.
- Commits authored by the maintainer only (no co-author trailers).
- Each phase ships its own tests; a phase is "done" only when its gate is green.

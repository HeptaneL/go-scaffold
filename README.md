# Go Scaffold

This repository provides a lightweight project scaffolding tool for quickly bootstrapping Go backend services. It renders the files under `templates/` with the values you supply and conditionally includes API, admin, and task modules based on your needs.

## Prerequisites

- Go 1.25 or newer

## Getting Started

```bash
git clone https://github.com/<your-org>/go-scaffold.git
cd go-scaffold
go run ./scaffold.go create my-service \
  --module=github.com/acme/my-service \
  --port=8080 \
  --with=api,admin,task
```

This command generates a new project in `./my-service`. The scaffold embeds all template assets, so there is no need to copy the `templates/` directory manually.

### CLI Flags

- `--module` (required): Fully qualified Go module path for the new project.
- `--port` (optional): Default HTTP port (used inside rendered configuration files). Defaults to `8080`.
- `--with` (optional): Comma-separated list of components to include. Supported values:
  - `api`: REST API service under `cmd/api` and `internal/app/api`.
  - `admin`: Admin HTTP service under `cmd/admin` and `internal/app/admin`.
  - `task`: Background task runner under `cmd/task` and `internal/app/task`.

If a component is omitted, the scaffolder skips its `cmd/`, `internal/app/`, and related router files.

## Project Layout

Generated projects share a common layout:

```
cmd/            # Service entrypoints (api, admin, task)
internal/       # Application code (bootstrap, app modules, router, services)
pkg/            # Reusable libraries (db, cache, logging, integrations)
settings/       # Configuration files
scripts/        # Helper SQL and maintenance scripts
```

All `.tmpl` files originate from `templates/` and are rendered with your project/module details.

## Customising Templates

- Edit files under `templates/` to change scaffold output.
- Run `rename.sh` if you introduce new plain files and want to ensure they carry the `.tmpl` suffix.
- Use `update.sh` to mass-replace module placeholders inside templates.

## Next Steps

After scaffolding a project:

```bash
cd my-service
go mod tidy
go run ./cmd/api/main.go start -c settings/local.json
```

Each generated command provides a Cobra CLI (`start` subcommand) to launch the service with the configuration found in `settings/local.json`.

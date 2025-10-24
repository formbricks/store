# Repository Guidelines

## Project Structure & Module Organization
The repository is a pnpm/turbo monorepo. The Go API lives in `apps/store` with its entrypoint in `cmd/store` and domain code under `internal`. Documentation is generated from `apps/docs`, a Docusaurus site that also hosts OpenAPI outputs stored in `static/openapi`. Shared linting and TypeScript baselines sit in `packages/eslint-config` and `packages/typescript-config`. Docker and CI helper files are at the root (`docker-compose.yml`, `amplify.yml`, `turbo.json`).

## Build, Test, and Development Commands
Run `pnpm install` once to hydrate all workspaces. Use `pnpm dev` to launch the turbo-powered dev pipeline (spawns the Go service and docs watchers where configured). `pnpm build`, `pnpm test`, `pnpm lint`, and `pnpm check-types` fan out to the respective tasks across apps. Inside `apps/store`, use `make setup` to copy `.env`, install Go tools, and start PostgreSQL via Docker; `make dev` runs the API with hot reload; `make build` compiles the binary to `bin/store`; `make generate-openapi` refreshes the docs contract.

## Coding Style & Naming Conventions
Prettier enforces formatting (`pnpm format` for write, `pnpm format:check` for CI). JavaScript/TypeScript code follows the shared ESLint config; use 2-space indentation and PascalCase for React components while keeping file names kebab-case. Go code must stay `gofmt`-clean with package names in lower_snake_case; generated Ent code belongs under `internal/ent`. Keep environment example templates in `env.example` and store secrets only in local `.env`.

## Testing Guidelines
Primary tests are Go unit/integration suites located beside implementations as `*_test.go`. Run `pnpm test` for a repository-wide check or `make test` inside `apps/store` for verbose race-and-cover runs. Generate HTML coverage when needed with `make test-coverage` (outputs `coverage.html`). When touching docs components, at minimum run `pnpm --filter @formbricks/docs lint` and `pnpm --filter @formbricks/docs check-types` to confirm MDX and TS typings.

## Commit & Pull Request Guidelines
The project uses Conventional Commits (e.g., `feat(api): add survey endpoints`), enforced by commitlint. Keep subjects under 72 characters and scope by package or domain (`docs`, `docker`, `store`). Before opening a PR, ensure lint, tests, and type checks pass; include a concise summary, linked issue references, and screenshots for UI or dashboard changes. Note any configuration impacts (e.g., new env vars) and mention follow-up migration steps when required.

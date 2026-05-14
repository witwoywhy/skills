---
name: go-backend-service
description: 0.0.1 | Use when building, reviewing, or modifying Go backend services that follow witwoywhy's clean architecture pattern, use github.com/witwoywhy/go-cores, define internal service/domain/handler/adaptor/repository packages, or need database migration SQL aligned with GORM domain models.
---

# Go Backend Service

Use this skill to build or review Go backend services in the local house style: clean architecture, hexagonal ports, `go-cores` runtime helpers, explicit dependency injection, and migration SQL that matches domain models.

## Usage

When this skill triggers:

1. Inspect the target repo before proposing structure or edits.
2. Identify what the user is asking for: new service, new endpoint, consumer, adaptor, repository, migration, review, or refactor.
3. Load the smallest relevant reference set from the section below.
4. Follow the repo's existing package names and patterns unless they conflict with this skill's core rules.
5. Make the change directly when the user asks for implementation.
6. Explain architectural tradeoffs only when they affect package boundaries, dependency direction, testability, or migration safety.
7. Before finishing, run the narrowest useful verification command available in the repo.

Use this skill especially for prompts like:

- "create a new Go backend service"
- "add an HTTP endpoint"
- "add a Kafka consumer"
- "implement service/repository/adaptor"
- "review this Go service structure"
- "generate migration SQL from this domain model"
- "align GORM model and migration"
- "wire go-cores logging/config/gin/kafka/gorm"

Do not use this skill for generic Go libraries, CLIs, frontend code, or unrelated infrastructure unless the repo clearly follows this backend service pattern.

## Reference Loading

Read only the references needed for the task.

- For service structure, package boundaries, handlers, services, orchestrators, adaptors, repositories, and unit tests: read `references/backend-design-pattern.md`.
- For `go-cores` APIs such as `vipers`, `apps`, `logs`, `gins`, `errs`, `reqs`, `kafka`, `gorms`, `storage`, and utilities: read `references/go-cores.md`.
- For migration repo layout, migration file naming, SQL type mapping, and domain-to-schema checks: read `references/db-migration-pattern.md`.

If the user asks for a new service, start with `backend-design-pattern.md`, then load `go-cores.md` for the runtime packages actually used. Load `db-migration-pattern.md` only when the service uses a database or the request mentions schema, migration, GORM, table, or SQL.

## Core Rules

- `vipers.Init()` runs in `init()` before any config-dependent package reads Viper.
- `main()` initializes only the infrastructure the service needs, then calls exactly one runtime entrypoint.
- Runtime packages start the process. They do not own business logic.
- Handlers and consumers are edge adapters. They bind, parse, validate, construct dependencies, call services, and return responses or ack/nack.
- Services own use case flow and business decisions.
- Orchestrators hold reusable business rules that would otherwise create service import cycles.
- Repositories integrate internal resources or services under the project's control.
- Adaptors integrate external systems outside the project's control.
- Domain and enum packages hold shared business language and must not depend on handler, service, orchestrator, adaptor, or repository packages.
- Use constructor injection. Build concrete dependencies at the edge and pass ports into services.
- Unit tests should assert public `Execute(...)` behavior with mocked dependencies, not private implementation.

## Implementation Workflow

1. Identify the runtime type: HTTP, consumer, scheduler, or worker.
2. Map required infrastructure: app and log are baseline; add validator, DB, circuit breaker, Kafka, storage, or outbound HTTP only when used.
3. Place code by responsibility:
   - `infrastructure/` for process-level singleton initialization.
   - `httpserv/`, `consumer/`, or `scheduler/` for runtime entrypoints.
   - `library/` for generic helpers that do not import `internal/`.
   - `internal/domain/` and `internal/enum/` for shared business types.
   - `internal/handler/` for HTTP edge binding and dependency construction.
   - `internal/service/<action>/` for use cases.
   - `internal/orchestrator/<name>/` for shared business rules.
   - `internal/repository/<name>/` for internal resources.
   - `internal/adaptor/<name>/` for external systems.
4. Define ports in the package that consumes them unless an existing local pattern says otherwise.
5. Keep package APIs small: `domain.go` for public request/response/interfaces, `service.go` for implementation, focused files for helpers or mappers.
6. Add or update tests where business flow, branching, dependency calls, or error mapping changed.
7. When DB models change, verify migrations match `TableName()` and every `gorm:"column:..."` tag.

## HTTP Pattern

For HTTP services:

- Use `httpserv.Run()` as the final call in `main()`.
- Create the Gin app with `gins.New()`.
- Register middleware before routes.
- Bind routes with `handler.Bind<Action>Route(app)`.
- Keep dependency construction inside handler bind functions.
- Use `app.WithLogger(...)` when the service only needs request and logger.
- Use `app.WithRouteContext(...)` only when the service needs HTTP route context.
- Return `errs.Error` from services when errors need standard HTTP mapping.

## Consumer Pattern

For message consumers:

- Use `consumer.Run()` instead of `httpserv.Run()`.
- Create the service before starting the consume loop.
- Configure the consumer group with `kafka.AddConfigKey(...)`.
- Keep the callback thin: unmarshal message, create trace-aware logger, call service, return ack/nack.
- Prefer `Execute(request, logger)` unless the use case truly needs extra context.

## Database And Migration Pattern

When a service uses GORM/domain persistence:

- Domain `TableName()` must match the SQL table name.
- Every persisted field's `gorm:"column:..."` tag must match the SQL column name.
- Store enum values as `VARCHAR`, not database enum types, unless the project already uses a stronger local convention.
- Use explicit ordered migration filenames: `{version}_{action}_{table_name}.up.sql` and `.down.sql`.
- Every `up.sql` must have a matching `down.sql`.
- Each `down.sql` reverses only its matching `up.sql`.
- Use one version for one schema or data action.

## Review Checklist

Before finishing, check:

- The runtime entrypoint matches the service type.
- Config is loaded before packages call `viper.UnmarshalKey(...)`.
- Business logic did not leak into infrastructure, handler, consumer, adaptor, or repository code.
- Dependencies are injected through constructors instead of created deep inside services.
- Repository versus adaptor naming reflects internal versus external ownership.
- Domain and enum packages stay dependency-light.
- Logs use `logger.Logger` and carry trace/span IDs when available.
- Error responses use `errs` and project error mapping conventions.
- Migrations and GORM models agree when persistence is involved.
- Tests cover the changed public behavior.

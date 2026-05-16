# FlowForge Database Schemas

Folder `db/schemas` berisi migration untuk membuat PostgreSQL schema per domain. Folder `db/migrations` dikelompokkan per schema dan berisi migration table sesuai kegunaan schema tersebut.

## Schema Domain

| Schema | Kegunaan | Tabel |
|---|---|---|
| `master` | Data master lintas module | `master.tenants` |
| `auth` | Identity, user, role, dan credential | `auth.users` |
| `workflow` | Definisi workflow dan versioning DAG | `workflow.workflows`, `workflow.workflow_versions` |
| `execution` | Runtime execution, step status, dan log | `execution.workflow_runs`, `execution.run_steps`, `execution.execution_logs` |
| `scheduler` | Jadwal cron workflow | `scheduler.schedules` |

## Struktur Folder Migration

```text
db/
  schemas/
    000001_master.up.sql
    000002_auth.up.sql
    000003_workflow.up.sql
    000004_execution.up.sql
    000005_scheduler.up.sql
  migrations/
    extensions/
      000002_enable_pgcrypto.up.sql
    master/
      000002_create_tenants_table.up.sql
    auth/
      000002_create_users_table.up.sql
    workflow/
      000002_create_workflows_table.up.sql
      000003_create_workflow_versions_table.up.sql
    execution/
      000002_create_workflow_runs_table.up.sql
      000003_create_run_steps_table.up.sql
      000004_create_execution_logs_table.up.sql
    scheduler/
      000002_create_schedules_table.up.sql
```

Table migration dimulai dari `000002` karena `000001` dipakai untuk pembuatan schema di `db/schemas`.

## Migration Order

Jalankan schema migration sebelum table migration:

```bash
./bin/migrate.sh up
```

Script tersebut menjalankan folder dengan urutan:

1. `db/schemas`
2. `db/migrations/extensions`
3. `db/migrations/master`
4. `db/migrations/auth`
5. `db/migrations/workflow`
6. `db/migrations/execution`
7. `db/migrations/scheduler`

Rollback menjalankan urutan sebaliknya. Setiap folder memakai migration tracking table sendiri agar numbering per schema tidak bentrok.

## Entity Mapping

| Entity | TableName |
|---|---|
| `Tenant` | `master.tenants` |
| `User` | `auth.users` |
| `Workflow` | `workflow.workflows` |
| `WorkflowVersion` | `workflow.workflow_versions` |
| `WorkflowRun` | `execution.workflow_runs` |
| `RunStep` | `execution.run_steps` |
| `ExecutionLog` | `execution.execution_logs` |
| `Schedule` | `scheduler.schedules` |

## Relationship Ringkas

- `auth.users.tenant_id` → `master.tenants.id`
- `workflow.workflows.tenant_id` → `master.tenants.id`
- `workflow.workflows.created_by` → `auth.users.id`
- `workflow.workflow_versions.workflow_id` → `workflow.workflows.id`
- `workflow.workflow_versions.tenant_id` → `master.tenants.id`
- `workflow.workflow_versions.created_by` → `auth.users.id`
- `execution.workflow_runs.workflow_id` → `workflow.workflows.id`
- `execution.workflow_runs.workflow_version_id` → `workflow.workflow_versions.id`
- `execution.workflow_runs.tenant_id` → `master.tenants.id`
- `execution.run_steps.run_id` → `execution.workflow_runs.id`
- `execution.run_steps.tenant_id` → `master.tenants.id`
- `execution.execution_logs.run_id` → `execution.workflow_runs.id`
- `execution.execution_logs.step_row_id` → `execution.run_steps.id`
- `execution.execution_logs.tenant_id` → `master.tenants.id`
- `scheduler.schedules.workflow_id` → `workflow.workflows.id`
- `scheduler.schedules.tenant_id` → `master.tenants.id`
- `scheduler.schedules.created_by` → `auth.users.id`

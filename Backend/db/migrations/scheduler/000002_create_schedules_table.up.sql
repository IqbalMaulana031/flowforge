CREATE TABLE IF NOT EXISTS scheduler.schedules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  workflow_id UUID NOT NULL REFERENCES workflow.workflows(id),
  tenant_id UUID NOT NULL REFERENCES master.tenants(id),
  cron_expression VARCHAR(120) NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  last_run_at TIMESTAMPTZ,
  next_run_at TIMESTAMPTZ,
  created_by UUID REFERENCES auth.users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_schedules_workflow_id ON scheduler.schedules(workflow_id);
CREATE INDEX IF NOT EXISTS idx_schedules_tenant_id ON scheduler.schedules(tenant_id);
CREATE INDEX IF NOT EXISTS idx_schedules_tenant_active ON scheduler.schedules(tenant_id, is_active);

CREATE TABLE IF NOT EXISTS execution.workflow_runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  workflow_id UUID NOT NULL REFERENCES workflow.workflows(id),
  tenant_id UUID NOT NULL REFERENCES master.tenants(id),
  workflow_version_id UUID NOT NULL REFERENCES workflow.workflow_versions(id),
  trigger_type VARCHAR(30) NOT NULL,
  status VARCHAR(30) NOT NULL,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  duration_ms BIGINT,
  error_message TEXT,
  trigger_payload JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_workflow_runs_workflow_id ON execution.workflow_runs(workflow_id);
CREATE INDEX IF NOT EXISTS idx_workflow_runs_tenant_id ON execution.workflow_runs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_workflow_runs_workflow_version_id ON execution.workflow_runs(workflow_version_id);
CREATE INDEX IF NOT EXISTS idx_workflow_runs_status ON execution.workflow_runs(status);
CREATE INDEX IF NOT EXISTS idx_workflow_runs_tenant_created ON execution.workflow_runs(tenant_id, created_at DESC) INCLUDE (status, workflow_id, trigger_type, duration_ms);
CREATE INDEX IF NOT EXISTS idx_workflow_runs_workflow_status ON execution.workflow_runs(workflow_id, status);

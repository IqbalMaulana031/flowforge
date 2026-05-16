CREATE TABLE IF NOT EXISTS execution.run_steps (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id UUID NOT NULL REFERENCES execution.workflow_runs(id),
  tenant_id UUID NOT NULL REFERENCES master.tenants(id),
  step_id VARCHAR(120) NOT NULL,
  step_name VARCHAR(200) NOT NULL,
  step_type VARCHAR(50) NOT NULL,
  status VARCHAR(30) NOT NULL,
  attempt_number INTEGER NOT NULL DEFAULT 0,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  duration_ms BIGINT,
  error_message TEXT,
  output JSONB
);

CREATE INDEX IF NOT EXISTS idx_run_steps_run_id ON execution.run_steps(run_id);
CREATE INDEX IF NOT EXISTS idx_run_steps_tenant_id ON execution.run_steps(tenant_id);
CREATE INDEX IF NOT EXISTS idx_run_steps_status ON execution.run_steps(status);
CREATE INDEX IF NOT EXISTS idx_run_steps_tenant_run ON execution.run_steps(tenant_id, run_id);

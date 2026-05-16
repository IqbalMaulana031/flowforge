CREATE TABLE IF NOT EXISTS execution.execution_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id UUID NOT NULL REFERENCES execution.workflow_runs(id),
  step_row_id UUID REFERENCES execution.run_steps(id),
  tenant_id UUID NOT NULL REFERENCES master.tenants(id),
  level VARCHAR(20) NOT NULL,
  message TEXT NOT NULL,
  metadata JSONB,
  logged_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_execution_logs_run_id ON execution.execution_logs(run_id);
CREATE INDEX IF NOT EXISTS idx_execution_logs_step_row_id ON execution.execution_logs(step_row_id);
CREATE INDEX IF NOT EXISTS idx_execution_logs_tenant_id ON execution.execution_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_execution_logs_logged_at ON execution.execution_logs(logged_at);
CREATE INDEX IF NOT EXISTS idx_execution_logs_tenant_run_logged ON execution.execution_logs(tenant_id, run_id, logged_at DESC);

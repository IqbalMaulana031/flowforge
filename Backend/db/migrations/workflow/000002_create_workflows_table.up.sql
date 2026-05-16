CREATE TABLE IF NOT EXISTS workflow.workflows (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES master.tenants(id),
  name VARCHAR(200) NOT NULL,
  description TEXT,
  current_version INTEGER NOT NULL DEFAULT 1,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_by UUID REFERENCES auth.users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_workflows_tenant_id ON workflow.workflows(tenant_id);
CREATE INDEX IF NOT EXISTS idx_workflows_tenant_created ON workflow.workflows(tenant_id, created_at DESC);
CREATE UNIQUE INDEX IF NOT EXISTS idx_workflows_tenant_name_active ON workflow.workflows(tenant_id, name) WHERE is_active = TRUE;

CREATE TABLE IF NOT EXISTS workflow.workflow_versions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  workflow_id UUID NOT NULL REFERENCES workflow.workflows(id),
  tenant_id UUID NOT NULL REFERENCES master.tenants(id),
  version_number INTEGER NOT NULL,
  dag_definition JSONB NOT NULL,
  changelog TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_by UUID REFERENCES auth.users(id),
  UNIQUE(workflow_id, version_number)
);

CREATE INDEX IF NOT EXISTS idx_workflow_versions_workflow_id ON workflow.workflow_versions(workflow_id);
CREATE INDEX IF NOT EXISTS idx_workflow_versions_tenant_id ON workflow.workflow_versions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_workflow_versions_tenant_workflow ON workflow.workflow_versions(tenant_id, workflow_id, version_number DESC);

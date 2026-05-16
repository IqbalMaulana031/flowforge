INSERT INTO master.tenants (name, slug, status)
VALUES ('Demo Tenant', 'demo', 'active')
ON CONFLICT (slug) DO NOTHING;

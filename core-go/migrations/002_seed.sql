-- +migrate Up
INSERT INTO plans (code, name, limits, is_default)
VALUES (
  'starter',
  'Starter',
  '{"max_sessions": 3, "max_users": 5, "max_messages_month": 10000}'::jsonb,
  TRUE
)
ON CONFLICT (code) DO NOTHING;

-- senha: ChangeMe123!
INSERT INTO users (email, password_hash, full_name, status)
VALUES (
  'admin@example.com',
  '$2a$10$4v0hHBM8VHbGSLf6M5Z8n.OYjpxwP4rM0WjHgfbx3Rp3XIanFkFBS',
  'Admin User',
  'active'
)
ON CONFLICT (email) DO NOTHING;

WITH starter_plan AS (
  SELECT id FROM plans WHERE code = 'starter' LIMIT 1
), tenant_row AS (
  INSERT INTO tenants (name, slug, status, plan_id)
  SELECT 'Tenant Demo', 'tenant-demo', 'active', starter_plan.id FROM starter_plan
  ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
  RETURNING id
), user_row AS (
  SELECT id FROM users WHERE email = 'admin@example.com' LIMIT 1
)
INSERT INTO tenant_users (tenant_id, user_id, role)
SELECT tenant_row.id, user_row.id, 'owner'
FROM tenant_row, user_row
ON CONFLICT (tenant_id, user_id) DO NOTHING;

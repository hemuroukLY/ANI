INSERT INTO users (
    tenant_id, username, email, password_hash,
    display_name, status, is_deleted, created_at, updated_at
)
VALUES
    -- 平台 root
    (
        NULL,
        'local:root',
        'root@platform.local',
        '$2a$12$4.kOj8KTyoFLnJpJFbN3zeJMAZiXIympHJ2MJSJfq0TLNXC6gLwW2',
        'Platform Root',
        'active', FALSE, NOW(), NOW()
    ),
    -- 租户 admin（tenant-admin）
    (
        '00000000-0000-0000-0000-000000000001',
        'local:admin',
        'admin@tenant-a.local',
        '$2a$12$4.kOj8KTyoFLnJpJFbN3zeJMAZiXIympHJ2MJSJfq0TLNXC6gLwW2',
        'Tenant Admin',
        'active', FALSE, NOW(), NOW()
    ),
    -- 新租户普通用户 demo（user）
    (
        '00000000-0000-0000-0000-000000000001',
        'local:demo',
        'demo@tenant-a.local',
        '$2a$12$4.kOj8KTyoFLnJpJFbN3zeJMAZiXIympHJ2MJSJfq0TLNXC6gLwW2',
        'Demo User',
        'active', FALSE, NOW(), NOW()
    )
ON CONFLICT DO NOTHING;

-- ===========================================================================
-- user_roles 关联：为上述用户分配内置角色
-- roles 来源：migration 20260502_003 内置角色
--   00000000-0000-0000-0000-000000000001  platform-admin (scope: platform, "*":"*")
--   00000000-0000-0000-0000-000000000002  tenant-admin   (scope: tenant,  "*":"*")
--   00000000-0000-0000-0000-000000000003  user           (普通用户)
-- 单角色约束：migration 20260827_001 在 user_roles(user_id) 上建立唯一索引
-- ===========================================================================
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r
    ON r.tenant_id IS NULL
    AND (
        (u.username = 'local:root'  AND r.name = 'platform-admin')
        OR (u.username = 'local:admin' AND r.name = 'tenant-admin')
        OR (u.username = 'local:demo'  AND r.name = 'user')
    )
WHERE u.tenant_id IS NULL
    OR u.tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (user_id, role_id) DO NOTHING;


    
-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
    old_role_id BIGINT;
    new_role_id BIGINT;
BEGIN
    SELECT id INTO old_role_id FROM roles WHERE name = 'customer' LIMIT 1;
    SELECT id INTO new_role_id FROM roles WHERE name = 'user' LIMIT 1;

    IF old_role_id IS NOT NULL AND new_role_id IS NULL THEN
        UPDATE roles SET name = 'user', title = 'کاربر عادی', updated_at = NOW() WHERE id = old_role_id;
    ELSIF old_role_id IS NOT NULL AND new_role_id IS NOT NULL THEN
        INSERT INTO user_roles (user_id, role_id)
        SELECT user_id, new_role_id FROM user_roles WHERE role_id = old_role_id
        ON CONFLICT (user_id, role_id) DO NOTHING;

        INSERT INTO role_permissions (role_id, permission_id)
        SELECT new_role_id, permission_id FROM role_permissions WHERE role_id = old_role_id
        ON CONFLICT (role_id, permission_id) DO NOTHING;

        DELETE FROM roles WHERE id = old_role_id;
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
UPDATE roles SET name = 'customer', title = 'مشتری / کاربر عادی', updated_at = NOW()
WHERE name = 'user' AND NOT EXISTS (SELECT 1 FROM roles WHERE name = 'customer');

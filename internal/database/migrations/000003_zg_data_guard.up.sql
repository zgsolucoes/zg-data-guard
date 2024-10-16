CREATE INDEX IF NOT EXISTS idx_access_permission_log_database_instance_id
	ON access_permission_log (database_instance_id);

CREATE INDEX IF NOT EXISTS idx_access_permission_log_user_id
	ON access_permission_log (user_id);

CREATE INDEX IF NOT EXISTS idx_access_permission_log_database_user_id
	ON access_permission_log (database_user_id);

CREATE INDEX IF NOT EXISTS idx_access_permission_log_database_id
	ON access_permission_log (database_id);

CREATE INDEX IF NOT EXISTS idx_database_instances_ecosystem_id
	ON database_instances (ecosystem_id);

CREATE INDEX IF NOT EXISTS idx_database_instances_database_technology_id
	ON database_instances (database_technology_id);

CREATE INDEX IF NOT EXISTS idx_databases_database_instance_id
	ON databases (database_instance_id);

CREATE INDEX IF NOT EXISTS idx_databases_created_by_user_id
	ON databases (created_by_user_id);

CREATE INDEX IF NOT EXISTS idx_database_users_database_role_id
	ON database_users (database_role_id);

CREATE INDEX IF NOT EXISTS idx_database_users_created_by_user_id
	ON database_users (created_by_user_id);

CREATE INDEX IF NOT EXISTS idx_access_permissions_database_id
	ON access_permissions (database_id);

CREATE INDEX IF NOT EXISTS idx_access_permissions_database_user_id
	ON access_permissions (database_user_id);

CREATE INDEX IF NOT EXISTS idx_access_permissions_granted_by_user_id
	ON access_permissions (granted_by_user_id);

CREATE INDEX IF NOT EXISTS idx_access_permissions_database_role_id
	ON access_permissions (database_role_id);

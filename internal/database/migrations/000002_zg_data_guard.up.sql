INSERT INTO application_users (name, email) VALUES ('zg-service', 'zg-service@email.com');

INSERT INTO database_roles (name, display_name, description, created_by_user_id, read_only)
VALUES ('user_ro', 'Read Only', 'User with read-only permission', (SELECT id FROM application_users WHERE email = 'zg-service@email.com'), TRUE),
	   ('developer', 'Developer', 'User with permission to manipulate data from tables', (SELECT id FROM application_users WHERE email = 'zg-service@email.com'), FALSE),
	   ('devops', 'DevOps', 'User with permission to manipulate data from tables and change the database structure',
		(SELECT id FROM application_users WHERE email = 'zg-service@email.com'), FALSE),
	   ('application', 'Application', 'User with permission to manipulate data from tables and change the database structure, intended only for service use',
		(SELECT id FROM application_users WHERE email = 'zg-service@email.com'), FALSE);

INSERT INTO forbidden_databases (database_name, description, created_at, created_by_user_id, updated_at)
VALUES ('zg-data-guard', 'Database used by the ZG Data Guard service, should not be used or accessed by any other service or user',
        NOW(), (SELECT id FROM application_users WHERE email = 'zg-service@email.com'), NOW());

ALTER TABLE database_instances
	ADD COLUMN IF NOT EXISTS roles_created BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE databases
	ADD COLUMN IF NOT EXISTS roles_configured BOOLEAN NOT NULL DEFAULT FALSE;

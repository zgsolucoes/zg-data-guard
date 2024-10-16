DELETE
FROM database_roles
WHERE name in ('user_ro', 'developer', 'devops', 'application');

DELETE
FROM forbidden_databases
WHERE database_name = 'zg-data-guard';

ALTER TABLE database_instances
	DROP COLUMN IF EXISTS roles_created;

ALTER TABLE databases
	DROP COLUMN IF EXISTS roles_configured;

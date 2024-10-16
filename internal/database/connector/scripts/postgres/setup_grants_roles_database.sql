-- Description: Script to set up grants to the roles in the current database. For further information, please check the README.
DO
$$
	DECLARE
		role_user_ro        text := 'user_ro';
		role_developer      text := 'developer';
		role_devops         text := 'devops';
		role_application    text := 'application';
		current_schema_name text;

	BEGIN
		EXECUTE FORMAT('REVOKE CONNECT ON DATABASE %I FROM PUBLIC', CURRENT_DATABASE());
		EXECUTE FORMAT('GRANT CREATE, TEMP ON DATABASE %I TO %I', CURRENT_DATABASE(), role_application);

		FOR current_schema_name IN (SELECT schema_name
									FROM information_schema.schemata
									WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast', 'pg_temp', 'pg_statistic')
									  AND schema_name NOT LIKE 'pg_temp_%'
									  AND schema_name NOT LIKE 'pg_toast_%')
			LOOP
				EXECUTE FORMAT('REVOKE CREATE ON SCHEMA %I FROM PUBLIC', current_schema_name, CURRENT_DATABASE());

				-- ROLE USER_RO
				EXECUTE FORMAT('GRANT USAGE ON SCHEMA %I TO %I', current_schema_name, role_user_ro);
				EXECUTE FORMAT('GRANT SELECT ON ALL TABLES IN SCHEMA %I TO %I', current_schema_name, role_user_ro);
				EXECUTE FORMAT('GRANT SELECT ON ALL SEQUENCES IN SCHEMA %I TO %I', current_schema_name, role_user_ro);
				EXECUTE FORMAT('GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA %I TO %I', current_schema_name, role_user_ro);

				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT SELECT ON TABLES TO %I', current_schema_name, role_user_ro);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT SELECT ON SEQUENCES TO %I', current_schema_name, role_user_ro);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT EXECUTE ON FUNCTIONS TO %I', current_schema_name, role_user_ro);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT USAGE ON TYPES TO %I', current_schema_name, role_user_ro);

				-- ROLE DEVELOPER
				EXECUTE FORMAT('GRANT USAGE ON SCHEMA %I TO %I', current_schema_name, role_developer);
				EXECUTE FORMAT('GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA %I TO %I', current_schema_name, role_developer);
				EXECUTE FORMAT('GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA %I TO %I', current_schema_name, role_developer);
				EXECUTE FORMAT('GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA %I TO %I', current_schema_name, role_developer);

				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO %I', current_schema_name, role_developer);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT USAGE, SELECT ON SEQUENCES TO %I', current_schema_name, role_developer);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT EXECUTE ON FUNCTIONS TO %I', current_schema_name, role_developer);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT USAGE ON TYPES TO %I', current_schema_name, role_developer);

				-- ROLE DEVOPS
				EXECUTE FORMAT('GRANT USAGE, CREATE ON SCHEMA %I TO %I', current_schema_name, role_devops);
				EXECUTE FORMAT('GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE, REFERENCES, TRIGGER ON ALL TABLES IN SCHEMA %I TO %I', current_schema_name, role_devops);
				EXECUTE FORMAT('GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA %I TO %I', current_schema_name, role_devops);
				EXECUTE FORMAT('GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA %I TO %I', current_schema_name, role_devops);

				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE, REFERENCES, TRIGGER ON TABLES TO %I', current_schema_name, role_devops);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO %I', current_schema_name, role_devops);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT EXECUTE ON FUNCTIONS TO %I', current_schema_name, role_devops);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT USAGE ON TYPES TO %I', current_schema_name, role_devops);

				-- ROLE APPLICATION
				EXECUTE FORMAT('GRANT USAGE, CREATE ON SCHEMA %I TO %I', current_schema_name, role_application);
				EXECUTE FORMAT('GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE, REFERENCES, TRIGGER ON ALL TABLES IN SCHEMA %I TO %I', current_schema_name, role_application);
				EXECUTE FORMAT('GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA %I TO %I', current_schema_name, role_application);
				EXECUTE FORMAT('GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA %I TO %I', current_schema_name, role_application);

				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT SELECT, INSERT, UPDATE, DELETE, TRUNCATE, REFERENCES, TRIGGER ON TABLES TO %I', current_schema_name,
							   role_application);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO %I', current_schema_name, role_application);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT EXECUTE ON FUNCTIONS TO %I', current_schema_name, role_application);
				EXECUTE FORMAT('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT USAGE ON TYPES TO %I', current_schema_name, role_application);
			END LOOP;
	END
$$;

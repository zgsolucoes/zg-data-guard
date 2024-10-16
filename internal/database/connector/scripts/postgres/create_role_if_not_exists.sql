DO
$$
	DECLARE
		role_name text := '%s';
	BEGIN
		IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = role_name) THEN
			EXECUTE FORMAT('CREATE ROLE %%I', role_name);
		END IF;
	END
$$;

DO
$$
	DECLARE
		user_to_drop text := '%s';
		database_rec record;
	BEGIN

		FOR database_rec IN (SELECT datname AS name
							 FROM pg_database
							 WHERE datistemplate = FALSE
							   AND datname != 'postgres'
							   AND pg_catalog.has_database_privilege((SELECT oid FROM pg_roles WHERE rolname = user_to_drop), oid, 'CONNECT')
							 ORDER BY datname)
			LOOP
				EXECUTE FORMAT('REVOKE ALL PRIVILEGES ON DATABASE %%I FROM %%I', database_rec.name, user_to_drop);
			END LOOP;

		EXECUTE FORMAT('DROP USER IF EXISTS %%I', user_to_drop);
	END
$$;

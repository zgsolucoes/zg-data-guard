SELECT db.id,
	   db.name,
	   db.current_size,
	   di.id,
	   di.name,
	   e.id,
	   e.display_name,
	   dt.id,
	   dt.name,
	   dt.version,
	   db.enabled,
	   db.roles_configured,
	   db.description,
	   db.created_by_user_id,
	   au.name,
	   db.created_at,
	   db.updated_at,
	   di.last_database_sync,
	   db.disabled_at
FROM databases db
	 JOIN database_instances di
		ON db.database_instance_id = di.id
	 JOIN application_users au
		ON db.created_by_user_id = au.id
	 JOIN database_technologies dt
		ON di.database_technology_id = dt.id
	 JOIN ecosystems e
		ON di.ecosystem_id = e.id
WHERE db.id = $1

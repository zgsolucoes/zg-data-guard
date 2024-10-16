SELECT di.id,
	   di.name,
	   hci.host,
	   hci.port,
	   hci.host_connection,
	   hci.port_connection,
	   hci.admin_username,
	   hci.admin_password,
	   di.ecosystem_id,
	   e.display_name,
	   di.database_technology_id,
	   dt.name,
	   dt.version,
	   di.enabled,
	   di.roles_created,
	   di.note,
	   di.created_at,
	   di.created_by_user_id,
	   au.name,
	   di.updated_at,
	   di.disabled_at,
	   di.last_database_sync,
	   di.connection_status,
	   di.last_connection_test,
	   di.last_connection_result
FROM database_instances di
	 JOIN host_connection_info hci
		ON di.host_connection_info_id = hci.id
	 JOIN application_users au
		ON di.created_by_user_id = au.id
	 JOIN database_technologies dt
		ON di.database_technology_id = dt.id
	 JOIN ecosystems e
		ON di.ecosystem_id = e.id
WHERE di.id = $1

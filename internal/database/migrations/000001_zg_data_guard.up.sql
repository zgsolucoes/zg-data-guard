CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the database tables and unique indexes for the application
CREATE TABLE application_users
(
	id          uuid               DEFAULT uuid_generate_v4() PRIMARY KEY,
	name        TEXT      NOT NULL,
	email       TEXT      NOT NULL UNIQUE,
	enabled     BOOLEAN   NOT NULL DEFAULT TRUE,
	created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	disabled_at TIMESTAMP
);

CREATE TABLE ecosystems
(
	id                 uuid               DEFAULT uuid_generate_v4() PRIMARY KEY,
	code               TEXT      NOT NULL UNIQUE,
	display_name       TEXT      NOT NULL,
	created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by_user_id uuid      NOT NULL,
	FOREIGN KEY (created_by_user_id) REFERENCES application_users (id)
);

CREATE TABLE database_technologies
(
	id                 uuid               DEFAULT uuid_generate_v4() PRIMARY KEY,
	name               TEXT      NOT NULL,
	version            TEXT      NOT NULL,
	created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by_user_id uuid      NOT NULL,
	FOREIGN KEY (created_by_user_id) REFERENCES application_users (id)
);
CREATE UNIQUE INDEX idx_unique_name_version_tecnologies
	ON database_technologies (name, version);

CREATE TABLE host_connection_info
(
	id              uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
	host            TEXT NOT NULL,
	port            TEXT NOT NULL,
	host_connection TEXT NOT NULL,
	port_connection TEXT NOT NULL,
	admin_username  TEXT NOT NULL,
	admin_password  TEXT NOT NULL
);
CREATE UNIQUE INDEX idx_unique_host_port_connection
	ON host_connection_info (host, port);

CREATE TABLE database_instances
(
	id                      uuid               DEFAULT uuid_generate_v4() PRIMARY KEY,
	ecosystem_id            uuid      NOT NULL,
	name                    TEXT      NOT NULL,
	host_connection_info_id uuid      NOT NULL UNIQUE,
	database_technology_id   uuid      NOT NULL,
	created_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by_user_id      uuid      NOT NULL,
	enabled                 BOOLEAN   NOT NULL DEFAULT TRUE,
	note                    TEXT,
	connection_status       TEXT      NOT NULL DEFAULT 'NOT_TESTED',
	disabled_at             TIMESTAMP,
	last_connection_result  TEXT,
	last_connection_test    TIMESTAMP,
	last_database_sync      TIMESTAMP,
	FOREIGN KEY (ecosystem_id) REFERENCES ecosystems (id),
	FOREIGN KEY (host_connection_info_id) REFERENCES host_connection_info (id),
	FOREIGN KEY (database_technology_id) REFERENCES database_technologies (id),
	FOREIGN KEY (created_by_user_id) REFERENCES application_users (id)
);

CREATE TABLE databases
(
	id                   uuid               DEFAULT uuid_generate_v4() PRIMARY KEY,
	name                 TEXT      NOT NULL,
	description          TEXT,
	current_size         TEXT,
	database_instance_id uuid      NOT NULL,
	created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	enabled              BOOLEAN   NOT NULL DEFAULT TRUE,
	created_by_user_id   uuid      NOT NULL,
	disabled_at          TIMESTAMP,
	FOREIGN KEY (database_instance_id) REFERENCES database_instances (id),
	FOREIGN KEY (created_by_user_id) REFERENCES application_users (id)
);
CREATE UNIQUE INDEX idx_unique_name_database_instance_id_databases
	ON databases (name, database_instance_id);

CREATE TABLE database_roles
(
	id                 uuid               DEFAULT uuid_generate_v4() PRIMARY KEY,
	name               TEXT      NOT NULL UNIQUE,
	display_name       TEXT      NOT NULL,
	description        TEXT,
	created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by_user_id uuid      NOT NULL,
	read_only          BOOLEAN   NOT NULL DEFAULT TRUE,
	FOREIGN KEY (created_by_user_id) REFERENCES application_users (id)
);

CREATE TABLE database_users
(
	id                 uuid               DEFAULT uuid_generate_v4() PRIMARY KEY,
	name               TEXT      NOT NULL,
	email              TEXT      NOT NULL UNIQUE,
	username           TEXT      NOT NULL UNIQUE,
	password           TEXT,
	position           TEXT,
	team               TEXT,
	enabled            BOOLEAN   NOT NULL DEFAULT TRUE,
	expired            BOOLEAN   NOT NULL DEFAULT FALSE,
	database_role_id   uuid      NOT NULL,
	created_by_user_id uuid      NOT NULL,
	created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at         TIMESTAMP,
	disabled_at        TIMESTAMP,
	FOREIGN KEY (database_role_id) REFERENCES database_roles (id),
	FOREIGN KEY (created_by_user_id) REFERENCES application_users (id)
);

CREATE TABLE access_permissions
(
	id                 uuid               DEFAULT uuid_generate_v4() PRIMARY KEY,
	database_id        uuid      NOT NULL,
	database_role_id   uuid,
	database_user_id   uuid      NOT NULL,
	granted_by_user_id uuid      NOT NULL,
	granted_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (database_id) REFERENCES databases (id),
	FOREIGN KEY (database_role_id) REFERENCES database_roles (id),
	FOREIGN KEY (database_user_id) REFERENCES database_users (id),
	FOREIGN KEY (granted_by_user_id) REFERENCES application_users (id)
);
CREATE UNIQUE INDEX idx_unique_database_user_permissions
	ON access_permissions (database_id, database_user_id);

CREATE TABLE access_permission_log
(
	id                   uuid               DEFAULT uuid_generate_v4() PRIMARY KEY,
	database_instance_id uuid      NOT NULL,
	database_id          uuid,
	database_user_id     uuid,
	message              TEXT      NOT NULL,
	success              BOOLEAN   NOT NULL,
	date                 TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id              uuid      NOT NULL,
	FOREIGN KEY (database_user_id) REFERENCES database_users (id),
	FOREIGN KEY (database_instance_id) REFERENCES database_instances (id),
	FOREIGN KEY (database_id) REFERENCES databases (id),
	FOREIGN KEY (user_id) REFERENCES application_users (id)
);

CREATE TABLE forbidden_databases
(
	id                 uuid               DEFAULT uuid_generate_v4() PRIMARY KEY,
	database_name      TEXT      NOT NULL,
	description        TEXT,
	created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by_user_id uuid      NOT NULL,
	updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (created_by_user_id) REFERENCES application_users (id)
);

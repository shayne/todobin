CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE todos (
    id uuid DEFAULT uuid_generate_v4 (),
	list_id uuid,
	todo VARCHAR(255) NOT NULL,
	done BOOLEAN,
	created_at TIMESTAMP DEFAULT current_timestamp,
	updated_at TIMESTAMP DEFAULT current_timestamp,
    PRIMARY KEY (id)
);


CREATE TABLE lists (
	id uuid DEFAULT uuid_generate_v4 (),
	name VARCHAR(255) NOT NULL,
	created_at TIMESTAMP DEFAULT current_timestamp,
	updated_at TIMESTAMP DEFAULT current_timestamp,
	PRIMARY KEY (id)
);

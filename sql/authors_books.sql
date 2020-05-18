CREATE TABLE IF NOT EXISTS users (
	id varchar(36) NOT NULL,
	email varchar(256) NOT NULL,
	first_name varchar(64) NOT NULL,
	last_name varchar(64) NOT NULL,
	password varchar(256) NOT NULL,
	enabled boolean DEFAULT true NOT NULL,
	created_at timestamp NOT NULL DEFAULT NOW(),
	updated_at timestamp NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id),
	UNIQUE(email)
);

CREATE TABLE IF NOT EXISTS authors (
	id varchar(36) NOT NULL,
	first_name varchar(64) NOT NULL,
	middle_name varchar(64) NOT NULL,
	last_name varchar(64) NOT NULL,
	dob timestamp NOT NULL,
	updated_at timestamp NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id),
	UNIQUE(first_name, middle_name, last_name, dob)
);

CREATE TABLE IF NOT EXISTS books (
	id varchar(36) NOT NULL,
	title varchar(200) NOT NULL,
	isbn varchar(18) NOT NULL,
	updated_at timestamp NOT NULL DEFAULT NOW(),
	PRIMARY KEY(id),
	UNIQUE(isbn)
);

CREATE TABLE IF NOT EXISTS books_authors (
	book_id varchar(36) NOT NULL REFERENCES books (id) ON DELETE CASCADE,
	author_id varchar(36) NOT NULL REFERENCES authors (id) ON DELETE CASCADE,
	PRIMARY KEY(book_id, author_id)
);



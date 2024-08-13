CREATE TABLE users (
	id serial primary key,
	email varchar(255) not null unique,
	password_hash varchar(255) not null
);

CREATE TABLE refresh_tokens(
	id serial primary key,
	user_id integer,
	foreign key (user_id) references users (id) on delete cascade,
	token varchar not null,
    user_ip varchar(100) not null,
	expires_at timestamp
);
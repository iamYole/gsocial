CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXIStS users (
    ID        bigserial PRIMARY KEY,
    username varchar(255) UNIQUE NOT NULL,
	Email     citext UNIQUE NOT NULL, 
	Password  bytea NOT NULL, 
	created_at timestamp(0) with time zone NOT NULL Default NOW()
);
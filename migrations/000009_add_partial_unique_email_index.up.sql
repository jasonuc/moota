UPDATE users SET email = NULL WHERE email = '';

ALTER TABLE users DROP CONSTRAINT users_email_key;

CREATE UNIQUE INDEX users_email_key
	ON users (email)
	WHERE email IS NOT NULL AND email <> '';

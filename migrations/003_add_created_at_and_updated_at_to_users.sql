ALTER TABLE users ADD COLUMN created_at timestamp DEFAULT now() NOT NULL;
ALTER TABLE users ADD COLUMN updated_at timestamp DEFAULT now() NOT NULL;
---- create above / drop below ----
ALTER TABLE users DROP COLUMN created_at
ALTER TABLE users DROP COLUMN updated_at;

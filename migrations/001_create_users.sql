CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users(
  id uuid not null primary key,
  first_name varchar not null,
  last_name varchar not null
);

---- create above / drop below ----
DROP TABLE users;

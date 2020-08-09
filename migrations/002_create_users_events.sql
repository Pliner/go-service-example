CREATE TABLE users_events(
    id serial not null primary key,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    user_id uuid not null
);
---- create above / drop below ----
DROP TABLE users_events;

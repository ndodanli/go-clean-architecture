drop table if exists endpoint;
create table if not exists endpoint
(
    id          bigserial primary key,
    name        text      not null,
    description text,

    created_at  timestamp not null default now(),
    updated_at  timestamp not null default now(),
    deleted_at  timestamp not null default '0001-01-01T00:00:00Z'
);

CREATE UNIQUE INDEX uidx_endpoint_name ON endpoint(name);
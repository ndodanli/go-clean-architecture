drop table if exists role;
create table if not exists role
(
    id          bigserial primary key,
    name        varchar(255) not null,
    description varchar(255),
    endpoint_ids bigint[]    not null default '{}',

    created_at  timestamp    not null default now(),
    updated_at  timestamp    not null default now(),
    deleted_at  timestamp    not null default '0001-01-01T00:00:00Z'

);

CREATE UNIQUE INDEX role_name_uidx ON role(name);
CREATE INDEX role_endpoint_ids_idx ON role USING GIN(endpoint_ids);
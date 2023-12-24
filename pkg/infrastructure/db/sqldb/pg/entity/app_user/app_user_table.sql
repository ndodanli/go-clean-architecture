drop table if exists app_user cascade;
create table if not exists app_user
(
    id                    bigserial primary key,
    username              varchar(50)  not null,
    password              varchar(255) not null,
    first_name            varchar(50)  not null,
    last_name             varchar(50)  not null,
    phone_number          varchar(50)  not null,
    email                 varchar(50)  not null,
    email_confirmed       bool         not null default false,
    email_confirmation    jsonb        not null default '{
      "code": "",
      "expires_at": null
    }',
    fp_email_confirmation jsonb        not null default '{
      "code": "",
      "expires_at": null
    }',
    role_ids              int[]        not null default '{}',
    block_details         jsonb[]      not null default '{}',

    created_at            timestamp    not null default now(),
    updated_at            timestamp    not null default now(),
    deleted_at            timestamp    not null default '0001-01-01T00:00:00Z'
);

create unique index if not exists app_user_email_uidx
    on app_user (email);
create unique index if not exists app_user_username_uidx
    on app_user (username);
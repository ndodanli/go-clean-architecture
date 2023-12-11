drop table if exists app_user cascade;
create table if not exists app_user
(
    id         bigserial primary key,
    username   varchar(50)  not null,
    password   varchar(255) not null,
    email      varchar(50)  not null,
    email_confirmed bool not null default false,
    fp_email_confirmation_details jsonb not null default '{"code": "", "expires_at": null}',

    created_at timestamp    not null default now(),
    updated_at timestamp    not null default now(),
    deleted_at timestamp    not null default '0001-01-01T00:00:00Z'
);

--     test_string_not_null varchar(50) not null,
--     test_string_nullable varchar(50),
--     test_int_not_null int not null,
--     test_int_nullable int,
--     test_bool_not_null bool not null,
--     test_bool_nullable bool,
--     test_date_not_null date not null,
--     test_date_nullable date,
--     test_timestamp_not_null timestamp not null,
--     test_timestamp_nullable timestamp,
--     test_float_not_null float not null,
--     test_float_nullable float,
--     test_double_not_null double precision not null,
--     test_double_nullable double precision,
--     test_decimal_not_null decimal not null,
--     test_decimal_nullable decimal,
--     test_enum_not_null test_enum_type not null,
--     test_enum_nullable test_enum_type,
--     test_json_not_null json not null,
--     test_json_nullable json,
--     test_jsonb_not_null jsonb not null,
--     test_jsonb_nullable jsonb,
--     test_bytea_not_null bytea not null,
--     test_bytea_nullable bytea,
--     test_array_not_null int[] not null,
--     test_array_nullable int[],

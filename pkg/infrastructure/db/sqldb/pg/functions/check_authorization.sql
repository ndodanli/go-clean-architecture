CREATE OR REPLACE FUNCTION check_authorization(app_user_id_arg INT, endpoint_name_arg TEXT, endpoint_method_arg TEXT)
    RETURNS TABLE
            (
                authorized  BOOLEAN,
                "roleIds"   BIGINT[],
                "isBlocked" BOOLEAN
            )
AS
$$
DECLARE
    authorized_var  BOOLEAN;
    role_ids_var    BIGINT[];
    is_blocked_var  BOOLEAN;
    endpoint_id_var BIGINT;
BEGIN
    SELECT id
    INTO endpoint_id_var
    FROM endpoint
    WHERE name = endpoint_name_arg
      AND method = endpoint_method_arg
      AND deleted_at = '0001-01-01T00:00:00Z'
    LIMIT 1;

    -- Pre-expand user roles
    WITH app_user_var AS (SELECT role_ids,
                                 block_details
                          FROM app_user
                          WHERE id = app_user_id_arg
                            AND deleted_at = '0001-01-01T00:00:00Z'),
         expanded_roles AS (SELECT UNNEST(role_ids) AS role_id
                            FROM app_user_var),
         expanded_blocks AS (SELECT elem ->> 'untilAt' AS until_at
                             FROM app_user_var
                                      CROSS JOIN LATERAL unnest(app_user_var.block_details) AS elem)
    -- Check authorization and block status
    SELECT INTO authorized_var, role_ids_var, is_blocked_var EXISTS (SELECT 1
                                                                     FROM expanded_roles er
                                                                              JOIN role r ON r.id = er.role_id
                                                                     WHERE r.deleted_at = '0001-01-01T00:00:00Z'
                                                                       AND endpoint_id_var = ANY (r.endpoint_ids)),
                                                             ARRAY(SELECT role_id
                                                                   FROM expanded_roles),
                                                             EXISTS (SELECT 1
                                                                     FROM expanded_blocks eb
                                                                     WHERE eb.until_at::timestamptz > NOW());


    RETURN QUERY SELECT authorized_var, role_ids_var, is_blocked_var;
END;
$$ LANGUAGE plpgsql;

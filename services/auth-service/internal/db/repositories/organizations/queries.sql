-- name: CreateOrganization :exec
INSERT INTO organizations (
    id,
    name,
    slug,
    created_at,
    created_by

) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: CreateOrganizationMembership :exec
INSERT INTO users_organizations (
    organization_id,
    user_id,
    role_id,
    created_at
) VALUES (
    $1,
    $2,
    $3,
    $4
);

-- name: ListOrganizationMembershipsByUserId :many
SELECT
    uo.organization_id AS organization_id,
    o.slug AS organization_slug,
    uo.role_id AS role_id,
    uo.created_at AS created_at
FROM
    users_organizations uo
JOIN
    organizations o ON o.id = uo.organization_id
WHERE
    uo.user_id = $1
ORDER BY
    uo.created_at DESC;
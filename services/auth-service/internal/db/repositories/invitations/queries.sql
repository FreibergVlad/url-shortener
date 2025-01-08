-- name: CreateInvitation :exec
INSERT INTO invitations (
    id,
    status,
    organization_id,
    email,
    role_id,
    created_at,
    expires_at,
    created_by

) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
);

-- name: GetInvitationById :one
SELECT * FROM invitations WHERE id = $1;

-- name: UpdateInvitationStatusById :exec
UPDATE invitations SET status = $1 WHERE id = $2;
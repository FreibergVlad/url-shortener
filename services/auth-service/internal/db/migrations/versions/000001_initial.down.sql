BEGIN;

DROP INDEX IF EXISTS organizations_slug_uidx;
DROP INDEX IF EXISTS invitations_org_id_email_uidx;
DROP INDEX IF EXISTS users_email_uidx;

DROP TABLE IF EXISTS users_organizations;
DROP TABLE IF EXISTS invitations;
DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS users;

COMMIT;
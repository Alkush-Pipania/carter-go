-- Drop tables in reverse order of creation (to handle foreign key constraints)
DROP TABLE IF EXISTS fallback_images;
DROP TABLE IF EXISTS forgot_passwords;
DROP TABLE IF EXISTS verifications;
DROP TABLE IF EXISTS linkforms;
DROP TABLE IF EXISTS folders;
DROP TABLE IF EXISTS users;

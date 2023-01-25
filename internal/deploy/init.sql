BEGIN;

CREATE SCHEMA ozon;

CREATE TABLE ozon.links(
    short_link VARCHAR(10) PRIMARY KEY,
    original_link VARCHAR(2048) UNIQUE
);

COMMIT;
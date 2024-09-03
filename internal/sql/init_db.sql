-- Create a new UTF-8 `snippetbox` database.
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Switch to using the `snippetbox` database.
USE snippetbox;

-- Create a `snippet` table.
CREATE TABLE snippet (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

-- Add an index on the created column.
CREATE INDEX idx_snippet_created ON snippet(created);

INSERT INTO snippet (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
INSERT INTO snippet (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
INSERT INTO snippet (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);
COMMIT;


-- This user gets 'Access denied' error when connect to database with 
-- DSN "web:webpwd@tcp(localhost:3306)/snippetbox?parseTime=true" .
------ CREATE USER 'web'@'localhost';
------ GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
------ ALTER USER 'web'@'localhost' IDENTIFIED BY 'webpwd';


-- This user can successfully connect to database with 
-- DSN "zeb:zebpwd@tcp(localhost:3306)/snippetbox?parseTime=true" .
CREATE USER zeb;
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO zeb;
ALTER USER zeb IDENTIFIED BY 'zebpwd';
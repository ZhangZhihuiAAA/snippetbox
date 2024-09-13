CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;


USE snippetbox;


CREATE TABLE snippet (
    id      INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title   VARCHAR(100) NOT NULL,
    content TEXT         NOT NULL,
    created DATETIME     NOT NULL,
    expires DATETIME     NOT NULL
);

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
-- DSN "zzh:zzhpwd@tcp(localhost:3306)/snippetbox?parseTime=true" .
CREATE USER zzh;
GRANT SELECT, INSERT, UPDATE, DELETE ON zsnippetbox.* TO zzh;
ALTER USER zzh IDENTIFIED BY 'zzhpwd';


CREATE TABLE sessions (
    token  CHAR(43)     PRIMARY KEY,
    data   BLOB         NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX idx_sessions_expiry ON sessions (expiry);


CREATE TABLE user (
    id              INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    hashed_password CHAR(60)     NOT NULL,
    created         DATETIME     NOT NULL
);

ALTER TABLE user ADD CONSTRAINT uc_user_email UNIQUE (email);



CREATE DATABASE test_snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER test_web;
GRANT CREATE, DROP, ALTER, INDEX, SELECT, INSERT, UPDATE, DELETE ON test_snippetbox.* TO test_web;
ALTER USER test_web IDENTIFIED BY 'test';
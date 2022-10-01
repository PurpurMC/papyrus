CREATE TABLE projects
(
    id         TEXT PRIMARY KEY NOT NULL,
    name       TEXT             NOT NULL UNIQUE,
    created_at DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX projects_name_idx ON projects (name);

CREATE TABLE versions
(
    id         TEXT PRIMARY KEY NOT NULL,
    name       TEXT             NOT NULL,
    project_id TEXT             NOT NULL,
    created_at DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (project_id, name),
    FOREIGN KEY (project_id) REFERENCES projects (id)
);

CREATE INDEX versions_name_idx ON versions (name);
CREATE INDEX versions_project_id_idx ON versions (project_id);

CREATE TABLE builds
(
    id         TEXT PRIMARY KEY NOT NULL,
    name       TEXT             NOT NULL,
    version_id TEXT             NOT NULL,
    result     TEXT             NOT NULL,
    duration   INTEGER          NOT NULL,
    timestamp  INTEGER          NOT NULL,
    created_at DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (version_id, name),
    FOREIGN KEY (version_id) REFERENCES versions (id)
);

CREATE INDEX builds_name_idx ON builds (name);
CREATE INDEX builds_version_id_idx ON builds (version_id);

CREATE TABLE commits
(
    id          TEXT PRIMARY KEY NOT NULL,
    build_id    TEXT             NOT NULL,
    author      TEXT             NOT NULL,
    email       TEXT             NOT NULL,
    description TEXT             NOT NULL,
    hash        TEXT             NOT NULL,
    timestamp   INTEGER          NOT NULL,
    created_at DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (build_id) REFERENCES builds (id)
);

CREATE INDEX commits_build_id_idx ON commits (build_id);

CREATE TABLE files
(
    id         TEXT PRIMARY KEY NOT NULL,
    build_id   TEXT             NOT NULL,
    hash       TEXT             NOT NULL,
    extension  TEXT             NOT NULL DEFAULT '',
    created_at DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (build_id) REFERENCES builds (id)
);

CREATE INDEX files_build_id_idx ON files (build_id);

CREATE TABLE temp_files
(
    id         TEXT PRIMARY KEY NOT NULL,
    extension  TEXT             NOT NULL,
    created_at DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX temp_files_created_at_idx ON temp_files (created_at);

CREATE TABLE IF NOT EXISTS sets (
    id BIGSERIAL PRIMARY KEY,
    set_uid bytea UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS sets_datas (
    id BIGSERIAL PRIMARY KEY,
    data bigint,
    set_uid bytea NOT NULL REFERENCES sets (set_uid) ON DELETE CASCADE
);
CREATE INDEX set_datas_uid_index ON sets_datas (set_uid);
CREATE TABLE IF NOT EXISTS intersecting_sets (
    id BIGSERIAL PRIMARY KEY,
    set_uid bytea NOT NULL REFERENCES sets (set_uid) ON DELETE CASCADE,
    intersectingset_uid bytea NOT NULL REFERENCES sets (set_uid) ON DELETE CASCADE
);

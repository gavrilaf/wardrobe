BEGIN;

CREATE TABLE IF NOT EXISTS info_objects
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(512) NOT NULL,
    author          VARCHAR(512) NOT NULL,
    source          VARCHAR(512) NOT NULL,
    published       TIMESTAMPTZ NOT NULL,
    created         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finalized       TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS info_objects_name_index ON info_objects (name);
CREATE INDEX IF NOT EXISTS info_objects_author_index ON info_objects (author);
CREATE INDEX IF NOT EXISTS info_objects_source_index ON info_objects (source);
CREATE INDEX IF NOT EXISTS info_objects_published_index ON info_objects (published);
CREATE INDEX IF NOT EXISTS info_objects_created_index ON info_objects (created);

CREATE TABLE IF NOT EXISTS files
(
    id              SERIAL PRIMARY KEY,
    info_object_id  INT NOT NULL,
    bucket          VARCHAR(256) NOT NULL,
    name            VARCHAR(512) NOT NULL,
    content_type    VARCHAR(20) NOT NULL,
    size            INT NOT NULL,
    uploaded        TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unq_file_path UNIQUE (bucket, name),
    CONSTRAINT fk_file_info_object
        FOREIGN KEY (info_object_id)
            REFERENCES info_objects (id)
);

CREATE TABLE IF NOT EXISTS tags
(
    id    SERIAL PRIMARY KEY,
    value VARCHAR(100) NOT NULL UNIQUE
);
CREATE UNIQUE INDEX IF NOT EXISTS unq_tags_value_index ON tags (value);

CREATE TABLE IF NOT EXISTS info_objects_tags
(
    id              SERIAL PRIMARY KEY,
    info_object_id  INTEGER NOT NULL,
    tag_id          INTEGER NOT NULL,
    CONSTRAINT unq_info_object_tag UNIQUE (info_object_id, tag_id),
    CONSTRAINT fk_info_object_tag_info_object_id
        FOREIGN KEY (info_object_id)
            REFERENCES info_objects (id),
    CONSTRAINT fk_info_object_tag_tag_id
        FOREIGN KEY (tag_id)
            REFERENCES tags (id)
);

END;
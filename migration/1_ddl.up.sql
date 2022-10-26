BEGIN;

CREATE TABLE IF NOT EXISTS file_objects
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(512) NOT NULL ,
    content_type    VARCHAR(20) NOT NULL ,
    size            INT NOT NULL DEFAULT 0,
    created         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    uploaded        TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS unq_file_objects_name_index ON file_objects (name);

CREATE TABLE IF NOT EXISTS tags
(
    id    SERIAL PRIMARY KEY,
    value VARCHAR(150) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS file_objects_tags
(
    id              SERIAL PRIMARY KEY,
    file_object_id  INTEGER NOT NULL,
    tag_id          INTEGER NOT NULL,
    CONSTRAINT unq_file_object_tag UNIQUE (file_object_id, tag_id),
    CONSTRAINT fk_file_object_tag_file_object_id
        FOREIGN KEY (file_object_id)
            REFERENCES file_objects (id),
    CONSTRAINT fk_file_object_tag_tag_id
        FOREIGN KEY (tag_id)
            REFERENCES tags (id)
);

END;
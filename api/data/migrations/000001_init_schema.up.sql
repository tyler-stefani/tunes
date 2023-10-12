CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE,
    email VARCHAR UNIQUE,
    password VARCHAR
);
CREATE TABLE artist (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);
CREATE TABLE project (
    id BIGINT PRIMARY KEY,
    title VARCHAR NOT NULL,
    form VARCHAR NOT NULL,
    release DATE NOT NULL
);
CREATE TABLE track (
    id BIGINT PRIMARY KEY,
    title VARCHAR NOT NULL,
    primary_project_id BIGINT
);
CREATE TABLE spin (
    id SERIAL PRIMARY KEY,
    time TIMESTAMP NOT NULL,
    user_id BIGINT,
    track_id BIGINT
);
CREATE TABLE artist_project (artist_id BIGINT, project_id BIGINT);
CREATE TABLE artist_track (artist_id BIGINT, track_id BIGINT);
CREATE TABLE project_track (project_id BIGINT, track_id BIGINT);
ALTER TABLE track
ADD FOREIGN KEY (primary_project_id) REFERENCES project (id);
ALTER TABLE spin
ADD FOREIGN KEY (user_id) REFERENCES "user" (id);
ALTER TABLE spin
ADD FOREIGN KEY (track_id) REFERENCES track (id);
ALTER TABLE artist_project
ADD FOREIGN KEY (artist_id) REFERENCES artist (id);
ALTER TABLE artist_project
ADD FOREIGN KEY (project_id) REFERENCES project (id);
ALTER TABLE artist_track
ADD FOREIGN KEY (artist_id) REFERENCES artist (id);
ALTER TABLE artist_track
ADD FOREIGN KEY (track_id) REFERENCES track (id);
ALTER TABLE project_track
ADD FOREIGN KEY (project_id) REFERENCES project (id);
ALTER TABLE project_track
ADD FOREIGN KEY (track_id) REFERENCES track (id);
-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE faculties
(
    id     uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    number text NOT NULL,
    name   text NOT NULL
);

CREATE TABLE departments
(
    id        uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    number    text NOT NULL,
    name      text NOT NULL,
    facultyId uuid REFERENCES faculties (id)
);

CREATE TABLE groups
(
    id           uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    number       text NOT NULL,
    course       int  NOT NULL,
    isMagistracy bool                      DEFAULT FALSE,
    departmentId uuid REFERENCES departments (id)
);

CREATE TABLE students
(
    id      uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    groupId uuid NOT NULL REFERENCES groups (id),
    number  text NOT NULL
);

CREATE TABLE teachers
(
    id         uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    speciality text,
    name       text NOT NULL
);


CREATE TABLE users
(
    id       uuid    NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    type     integer NOT NULL,
    name     text    NOT NULL,
    login    text    NOT NULL,
    password text    NOT NULL,
    ownerId  uuid    NOT NULL,
    token    text    NOT NULL
);

CREATE TABLE subjects
(
    id        uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacherId uuid NOT NULL REFERENCES teachers (id),
    groupId   uuid NOT NULL REFERENCES groups (id),
    name      text,
    type      int                       DEFAULT 1
);

CREATE UNIQUE INDEX ON subjects (groupId, name);

CREATE TABLE lesson
(
    id            uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    type          integer,
    subjectId     uuid NOT NULL REFERENCES subjects (id),
    name          text,
    couple        int  NOT NULL,
    day           int  NOT NULL,
    groupId       uuid NOT NULL REFERENCES groups (id),
    teacherId     uuid REFERENCES teachers (id),
    cabinet       text,
    isDenominator bool                      DEFAULT FALSE,
    isNumerator   bool                      DEFAULT FALSE
);

CREATE UNIQUE INDEX ON lesson (teacherId, couple, day, isNumerator, isDenominator);
-- CREATE UNIQUE INDEX ON lesson (teacherId, couple, day, isDenominator);
CREATE UNIQUE INDEX ON lesson (groupid, couple, day, isNumerator);
CREATE UNIQUE INDEX ON lesson (groupid, couple, day, isDenominator);

CREATE TABLE lessons_progress
(
    id        uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    lessonId  uuid NOT NULL REFERENCES lesson (id),
    studentId uuid NOT NULL REFERENCES students (id),
    points    integer
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

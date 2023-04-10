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
    id         uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    number     text NOT NULL,
    name       text NOT NULL,
    faculty_id uuid REFERENCES faculties (id)
);

CREATE TABLE groups
(
    id            uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    number        text NOT NULL,
    course        int  NOT NULL,
    is_magistracy bool                      DEFAULT FALSE,
    department_id uuid REFERENCES departments (id)
);

CREATE TABLE students
(
    id       uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id uuid NOT NULL REFERENCES groups (id),
    number   text NOT NULL
);

CREATE TABLE teachers
(
    id         uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    speciality text NOT NULL
);


CREATE TABLE users
(
    id       uuid    NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    type     integer NOT NULL,
    name     text    NOT NULL,
    login    text    NOT NULL,
    password text    NOT NULL,
    owner_id uuid    NOT NULL
);

CREATE TABLE subjects
(
    id         uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacher_id uuid NOT NULL REFERENCES teachers (id),
    group_id   uuid NOT NULL REFERENCES groups (id)
);

CREATE TABLE lesson
(
    id         uuid    NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    type       integer NOT NULL,
    subject_id uuid    NOT NULL REFERENCES subjects (id),
    module     int,
    name       text
);

CREATE TABLE lessons_progress
(
    id         uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    lesson_id  uuid NOT NULL REFERENCES lesson (id),
    student_id uuid NOT NULL REFERENCES students (id),
    points     integer
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

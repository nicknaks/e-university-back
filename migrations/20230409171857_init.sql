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
    name    text NOT NULL
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
    login    text    NOT NULL,
    password text    NOT NULL,
    ownerId  uuid    NOT NULL,
    token    text    NOT NULL
);

CREATE TABLE subjects
(
    id           uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacherId    uuid NOT NULL REFERENCES teachers (id),
    groupId      uuid NOT NULL REFERENCES groups (id),
    name         text,
    type         int                       DEFAULT 1,
    addTeacherId uuid REFERENCES teachers (id)
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
    isNumerator   bool                      DEFAULT FALSE,
    addTeacherId  uuid REFERENCES teachers (id)
);

CREATE UNIQUE INDEX ON lesson (groupid, couple, day, isNumerator);
CREATE UNIQUE INDEX ON lesson (groupid, couple, day, isDenominator);

CREATE TABLE classes
(
    id        uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    day       date,
    type      int,
    comment   text,
    name      text,
    subType   int,
    module    int,
    subjectId uuid NOT NULL REFERENCES subjects (id),
    lessonId  uuid,
    groupId   uuid NOT NULL REFERENCES groups (id)
);

CREATE TABLE classes_progresses
(
    id        uuid  NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    classId   uuid  NOT NULL REFERENCES classes (id),
    studentId uuid  NOT NULL REFERENCES students (id),
    isAbsent  bool  NOT NULL             DEFAULT false,
    teacherId uuid REFERENCES teachers (id),
    mark      float NOT NULL             DEFAULT 0
);

CREATE TABLE subjects_results
(
    id                      uuid  NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    studentId               uuid  NOT NULL REFERENCES students (id),
    subjectId               uuid  NOT NULL REFERENCES subjects (id),
    firstModuleMark         float NOT NULL             DEFAULT 0,
    firstModuleMarkComment  text,
    secondModuleMark        float NOT NULL             DEFAULT 0,
    secondModuleMarkComment text,
    thirdModuleMark         float NOT NULL             DEFAULT 0,
    thirdModuleMarkComment  text,
    mark                    float NOT NULL             DEFAULT 0,
    examResult              int   NOT NULL             DEFAULT 0,
    examResultComment       text
);

CREATE OR REPLACE FUNCTION update_results_proc()
    RETURNS trigger AS
$$
DECLARE
    currentSubjectID uuid;
    currentModule    int;
BEGIN
    SELECT classes.subjectId, classes.module
    into currentSubjectID, currentModule
    FROM classes
    WHERE classes.id = new.classId;

    CASE
        WHEN currentModule = 1 THEN UPDATE subjects_results
                                    SET firstModuleMark = firstModuleMark + new.mark - old.mark,
                                        mark            = mark + new.mark - old.mark
                                    WHERE subjects_results.subjectId = currentSubjectID
                                      AND subjects_results.studentId = new.studentId;
        WHEN currentModule = 2 THEN UPDATE subjects_results
                                    SET secondModuleMark = secondModuleMark + new.mark - old.mark,
                                        mark             = mark + new.mark - old.mark
                                    WHERE subjects_results.subjectId = currentSubjectID
                                      AND subjects_results.studentId = new.studentId;
        WHEN currentModule = 3 THEN UPDATE subjects_results
                                    SET thirdModuleMark = thirdModuleMark + new.mark - old.mark,
                                        mark            = mark + new.mark - old.mark
                                    WHERE subjects_results.subjectId = currentSubjectID
                                      AND subjects_results.studentId = new.studentId;
        END CASE;

    RETURN new;
END;
$$
    LANGUAGE 'plpgsql';


CREATE OR REPLACE FUNCTION update_subject_result()
    RETURNS trigger AS
$$
BEGIN
    CASE
        WHEN old.examResult != new.examResult THEN new.mark = new.mark + new.examResult - old.examResult;
        WHEN old.examResult = new.examResult THEN RETURN new;
        END CASE;

    RETURN new;
END;
$$
    LANGUAGE 'plpgsql';


CREATE TRIGGER update_exam
    BEFORE UPDATE
    ON subjects_results
    FOR EACH ROW
EXECUTE PROCEDURE update_subject_result();

CREATE TRIGGER update_results
    BEFORE UPDATE
    ON classes_progresses
    FOR EACH ROW
EXECUTE PROCEDURE update_results_proc();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

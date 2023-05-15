// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Class struct {
	ID              string           `json:"id"`
	Day             string           `json:"day"`
	Type            LessonType       `json:"type"`
	Comment         *string          `json:"comment"`
	Name            *string          `json:"name"`
	Module          int              `json:"module"`
	SubjectID       string           `json:"subjectID"`
	LessonID        string           `json:"lessonID"`
	GroupID         string           `json:"groupID"`
	StudentProgress []*ClassProgress `json:"studentProgress"`
}

type ClassProgress struct {
	ID        string  `json:"id"`
	ClassID   string  `json:"classID"`
	StudentID string  `json:"studentID"`
	IsAbsent  bool    `json:"isAbsent"`
	TeacherID *string `json:"teacherID"`
	Mark      int     `json:"mark"`
}

type Department struct {
	ID        string   `json:"id"`
	Number    string   `json:"number"`
	Name      string   `json:"name"`
	FacultyID string   `json:"facultyID"`
	Groups    []*Group `json:"groups"`
}

type Faculty struct {
	ID          string        `json:"id"`
	Number      string        `json:"number"`
	Name        string        `json:"name"`
	Departments []*Department `json:"departments"`
}

type Group struct {
	ID       string     `json:"id"`
	Number   string     `json:"number"`
	Course   int        `json:"course"`
	Students []*Student `json:"students"`
}

type Lesson struct {
	ID            string     `json:"id"`
	Type          LessonType `json:"type"`
	SubjectID     string     `json:"subjectID"`
	Name          *string    `json:"name"`
	Couple        int        `json:"couple"`
	Day           int        `json:"day"`
	GroupID       string     `json:"groupID"`
	TeacherID     *string    `json:"teacherID"`
	Cabinet       *string    `json:"cabinet"`
	IsDenominator bool       `json:"isDenominator"`
	IsNumerator   bool       `json:"isNumerator"`
	Teacher       *Teacher   `json:"teacher"`
	Group         *Group     `json:"group"`
}

type Student struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	GroupID string `json:"groupId"`
}

type Subject struct {
	ID        string      `json:"id"`
	TeacherID *string     `json:"teacherID"`
	GroupID   string      `json:"groupID"`
	Name      *string     `json:"name"`
	Group     *Group      `json:"group"`
	Teacher   *Teacher    `json:"teacher"`
	Type      SubjectType `json:"type"`
}

type SubjectResult struct {
	ID               string     `json:"id"`
	StudentID        string     `json:"studentID"`
	SubjectID        string     `json:"subjectID"`
	Subject          []*Subject `json:"subject"`
	FirstModuleMark  int        `json:"firstModuleMark"`
	SecondModuleMark int        `json:"secondModuleMark"`
	ThirdModuleMark  int        `json:"thirdModuleMark"`
	Mark             int        `json:"mark"`
	//  оценка за предмет
	Total int `json:"total"`
}

type Teacher struct {
	ID   string  `json:"id"`
	Name *string `json:"name"`
}

type User struct {
	ID      string   `json:"id"`
	OwnerID *string  `json:"owner_id"`
	Type    UserType `json:"type"`
}

type AbsentSetInput struct {
	ClassProgressID []string `json:"classProgressID"`
}

type ClassesFilter struct {
	Ids       []string `json:"ids"`
	SubjectID *string  `json:"subjectID"`
	GroupID   *string  `json:"groupID"`
}

type ClassesProgressFilter struct {
	ClassID *string `json:"classID"`
}

type GroupsFilter struct {
	IDIn         []string `json:"idIn"`
	DepartmentID *string  `json:"departmentID"`
	Course       *int     `json:"course"`
	IsMagistracy *bool    `json:"isMagistracy"`
}

type LessonCreateInput struct {
	SubjectID     string     `json:"subjectID"`
	Type          LessonType `json:"type"`
	Couple        int        `json:"couple"`
	Day           int        `json:"day"`
	Cabinet       *string    `json:"cabinet"`
	IsDenominator bool       `json:"isDenominator"`
	IsNumerator   bool       `json:"isNumerator"`
}

type MarkCreateInput struct {
	ClassProgressID string `json:"classProgressID"`
	Mark            int    `json:"mark"`
}

type ScheduleFilter struct {
	GroupID   *string `json:"groupID"`
	TeacherID *string `json:"teacherID"`
}

type StudentCreateInput struct {
	Name    string `json:"name"`
	GroupID string `json:"groupID"`
}

type StudentsFilter struct {
	GroupID   *string  `json:"groupID"`
	SubjectID *string  `json:"subjectID"`
	IDIn      []string `json:"idIn"`
}

type SubjectCreateInput struct {
	Name      string      `json:"name"`
	Type      SubjectType `json:"type"`
	TeacherID string      `json:"teacherID"`
	GroupID   string      `json:"groupID"`
}

type SubjectResultsFilter struct {
	SubjectID *string `json:"subjectID"`
	StudentID *string `json:"studentID"`
}

type SubjectTypeChangeInput struct {
	ID   string      `json:"id"`
	Type SubjectType `json:"type"`
}

type SubjectsFilter struct {
	ID        []string `json:"ID"`
	GroupID   *string  `json:"groupID"`
	TeacherID *string  `json:"teacherID"`
}

type TeachersFilter struct {
	IDIn []string `json:"idIn"`
}

type LessonType string

const (
	LessonTypeDefault LessonType = "DEFAULT"
	LessonTypeLab     LessonType = "LAB"
	LessonTypeLec     LessonType = "LEC"
	LessonTypeSem     LessonType = "SEM"
)

var AllLessonType = []LessonType{
	LessonTypeDefault,
	LessonTypeLab,
	LessonTypeLec,
	LessonTypeSem,
}

func (e LessonType) IsValid() bool {
	switch e {
	case LessonTypeDefault, LessonTypeLab, LessonTypeLec, LessonTypeSem:
		return true
	}
	return false
}

func (e LessonType) String() string {
	return string(e)
}

func (e *LessonType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = LessonType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid LessonType", str)
	}
	return nil
}

func (e LessonType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SubjectType string

const (
	SubjectTypeUnknown    SubjectType = "UNKNOWN"
	SubjectTypeCredit     SubjectType = "CREDIT"
	SubjectTypeExam       SubjectType = "EXAM"
	SubjectTypeCourseWork SubjectType = "COURSE_WORK"
	SubjectTypePractical  SubjectType = "PRACTICAL"
)

var AllSubjectType = []SubjectType{
	SubjectTypeUnknown,
	SubjectTypeCredit,
	SubjectTypeExam,
	SubjectTypeCourseWork,
	SubjectTypePractical,
}

func (e SubjectType) IsValid() bool {
	switch e {
	case SubjectTypeUnknown, SubjectTypeCredit, SubjectTypeExam, SubjectTypeCourseWork, SubjectTypePractical:
		return true
	}
	return false
}

func (e SubjectType) String() string {
	return string(e)
}

func (e *SubjectType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SubjectType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SubjectType", str)
	}
	return nil
}

func (e SubjectType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UserType string

const (
	UserTypeUnknown UserType = "UNKNOWN"
	UserTypeTeacher UserType = "TEACHER"
	UserTypeStudent UserType = "STUDENT"
	UserTypeAdmin   UserType = "ADMIN"
)

var AllUserType = []UserType{
	UserTypeUnknown,
	UserTypeTeacher,
	UserTypeStudent,
	UserTypeAdmin,
}

func (e UserType) IsValid() bool {
	switch e {
	case UserTypeUnknown, UserTypeTeacher, UserTypeStudent, UserTypeAdmin:
		return true
	}
	return false
}

func (e UserType) String() string {
	return string(e)
}

func (e *UserType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserType", str)
	}
	return nil
}

func (e UserType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

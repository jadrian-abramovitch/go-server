package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type DbService struct {
	db *sql.DB
}

type DbCourse struct {
	CourseID      int
	Instructor    string
	Name          string
	Prerequisites string
}

func (course Course) toDb() DbCourse {
	var prereqs strings.Builder
	for i, val := range course.Prerequisites {
		if i != 0 {
			prereqs.WriteString(",")
		}
		prereqs.WriteString(strconv.Itoa(val))
	}
	return DbCourse{CourseID: course.CourseID, Instructor: course.Instructor, Name: course.Name, Prerequisites: prereqs.String()}

}

func (dbCourse DbCourse) toApi() (Course, error) {
	fmt.Print("test to see logs")
	fmt.Printf("prereqs: %s", dbCourse.Prerequisites)
	splicedString := strings.Split(dbCourse.Prerequisites, ",")
	var prereqArray []int
	for _, val := range splicedString {
		courseId, err := strconv.Atoi(val)
		if err != nil {
			return Course{}, err
		}
		prereqArray = append(prereqArray, courseId)
	}
	return Course{CourseID: dbCourse.CourseID, Instructor: dbCourse.Instructor, Name: dbCourse.Name, Prerequisites: prereqArray}, nil
}

func (dbService DbService) getAllCourses() ([]Course, error) {
	rows, err := dbService.db.Query("SELECT * FROM courses")
	if err != nil {
		return []Course{}, err
	}
	var courses []Course
	for rows.Next() {
		var result DbCourse
		if err := rows.Scan(&result.CourseID, &result.Instructor, &result.Name, &result.Prerequisites); err != nil {
			return []Course{}, err
		}
		translatedCourse, err := result.toApi()
		if err != nil {
			return []Course{}, err
		}
		courses = append(courses, translatedCourse)
	}
	return courses, nil
}

func (dbService DbService) createCourse(newCourse Course) error {
	course := newCourse.toDb()
	if _, err := dbService.db.Exec("INSERT INTO courses (CourseId, Instructor, Name, Prerequisites) VALUES (?,?,?,?)", course.CourseID, course.Instructor, course.Name, course.Prerequisites); err != nil {
		return err
	}
	return nil
}

func (dbService DbService) deleteCourse(courseId int, db *sql.DB) error {
	if _, err := db.Exec("DELETE FROM courses WHERE CourseID=?", courseId); err != nil {
		return err
	}
	return nil
}

func (dbService DbService) updateCourse(courseId int, updatedCourse Course) error {
	course := updatedCourse.toDb()
	if _, err := dbService.db.Exec("UPDATE courses SET Instructor = ?, Name = ?, Prerequisites = ? WHERE CourseID = ?", course.Instructor, course.Name, course.Prerequisites, courseId); err != nil {
		return err
	}

	return nil

}

func (dbService DbService) getCourse(courseId int) (Course, error) {
	row := dbService.db.QueryRow("SELECT * FROM courses WHERE CourseID = ?", courseId)
	var dbCourse DbCourse
	if err := row.Scan(&dbCourse.CourseID, &dbCourse.Instructor, &dbCourse.Name, &dbCourse.Prerequisites); err != nil {
		return Course{}, err
	}
	course, err := dbCourse.toApi()
	if err != nil {
		return Course{}, err
	}
	return course, nil
}

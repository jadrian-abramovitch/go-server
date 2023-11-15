package main

import "database/sql"

type DbService struct {
	db *sql.DB
}

func (dbService DbService) getAllCourses() ([]Course, error) {
	rows, err := dbService.db.Query("SELECT * FROM courses")
	if err != nil {
		return []Course{}, err
	}
	var courses []Course
	for rows.Next() {
		var result Course
		if err := rows.Scan(&result.CourseID, &result.Instructor, &result.Name); err != nil {
			return []Course{}, err
		}
		courses = append(courses, result)
	}
	return courses, nil
}

func (dbService DbService) createCourse(newCourse Course) error {
	if _, err := dbService.db.Exec("INSERT INTO courses (CourseId, Instructor, Name) VALUES (?,?,?)", newCourse.CourseID, newCourse.Instructor, newCourse.Name); err != nil {
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
	if _, err := dbService.db.Exec("UPDATE courses SET Instructor = ?, Name = ? WHERE CourseID = ?", updatedCourse.Instructor, updatedCourse.Name, courseId); err != nil {
		return err
	}

	return nil

}

func (dbService DbService) getCourse(courseId int) (Course, error) {
	row := dbService.db.QueryRow("SELECT * FROM courses WHERE CourseID = ?", courseId)
	var course Course
	if err := row.Scan(&course.CourseID, &course.Instructor, &course.Name); err != nil {
		return Course{}, err
	}
	return course, nil
}

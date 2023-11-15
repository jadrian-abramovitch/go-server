package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Course struct {
	CourseID   int    `json:"courseId" db:"courseId"`
	Instructor string `json:"instructor" db:"instructor"`
	Name       string `json:"name" db:"name"`
	// PreRequisites []int  `json:"preRequisites"`
}

var initialCourses = []Course{
	{1, "Smith", "Calculus"},
	{2, "Chen", "Philosophy"},
	{3, "Anderson", "Calculus 2"},
}

func setupTable(db *sql.DB) error {
	if _, err := db.Exec("DROP TABLE IF EXISTS courses"); err != nil {
		return err
	}
	if _, err := db.Exec("CREATE TABLE COURSES (CourseID INT PRIMARY KEY NOT NULL, Instructor CHAR(50), Name CHAR(50))"); err != nil {
		return err
	}
	for _, course := range initialCourses {
		if _, err := db.Exec("INSERT INTO courses (CourseId, Instructor, Name) VALUES (?,?,?)", course.CourseID, course.Instructor, course.Name); err != nil {
			return err
		}
	}

	_, err := db.Query("SELECT * FROM courses")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "file:test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := setupTable(db); err != nil {
		fmt.Println("Couldn't setup db", err)
		return
	}

	dbService := DbService{db}
	handlers := HandlerService{dbService}

	r := gin.Default()
	r.GET("/course/:id", handlers.getCourseHandler)
	r.GET("/courses", handlers.getAllCoursesHandler)
	r.POST("/newCourse", handlers.createCourseHandler)
	r.DELETE("/course/:id", handlers.deleteCourseHandler)
	r.PATCH("/course/:id", handlers.updateCourseHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

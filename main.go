package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Course struct {
	CourseID      int    `json:"courseId" db:"courseId" form:"courseId"`
	Instructor    string `json:"instructor" db:"instructor" form:"instructor"`
	Name          string `json:"name" db:"name" form:"name"`
	Prerequisites []int  `json:"preRequisites" form:"preRequisites"`
}

var initialCourses = []Course{
	{1, "Smith", "Calculus", []int{2, 3}},
	{2, "Chen", "Philosophy", []int{3, 1}},
	{3, "Anderson", "Calculus 2", []int{1}},
}

func setupTable(db *sql.DB, dbService DbService) error {
	if _, err := db.Exec("DROP TABLE IF EXISTS courses"); err != nil {
		return err
	}
	if _, err := db.Exec("CREATE TABLE COURSES (CourseID INT PRIMARY KEY NOT NULL, Instructor CHAR(50), Name CHAR(50), Prerequisites CHAR(50))"); err != nil {
		return err
	}
	for _, course := range initialCourses {
		if err := dbService.createCourse(course); err != nil {
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
	dbService := DbService{db}
	if err := setupTable(db, dbService); err != nil {
		fmt.Println("Couldn't setup db", err)
		return
	}

	handlers := HandlerService{dbService}

	r := gin.Default()
	r.GET("/course/:id", handlers.getCourseHandler)
	r.GET("/courses", handlers.getAllCoursesHandler)
	r.POST("/newCourse", handlers.createCourseHandler)
	r.DELETE("/course/:id", handlers.deleteCourseHandler)
	r.PATCH("/course/:id", handlers.updateCourseHandler)

	r.LoadHTMLFiles("index.html", "edit.html")

	r.GET("/", func(c *gin.Context) {
		allCourses, err := dbService.getAllCourses()
		if err != nil {
			fmt.Println("no courses")
		}
		c.HTML(http.StatusOK, "index.html", allCourses)
	})
	r.GET("/edit/:id", func(c *gin.Context) {
		courseId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Println("could not parse id")
		}
		course, err := dbService.getCourse(courseId)
		if err != nil {
			fmt.Println("no courses")
		}
		c.HTML(http.StatusOK, "edit.html", course)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

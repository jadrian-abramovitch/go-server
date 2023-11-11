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
	CourseID   int    `json:"courseId" db:"courseId"`
	Instructor string `json:"instructor" db:"instructor"`
	Name       string `json:"name" db:"name"`
	// PreRequisites []int  `json:"preRequisites"`
}

var courses = []Course{
	{1, "Smith", "Calculus"},
	{2, "Chen", "Philosophy"},
	{3, "Anderson", "Calculus 2"},
}

// var courses = []Course{
// 	{1, "Smith", "Calculus", []int{}},
// 	{2, "Chen", "Philosophy", []int{}},
// 	{3, "Anderson", "Calculus 2", []int{1, 2}},
// }

func getAllCourses(c *gin.Context) {
	// c.IndentedJSON(http.StatusOK, courses)
}

func createCourse(c *gin.Context) {
	var newCourse Course
	if err := c.BindJSON(&newCourse); err != nil {
		return
	}

	courses = append(courses, newCourse)
	c.IndentedJSON(http.StatusCreated, newCourse)
}

func deleteCourse(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Could not parse ID")
		return
	}

	for i, course := range courses {
		if course.CourseID == courseId {
			courses = append(courses[:i], courses[i+1:]...)
			c.String(http.StatusOK, fmt.Sprintf("Course %d deleted", courseId))
			return
		}
	}
	c.String(http.StatusNotFound, "Could not find course with that ID")
}

func updateCourse(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Could not parse ID")
		return
	}

	var updatedCourse Course
	if err := c.BindJSON(&updatedCourse); err != nil {
		c.String(http.StatusNotFound, "Could not parse input")
		return
	}
	for i, course := range courses {
		if course.CourseID == courseId {
			courses = append(append(courses[:i], updatedCourse), courses[i+1:]...)
			c.String(http.StatusOK, fmt.Sprintf("Course %d updated", courseId))
			return
		}
	}
	c.String(http.StatusNotFound, "Could not find course with that ID")
}

func setupTable(db *sql.DB) error {
	fmt.Println("results for insert: ", courses[1])
	_, err := db.Exec("INSERT INTO courses (CourseId, Instructor, Name) VALUES (?,?,?)", courses[1].CourseID, courses[1].Instructor, courses[1].Name)
	if err != nil {
		return err
	}
	rows, err := db.Query("SELECT * FROM courses")
	if err != nil {
		return err
	}
	for rows.Next() {
		var result Course
		if err := rows.Scan(&result.CourseID, &result.Instructor, &result.Name); err != nil {
			return err
		}
		fmt.Println("results from select: ", result)
	}

	return nil
}

func main() {
	db, err := sql.Open("sqlite3",
		"file:test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := setupTable(db); err != nil {
		fmt.Println("Couldn't insert data", err)
		return
	}
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/courses", getAllCourses)
	r.POST("/newCourse", createCourse)
	r.DELETE("/course/:id", deleteCourse)
	r.PATCH("/course/:id", updateCourse)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

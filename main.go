package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Course struct {
	CourseID      int    `json:"courseId"`
	Instructor    string `json:"instructor"`
	Name          string `json:"name"`
	PreRequisites []int  `json:"preRequisites"`
}

var courses = []Course{
	{1, "Smith", "Calculus", []int{}},
	{2, "Chen", "Philosophy", []int{}},
	{3, "Anderson", "Calculus 2", []int{1, 2}},
}

func getAllCourses(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, courses)
}
func createCourse(c *gin.Context) {
	var newCourse Course

	if err := c.BindJSON(&newCourse); err != nil {
		return
	}

	courses = append(courses, newCourse)
	c.IndentedJSON(http.StatusCreated, newCourse)
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/courses", getAllCourses)
	r.POST("/newCourse", createCourse)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

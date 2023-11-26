package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerService struct {
	dbService DbService
}

func (handler HandlerService) getCourseHandler(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Could not parse ID")
		return
	}
	course, err := handler.dbService.getCourse(courseId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not find course")
		return
	}
	c.IndentedJSON(http.StatusOK, course)
}

func (handler HandlerService) getAllCoursesHandler(c *gin.Context) {
	result, err := handler.dbService.getAllCourses()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error()) // dev mode only
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}

func (handler HandlerService) createCourseHandler(c *gin.Context) {
	var newCourse Course
	if err := c.Bind(&newCourse); err != nil {
		fmt.Print(err)
		c.String(http.StatusInternalServerError, err.Error()) // dev mode only
		return
	}

	err := handler.dbService.createCourse(newCourse)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error()) // dev mode only
		return
	}
	c.String(http.StatusOK, "Course created succesfully")
}

func (handler HandlerService) deleteCourseHandler(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Could not parse ID")
		return
	}
	if err := handler.dbService.deleteCourse(courseId, handler.dbService.db); err != nil {
		c.String(http.StatusInternalServerError, "error deleting course")
	}
	c.String(http.StatusOK, fmt.Sprintf("Course with id %d deleted successfully", courseId))

}

func (handler HandlerService) updateCourseHandler(c *gin.Context) {
	courseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Could not parse ID")
		return
	}

	var updatedCourse Course
	if err := c.BindJSON(&updatedCourse); err != nil {
		c.String(http.StatusNotFound, err.Error())
		fmt.Printf("could not parse!!!, err: %s", err.Error())
		return
	}
	if err := handler.dbService.updateCourse(courseId, updatedCourse); err != nil {
		fmt.Printf(err.Error())
		c.String(http.StatusInternalServerError, "Could not update course")
		return
	}
	c.String(http.StatusOK, "Updated course succesfully")
}

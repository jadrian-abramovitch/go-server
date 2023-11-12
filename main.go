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

var initialCourses = []Course{
	{1, "Smith", "Calculus"},
	{2, "Chen", "Philosophy"},
	{3, "Anderson", "Calculus 2"},
}

func getAllCourses(db *sql.DB) ([]Course, error) {
	rows, err := db.Query("SELECT * FROM courses")
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

func createCourse(c *gin.Context, db *sql.DB) error {
	var newCourse Course
	if err := c.BindJSON(&newCourse); err != nil {
		return err
	}

	if _, err := db.Exec("INSERT INTO courses (CourseId, Instructor, Name) VALUES (?,?,?)", newCourse.CourseID, newCourse.Instructor, newCourse.Name); err != nil {
		return err
	}
	return nil
}

func deleteCourse(c *gin.Context, db *sql.DB) {
	courseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Could not parse ID")
		return
	}

	if _, err := db.Exec("DELETE FROM courses WHERE CourseID=?", courseId); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Course with id %d deleted successfully", courseId))
}

func updateCourse(c *gin.Context, db *sql.DB) {
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

	if _, err := db.Exec("UPDATE courses SET Instructor = ?, Name = ? WHERE CourseID = ?", updatedCourse.Instructor, updatedCourse.Name, courseId); err != nil {
		c.String(http.StatusInternalServerError, "Could not update value")
		return
	}

	c.String(http.StatusOK, "Course patched succesfully")
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

	rows, err := db.Query("SELECT * FROM courses")
	if err != nil {
		return err
	}

	//delete this later
	for rows.Next() {
		var result Course
		if err := rows.Scan(&result.CourseID, &result.Instructor, &result.Name); err != nil {
			return err
		}
		fmt.Println("results from select: ", result)
	}

	return nil
}

func getCourse(courseId int, db *sql.DB) (Course, error) {
	row := db.QueryRow("SELECT * FROM courses WHERE CourseID = ?", courseId)
	var course Course
	if err := row.Scan(&course.CourseID, &course.Instructor, &course.Name); err != nil {
		return Course{}, err
	}
	return course, nil
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

	r := gin.Default()
	r.GET("/course/:id", func(c *gin.Context) {
		courseId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(http.StatusBadRequest, "Could not parse ID")
			return
		}
		course, err := getCourse(courseId, db)
		if err != nil {
			c.String(http.StatusInternalServerError, "Could not find course")
			return
		}
		c.IndentedJSON(http.StatusOK, course)
	})
	r.GET("/courses", func(c *gin.Context) {
		result, err := getAllCourses(db)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error()) // dev mode only
			return
		}
		c.IndentedJSON(http.StatusOK, result)
	})
	r.POST("/newCourse", func(c *gin.Context) {
		err := createCourse(c, db)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error()) // dev mode only
			return
		}
		c.String(http.StatusOK, "Course created succesfully")
	})
	r.DELETE("/course/:id", func(c *gin.Context) {
		deleteCourse(c, db)
	})
	r.PATCH("/course/:id", func(c *gin.Context) {
		updateCourse(c, db)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

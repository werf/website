package controllers

import (
	"040_db/internal/services"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func RememberController(c *gin.Context) {
	dbType, dbPath := services.GetDBCredentials()

	db, err := sql.Open(dbType, dbPath)
	if err != nil {
		panic(err)
	}

	answer := c.Query("answer")
	name := c.Query("name")
	_, err = db.Exec("INSERT INTO talkers (answer, name) VALUES (?, ?)",
		answer, name)
	if err != nil {
		panic(err)
	}

	c.String(http.StatusOK, "Got it.\n")

	defer db.Close()
}

func SayController(c *gin.Context) {
	dbType, dbPath := services.GetDBCredentials()

	db, err := sql.Open(dbType, dbPath)
	if err != nil {
		panic(err)
	}

	result, err := db.Query("SELECT * FROM talkers")
	if err != nil {
		panic(err)
	}

	for result.Next() {
		var id int
		var answer string
		var name string

		err = result.Scan(&id, &answer, &name)
		if err != nil {
			panic(err)
		}

		c.String(http.StatusOK, answer+", "+name+"!\n")
		break
	}
}

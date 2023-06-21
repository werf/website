package controllers

import (
	"040_db/internal/services"
	"database/sql"
	"github.com/gin-gonic/gin"

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
	_, err = db.Exec("insert into talkers (answer, name) values (?, ?)",
		answer, name)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func SayController(*gin.Context) {

}

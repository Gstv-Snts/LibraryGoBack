package main

import (
	"library/controllers"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()

	//login
	router.GET("/user", controllers.UserHandler)
	//post new book
	router.POST("/book", controllers.PostBookHandler)
	//post new rent
	router.POST("/rent/:bookId", controllers.PostRentHandler)
	//get all books
	router.GET("/book", controllers.GetBooksHandler)
	//update book
	router.PUT("/book/:bookId", controllers.PutBookHandler)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal("Error on running router: ", err)
	}
}

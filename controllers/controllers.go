package controllers

import (
	"database/sql"
	"fmt"
	"library/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
)

func UserHandler(c *gin.Context) {
	//ARRANGE
	envMap, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error on reading .env file: ", err)
	}
	if c.Request.Header.Get("API_KEY") != envMap["API_KEY"] {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid API key"})
	} else {
		sqlDB, err := sql.Open("mysql", envMap["DSN"])
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "unable to connect to the database"})
		} else {
			defer sqlDB.Close()
			body := new(models.User)
			err = c.ShouldBindBodyWith(&body, binding.JSON)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid body"})
			} else {
				rows, err := sqlDB.Query(fmt.Sprintf("SELECT * FROM users WHERE email='%v';", body.Email))
				if err != nil {
					fmt.Println(err.Error())
					c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "email invalid"})
				} else {
					queryUser := new(models.User)
					rows.Next()
					err := rows.Scan(&queryUser.Id, &queryUser.Email, &queryUser.Password)
					if err != nil {
						log.Fatal("Error on scanning user from query: ", err)
					} else {
						if queryUser.Password == body.Password {
							c.IndentedJSON(http.StatusOK, gin.H{"auth": true})
						} else {
							c.IndentedJSON(http.StatusOK, gin.H{"auth": false})
						}
					}
				}
			}
		}
	}
}

func PostBookHandler(c *gin.Context) {
	envMap, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error on reading .env: ", err)
	}
	if c.Request.Header.Get("API_KEY") != envMap["API_KEY"] {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid API key"})
	} else {
		bodyBook := new(models.Book)

		err := c.ShouldBindBodyWith(&bodyBook, binding.JSON)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error on binding body"})
			log.Fatal(err.Error())
		} else {
			if bodyBook.Image == "" ||
				bodyBook.Tittle == "" ||
				bodyBook.Synopsis == "" ||
				bodyBook.Author == "" ||
				bodyBook.Genre == "" ||
				bodyBook.SystemEntryDate == "" ||
				bodyBook.Status.Description == "" {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid book"})
			} else {
				fmt.Println(bodyBook.SystemEntryDate)
				sqlDB, err := sql.Open("mysql", envMap["DSN"])
				if err != nil {
					log.Fatal("Error connecting to Planet Scale", err)
				}
				_, err = sqlDB.Exec(fmt.Sprintf("INSERT INTO books(tittle,author,genre,isActive,description,image,systemEntryDate,synopsis) VALUES('%v','%v','%v',%v,'%v','%v','%v','%v');",
					bodyBook.Tittle,
					bodyBook.Author,
					bodyBook.Genre,
					0,
					"Active",
					bodyBook.Image,
					bodyBook.SystemEntryDate,
					bodyBook.Synopsis))
				if err != nil {
					log.Fatal("Erro on the query: ", err.Error())
				} else {
					if err != nil {
						log.Fatal("Error on getting last inserted id: ", err)
					}
					c.IndentedJSON(http.StatusOK, gin.H{"message": "ok"})
				}
				err = sqlDB.Close()
				if err != nil {
					log.Fatal("Error closing Planet Scale db: ", err)
				}
			}
		}
	}
}

func PostRentHandler(c *gin.Context) {
	envMap, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error on reading .env: ", err)
	}
	if c.Request.Header.Get("API_KEY") != envMap["API_KEY"] {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid API key"})
	} else {
		newRent := new(models.Rent)
		err := c.ShouldBind(&newRent)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid rent"})
		}
		id := c.Param("bookId")
		if err != nil {
			log.Fatal("Error on getting param: ", err)
		}
		sqlDB, err := sql.Open("mysql", envMap["DSN"])
		if err != nil {
			log.Fatal("Error opening Plante Scale: ", err)
		}
		fmt.Println(newRent.WithdrawlDate)
		fmt.Println(newRent.DeliveryDate)
		_, err = sqlDB.Exec(fmt.Sprintf("INSERT INTO rents(bookId,studentName,class,withdrawlDate,deliveryDate) VALUES(%v,'%v','%v','%v','%v');",
			id, newRent.StudentName, newRent.Class, newRent.WithdrawlDate, newRent.DeliveryDate))
		if err != nil {
			log.Fatal("Error inserting new rent: ", err)
		}
		c.IndentedJSON(http.StatusAccepted, gin.H{"message": "inserted"})
	}
}

func GetBooksHandler(c *gin.Context) {
	envMap, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error reading .env: ", err)
	}
	if c.Request.Header.Get("API_KEY") != envMap["API_KEY"] {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid API key"})
	} else {
		sqlDB, err := sql.Open("mysql", envMap["DSN"])
		if err != nil {
			log.Fatal("Error opening Plante Scale: ", err)
		}
		rows, err := sqlDB.Query("SELECT * FROM books;")
		if err != nil {
			log.Fatal("Error querying books: ", err)
		}
		books := []models.Book{}
		for i := 0; rows.Next(); i++ {
			currentBook := new(models.Book)
			err := rows.Scan(
				&currentBook.Id,
				&currentBook.Tittle,
				&currentBook.Author,
				&currentBook.Genre,
				&currentBook.Status.IsActive,
				&currentBook.Status.Description,
				&currentBook.Image,
				&currentBook.SystemEntryDate,
				&currentBook.Synopsis)
			if err != nil {
				log.Fatal("Error scanning book: ", err)
			}
			rows, err := sqlDB.Query(fmt.Sprintf("SELECT studentName,class,withdrawlDate,deliveryDate FROM rents WHERE bookId=%v", currentBook.Id))
			if err != nil {
				log.Fatal("Error querying rents: ", err)
			}
			currentBookRents := []models.Rent{}
			for i := 0; rows.Next(); i++ {
				currentRent := new(models.Rent)
				err := rows.Scan(&currentRent.Class, &currentRent.DeliveryDate, &currentRent.StudentName, &currentRent.WithdrawlDate)
				if err != nil {
					log.Fatal("Erro scanning current rent: ", err)
				}
				currentBookRents = append(currentBookRents, *currentRent)
			}
			currentBook.RentHistory = currentBookRents
			books = append(books, *currentBook)
		}
		c.IndentedJSON(http.StatusOK, books)
	}
}

func PutBookHandler(c *gin.Context) {
	envMap, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error reading .env: ", err)
	}
	if c.Request.Header.Get("API_KEY") != envMap["API_KEY"] {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid API key"})
	} else {
		bookdId := c.Param("bookId")
		updatedBook := new(models.BookUpdate)
		err := c.ShouldBindWith(updatedBook, binding.JSON)
		if err != nil {
			log.Fatal("Error binding body: ", err)
		}

		sqlDB, err := sql.Open("mysql", envMap["DSN"])
		if err != nil {
			log.Fatal("Error connecting to Planet Scale: ", err)
		}
		_, err = sqlDB.Exec(fmt.Sprintf("UPDATE books SET tittle='%v',author='%v',genre='%v',image='%v',synopsis='%v',systemEntryDate='%v',description='%v',isActive=%v WHERE id=%v;",
			updatedBook.Tittle,
			updatedBook.Author,
			updatedBook.Genre,
			updatedBook.Image,
			updatedBook.Synopsis,
			updatedBook.SystemEntryDate,
			updatedBook.Status.Description,
			updatedBook.Status.IsActive,
			bookdId))
		if err != nil {
			fmt.Println("Error executing update: " + err.Error())
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid book"})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "book updated"})
		}
	}
}

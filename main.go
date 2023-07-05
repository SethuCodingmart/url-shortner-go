package main

import (
	"fmt"
	"log"
	database "urlShortner/database"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Printf("Hello Go URL Shortener !🚀")

	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	} else {
		fmt.Print("DATABASE CONNECTED!!")
	}

	r := gin.Default()

	// HEALTH CHECK
	r.GET("/health-check", func(res *gin.Context) {
		res.JSON(200, gin.H{
			"message": "URL SHORTNER UP AND RUNNING!!",
		})
	})

	// INIT PORT FOR SEVER
	errOnInit := r.Run(":3001")
	if errOnInit != nil {
		panic(fmt.Sprintf(`Server Not Started!! Reason:- %v`, err))
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	database "urlShortner/database"
	interfaceGo "urlShortner/interface"
	response "urlShortner/response"
	store "urlShortner/store"
	utils "urlShortner/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func staticCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set cache-related headers
		c.Header("Cache-Control", "public, max-age=3600")
		c.Header("Expires", time.Now().AddDate(1, 0, 0).Format(time.RFC1123))

		c.Next()
	}
}

func main() {
	fmt.Printf("Hello Go URL Shortener !🚀")
	// loadEnv()
	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	} else {
		fmt.Print("DATABASE CONNECTED!!")
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Use(staticCache())

	r.Static("/page", "./templates")
	// r.LoadHTMLGlob("templates/*.html")

	// HEALTH CHECK
	r.GET("/health-check", func(c *gin.Context) {
		response.Response(c, 200, "URL SHORTNER UP AND RUNNING!!", true, nil)
	})

	// ROUTES
	r.POST("/save-url", func(c *gin.Context) {
		var urlParameters interfaceGo.URLParameters
		if err := c.BindJSON(&urlParameters); err != nil {
			response.Response(c, 400, "SOMETHING WENT WRONG!!", false, err.Error())
			return
		}
		if urlParameters.Location == "" {
			response.Response(c, 400, "SOME FIELDS ARE INVALID!!", false, nil)
			return
		}
		if urlParameters.Alias == "" {
			result, err := utils.GenerateRandomString(5)
			if err != nil {
				response.Response(c, 400, "SOMETHING WENT WRONG!!", false, err.Error())
			}
			urlParameters.Alias = result
		}
		result, err := store.SaveURL(urlParameters)
		if err != nil {
			response.Response(c, 400, "SOMETHING WENT WRONG!!", false, err.Error())
			return
		}
		var finalResult interfaceGo.URLResultSuccess
		finalResult.Alias = urlParameters.Alias
		finalResult.Success = result
		response.Response(c, 200, "URL SHORTEN SUCCESS", true, finalResult)
	})

	r.GET("/:alias", func(c *gin.Context) {
		path := c.Param("alias")
		result := store.RedirectURL(path)
		c.Redirect(http.StatusTemporaryRedirect, result)
	})

	// r.GET("/page/404", func(c *gin.Context) {
	// 	c.Header("Cache-Control", "max-age=3600")
	// 	c.HTML(http.StatusOK, "404.html", gin.H{})
	// })

	// INIT PORT FOR SEVER
	errOnInit := r.Run(":3001")
	if errOnInit != nil {
		panic(fmt.Sprintf(`Server Not Started!! Reason:- %v`, err))
	}
}

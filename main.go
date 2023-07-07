package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	database "urlShortner/database"
	interfaceGo "urlShortner/interface"
	response "urlShortner/response"
	store "urlShortner/store"
	utils "urlShortner/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

type CustomError struct {
	message string
}

func (e CustomError) Error() string {
	return e.message
}

func GenerateToken(userId int, gmail string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": strconv.Itoa(userId),
		"gmail":  gmail,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func VerifyAndClaimToken(_tokenString string) (*string, error) {
	tokenString := strings.Split(_tokenString, " ")

	if tokenString[0] != "Bearer" {
		return nil, CustomError{"Invalid token"}
	}

	token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return nil, CustomError{"Invalid token" + err.Error()}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, CustomError{"Invalid token claims"}
	}
	userId, ok := claims["userId"].(string)
	if !ok {
		return nil, CustomError{"Invalid user ID in token"}
	}

	return &userId, nil
}

func staticCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set cache-related headers
		c.Header("Cache-Control", "public, max-age=3600")
		c.Header("Expires", time.Now().AddDate(1, 0, 0).Format(time.RFC1123))

		c.Next()
	}
}

func RequestIdMiddleware(c *gin.Context) {
	c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
	c.Next()
}

type AuthArrayStruct struct {
	path   string
	access []string
}

var AuthArray = []AuthArrayStruct{
	{
		path:   "/save-url",
		access: []string{"API", "WEB"},
	},
	{
		path:   "/user",
		access: []string{"WEB"},
	},
	{
		path:   "/create-authkey",
		access: []string{"WEB"},
	},
	{
		path:   "/get-urls",
		access: []string{"WEB"},
	},
}

func checkAccess(access []string, grantAccess string) bool {
	for _, r := range access {
		if r == grantAccess {
			return true
		}
	}
	return false
}

func authMiddleware(c *gin.Context) {
	sdkAuthorizationHeader := c.GetHeader("authkey")
	webAuthorizationHeader := c.GetHeader("Authorization")
	path := c.FullPath()
	accessCheck := false
	// Route Access
	if webAuthorizationHeader != "" {
		for _, r := range AuthArray {
			if path == r.path {
				accessCheck = checkAccess(r.access, "WEB")
			}
		}
		if !accessCheck {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			return
		}

		tokenVerifyResult, errToken := VerifyAndClaimToken(webAuthorizationHeader)
		if errToken != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":      "Unauthorized",
				"errMessage": errToken.Error(),
			})
			return
		}
		userId := &tokenVerifyResult
		intId, err := strconv.Atoi(**userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":      "Unauthorized",
				"errMessage": err.Error(),
			})
			return
		}
		userResult, err := store.GetUserWithId(intId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":      "Unauthorized",
				"errMessage": errToken.Error(),
			})
			return
		}
		c.Set("user", userResult)
		c.Next()
	} else if sdkAuthorizationHeader != "" {
		for _, r := range AuthArray {
			if path == r.path {
				accessCheck = checkAccess(r.access, "API")
			}
		}
		if !accessCheck {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			return
		}

		userResult, err := store.GetUserWithAuthKey(sdkAuthorizationHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":      "Unauthorized",
				"errMessage": err.Error(),
			})
			return
		}
		c.Set("user", userResult)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized := TNF",
		})
		return
	}
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func main() {
	fmt.Printf("Hello Go URL Shortener !🚀")
	LoadEnv()
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
	r.Use(RequestIdMiddleware)
	r.Static("/page", "./templates")
	// r.LoadHTMLGlob("templates/*.html")

	// HEALTH CHECK
	r.GET("/health-check", func(c *gin.Context) {
		result, err := store.HealthCheck()
		var resResult interfaceGo.HealthCheck
		if err == nil && result {
			resResult.DB = `DB RUNNING CHECK :- OKAY!!`
			resResult.APP = `DB RUNNING CHECK :- OKAY!!`
		} else {
			resResult.DB = `DB RUNNING CHECK :- FAILED!! ` + err.Error()
			resResult.APP = `DB RUNNING CHECK :- OKAY!!`
		}
		response.Response(c, 200, "URL SHORTNER UP AND RUNNING!!", true, resResult)
	})

	// ROUTES
	r.POST("/save-url", authMiddleware, func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists || user == nil {
			response.Response(c, http.StatusForbidden, "User Not Found", false, "USER NOT FOUND!!")
			return
		}
		userId := user.(*database.Users).Id

		var urlParameters interfaceGo.URLParameters
		urlParameters.Id = userId
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

	r.GET("/get-urls", authMiddleware, func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists || user == nil {
			response.Response(c, http.StatusForbidden, "User Not Found", false, "USER NOT FOUND!!")
			return
		}
		userId := user.(*database.Users).Id
		result, err := store.GetUrlsWithUserId(userId)
		if err != nil {
			response.Response(c, 400, "SOMETHING WENT WRONG!!", false, err.Error())
			return
		}
		response.Response(c, 200, "URLS FETCHEED SUCCESS", true, result)
	})

	r.POST("/register", func(c *gin.Context) {
		var user database.Users
		if err := c.BindJSON(&user); err != nil {
			response.Response(c, 400, "SOMETHING WENT WRONG!!", false, err.Error())
			return
		}
		hashedPassword, errHash := utils.HashPassword(user.Password)
		if errHash != nil {
			response.Response(c, 400, "SOMETHING WENT WRONG!!", false, errHash.Error())
			return
		}
		result, err := store.CreateUser(user.Gmail, hashedPassword)
		if err != nil {
			response.Response(c, 400, "SOMETHING WENT WRONG!!", false, err.Error())
			return
		}
		response.Response(c, 200, "USER CREATED SUCCESS", true, result)
	})

	r.POST("/login", func(c *gin.Context) {
		var creds interfaceGo.Login
		if err := c.BindJSON(&creds); err != nil {
			response.Response(c, 400, "SOMETHING WENT WRONG!!", false, err.Error())
			return
		}
		result, err := store.GetUserWithGmailAndPassword(creds.Gmail, creds.Password)
		if err != nil {
			response.Response(c, 400, "SOMETHING WENT WRONG!!", false, err.Error())
			return
		}
		token, err := GenerateToken(result.Id, result.Gmail)
		if err != nil {
			response.Response(c, 400, "SOMETHING WENT WRONG!!", false, err.Error())
			return
		}
		resultRes := make(map[string]interface{})
		resultRes["token"] = token
		resultRes["gmail"] = result.Gmail
		response.Response(c, 200, "USER LOGIN SUCCESS!!", true, resultRes)
	})

	r.GET("/user", authMiddleware, func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists || user == nil {
			response.Response(c, http.StatusForbidden, "User Not Found", false, "USER NOT FOUND!!")
			return
		}
		// userId := user.(*database.Users).Id
		// result, err := store.GetUserWithId(userId)
		// if err != nil {
		// 	response.Response(c, 400, "SOMETHING WENT WRONG!!", false, err.Error())
		// 	return
		// }
		response.Response(c, 200, "USER DATA FETCHED SUCCESS", true, user)
	})

	r.GET("/create-authkey", authMiddleware, func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists || user == nil {
			response.Response(c, http.StatusForbidden, "User Not Found", false, "USER NOT FOUND!!")
			return
		}
		userId := user.(*database.Users).Id
		// randString := randstr.String(20)
		randString := uuid.NewV4().String()
		result, err := store.GenerateAuthKey(randString, userId)
		if err != nil {
			response.Response(c, 400, "SOMETHING WENT WRONG!!", false, err.Error())
			return
		}
		resultResp := make(map[string]interface{})
		resultResp["authkey"] = randString
		resultResp["status"] = result
		response.Response(c, 200, "AUTH KEY CREATED SUCCESS", true, resultResp)
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

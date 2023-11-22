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

	"net/smtp"

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
		path:   "/user",
		access: []string{"WEB"},
	},
	{
		path:   "/transcations",
		access: []string{"WEB"},
	},
	{
		path:   "/create-workspace",
		access: []string{"WEB"},
	},
	// {
	// 	path:   "/save-url",
	// 	access: []string{"API", "WEB"},
	// },
	// {
	// 	path:   "/create-authkey",
	// 	access: []string{"WEB"},
	// },
	// {
	// 	path:   "/get-urls",
	// 	access: []string{"WEB"},
	// },
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
				"error":      "Forbidden",
				"errMessage": "Check API Path.",
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
				"errMessage": err.Error(),
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
				"error":      "Forbidden",
				"errMessage": "Check API Path.",
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

func sendMail(body string, subject string) bool {
	from := "rpsethu1471@gmail.com"
	pass := "rxhchzvnsrmlxuyc"
	to := "sethu1471@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject:" + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return false
	}

	log.Print("sent mail")
	return true
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
		AllowOrigins:     []string{"http://localhost:3001", "http://localhost:5173", "http://192.168.29.173:5173"},
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

	r.POST("/sendforgetpasswordotp", func(ctx *gin.Context) {
		var email interfaceGo.SendSignUpOTP
		if err := ctx.BindJSON(&email); err != nil {
			response.Response(ctx, 400, err.Error(), false, nil)
			return
		}
		checkUser, err := store.CheckUserExsist(email.Gmail)
		if !checkUser {
			response.Response(ctx, 400, "User not found.", false, nil)
			return
		}
		if err != nil {
			response.Response(ctx, 400, "Something Went Wrong. Please Try Again", false, nil)
			return
		}
		randomOtp, err := utils.GenerateRandomAlphaNumericString(6)
		if err != nil {
			response.Response(ctx, 400, "OTP Generation Failed.", false, nil)
			return
		}
		saveOtpParams := interfaceGo.SaveOTP{Key: email.Gmail, Value: randomOtp, Type: interfaceGo.FORGOT_PASSWORD}
		saveOtp, err := store.SaveOTP(saveOtpParams)
		if err != nil {
			response.Response(ctx, 400, "OTP Save Failed.", false, nil)
			return
		}
		if saveOtp {
			message := `Your Forget Password OTP for PlayURL is ` + randomOtp + ".\n\nThank You."
			if result := sendMail(message, "PlayURL - Forget Password"); !result {
				response.Response(ctx, 400, "Mail not sent", false, nil)
				return
			}
			response.Response(ctx, 200, "Mail Sent", true, nil)
		} else {
			response.Response(ctx, 400, "Mail Sent Failed. Try Again.", false, nil)
		}
	})

	r.POST("/forgetpassword", func(c *gin.Context) {
		var data interfaceGo.ForgetPassword
		if err := c.BindJSON(&data); err != nil {
			response.Response(c, 400, err.Error(), false, nil)
			return
		}
		otpVerify, err := store.GetOTP(data.Gmail, "FORGOT_PASSWORD")
		if err != nil {
			response.Response(c, 400, err.Error(), false, nil)
			return
		}
		if otpVerify.Value != data.Otp {
			response.Response(c, 400, "OTP not matched,", false, nil)
			return
		}
		hashedPassword, errHash := utils.HashPassword(data.Password)
		if errHash != nil {
			response.Response(c, 400, errHash.Error(), false, nil)
			return
		}
		result, err := store.UpdateUser(data.Gmail, hashedPassword)
		if err != nil {
			response.Response(c, 400, err.Error(), false, nil)
			return
		}
		response.Response(c, 200, "Password Changed.", true, result)
	})

	r.POST("/sendsignupotp", func(ctx *gin.Context) {
		var data interfaceGo.SendSignUpOTP
		if err := ctx.BindJSON(&data); err != nil {
			response.Response(ctx, 400, err.Error(), false, nil)
			return
		}
		checkUser, err := store.CheckUserExsist(data.Gmail)
		if checkUser {
			response.Response(ctx, 400, "User already registered.", false, nil)
			return
		}
		if err != nil {
			response.Response(ctx, 400, err.Error(), false, nil)
			return
		}
		randomOtp, err := utils.GenerateRandomAlphaNumericString(6)
		if err != nil {
			response.Response(ctx, 400, "OTP Generation Failed.", false, nil)
			return
		}
		saveOtpParams := interfaceGo.SaveOTP{Key: data.Gmail, Value: randomOtp, Type: interfaceGo.SIGNUP}
		saveOtp, err := store.SaveOTP(saveOtpParams)
		if err != nil {
			response.Response(ctx, 400, "OTP Save Failed.", false, nil)
			return
		}
		if saveOtp {
			message := `Your Verification OTP for PlayURL Signup is ` + randomOtp + ".\n\nThank You."
			if result := sendMail(message, "PlayURL Signup - Verification OTP"); !result {
				response.Response(ctx, 400, "Mail not sent", false, nil)
				return
			}
			response.Response(ctx, 200, "Mail Sent", true, nil)
		} else {
			response.Response(ctx, 400, "Mail Sent Failed. Try Again.", false, nil)
		}
	})

	r.POST("/register", func(c *gin.Context) {
		var data interfaceGo.Register
		if err := c.BindJSON(&data); err != nil {
			response.Response(c, 400, err.Error(), false, nil)
			return
		}
		otpVerify, err := store.GetOTP(data.Gmail, "SIGNUP")
		if err != nil {
			response.Response(c, 400, err.Error(), false, nil)
			return
		}
		if otpVerify.Value != data.Otp {
			response.Response(c, 400, "OTP not matched,", false, nil)
			return
		}
		hashedPassword, errHash := utils.HashPassword(data.Password)
		if errHash != nil {
			response.Response(c, 400, errHash.Error(), false, nil)
			return
		}
		result, err := store.CreateUser(data.Gmail, hashedPassword, data.Username, data.Name, data.Phone)
		if err != nil {
			response.Response(c, 400, err.Error(), false, nil)
			return
		}
		response.Response(c, 200, "User Created Success.", true, result)
	})

	r.POST("/login", func(c *gin.Context) {
		var creds interfaceGo.Login
		if err := c.BindJSON(&creds); err != nil {
			response.Response(c, 400, "Some fields are unfilled!!", false, err.Error())
			return
		}
		result, err := store.GetUserWithGmailAndPassword(creds.Gmail, creds.Password)
		if err != nil {
			response.Response(c, 400, err.Error(), false, nil)
			return
		}
		token, err := GenerateToken(result.Id, result.Gmail)
		if err != nil {
			response.Response(c, 400, err.Error(), false, nil)
			return
		}
		resultResponse := make(map[string]interface{})
		resultResponse["token"] = token
		resultResponse["gmail"] = result.Gmail
		resultResponse["username"] = result.Username
		resultResponse["id"] = result.Id
		response.Response(c, 200, "USER LOGIN SUCCESS!!", true, resultResponse)
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

	r.GET("/transcations", authMiddleware, func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists || user == nil {
			response.Response(c, http.StatusForbidden, "User Not Found", false, "USER NOT FOUND!!")
			return
		}
		userId := user.(*database.Users).Id
		result, err := store.GetTransactions(userId)
		if err != nil {
			response.Response(c, 400, err.Error(), false, nil)
			return
		}
		response.Response(c, 200, "Transcation Fetched Success.", true, result)
	})

	r.POST("/create-workspace", authMiddleware, func(ctx *gin.Context) {
		user, exists := ctx.Get("user")
		if !exists || user == nil {
			response.Response(ctx, http.StatusForbidden, "User Not Found", false, "USER NOT FOUND!!")
			return
		}
		userId := user.(*database.Users).Id
		var data interfaceGo.CreateWorkspace
		if err := ctx.BindJSON(&data); err != nil {
			response.Response(ctx, 400, "Some fields are unfilled!!", false, err.Error())
			return
		}
		if len(data.Shorthandname) < 4 {
			response.Response(ctx, 400, "ShortHandName should contain atleast 5 letters.", false, nil)
			return
		}
		result, err := store.SaveWorkspace(data.Name, data.Shorthandname, data.Description, userId)
		if err != nil {
			response.Response(ctx, 400, err.Error(), false, nil)
			return
		}
		returnData := make(map[string]interface{})
		returnData["workspaceId"] = result
		response.Response(ctx, 200, "Workspace Created Success.", true, returnData)
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
		resultData := make(map[string]interface{})
		resultData["authkey"] = randString
		resultData["status"] = result
		response.Response(c, 200, "AUTH KEY CREATED SUCCESS", true, resultData)
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

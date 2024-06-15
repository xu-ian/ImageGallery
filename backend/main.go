package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

var store cookie.Store = cookie.NewStore([]byte("secret"))

// SQL configurations
var cfg mysql.Config

// SQL database and errors
var db *sql.DB

// Image info structure for returned images
// Filename should be renamed Extension
type Image struct {
	Username string `json:"username"`
	Path     string `json:"filepath"`
	Filename string `json:"filename"`
	Id       string `json:"id"`
}

type Imglink struct {
	Link string
}

// Sign in info structure for login, signup, logout
type SigninInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Regex for detecting if a string is alphanumeric
func isAlphanumeric(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(str)
}

func isFileLink(str string) bool {
	return regexp.MustCompile(`^([a-zA-Z0-9_\-.:\/]+)$`).MatchString(str)
}

// Used to generate a random salt
func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func getFilenameExtension(filename string) string {
	splitname := strings.Split(filename, ".")
	return splitname[len(splitname)-1]
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Adds a new user to the database
func signup(c *gin.Context) {
	var signup SigninInfo
	var existingUsername int
	//Check if you can retrieve the username and password
	if err := c.BindJSON(&signup); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	//Check if username is alphanumeric
	if !isAlphanumeric(signup.Username) {
		c.JSON(http.StatusBadRequest, "Username can only contain alphanumeric characters")
		return
	}

	//Check whether the username already exists
	rows, err := db.Query("SELECT count(*) FROM users where username = ?", signup.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&existingUsername)
	if existingUsername == 1 {
		c.JSON(http.StatusUnprocessableEntity, "Username is taken")
		return
	}

	//Hash the password with a salt and insert the new user into the database
	salt := generateRandomString(32)

	hashedpass := sha256.Sum256([]byte(signup.Password + salt))
	stringhash := hex.EncodeToString(hashedpass[:]) + "." + salt
	log.Println(stringhash)
	_, err = db.Exec("INSERT INTO users(username, password) VALUES (?, ?)", signup.Username, stringhash)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "Cannot sign up, try again later")
		return
	}

	c.JSON(http.StatusOK, "SUCCESS")
}

// Logs in the user
func login(c *gin.Context) {
	var signin SigninInfo
	var dbinfo SigninInfo

	//Check if you can retrieve the username and password
	if err := c.BindJSON(&signin); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	//Check if username is alphanumeric
	if !isAlphanumeric(signin.Username) {
		c.JSON(http.StatusBadRequest, "Username can only contain alphanumeric characters")
		return
	}

	//Check whether the username exists
	rows, err := db.Query("SELECT * FROM users where username = ?", signin.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()
	if !rows.Next() {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}

	//Check whether the passwords match
	rows.Scan(&dbinfo.Username, &dbinfo.Password)
	dbpass := strings.Split(dbinfo.Password, ".")
	check := sha256.Sum256([]byte(signin.Password + dbpass[1]))
	if hex.EncodeToString(check[:]) != dbpass[0] {
		c.JSON(http.StatusUnauthorized, "Password is wrong")
		return
	}
	session := sessions.Default(c)
	session.Options(sessions.Options{HttpOnly: true, Secure: true})
	session.Set("username", signin.Username)
	session.Save()
	c.SetCookie("username", signin.Username, 3600, "/", "localhost", false, false)

	c.JSON(http.StatusOK, "success")
}

// Logs out the user
func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("username", "")
	session.Save()
	c.SetCookie("username", "", -1, "/", "localhost", false, false)
	c.JSON(http.StatusOK, "Logged out successfully")
}

// User posts an image to the database
func postImage(c *gin.Context) {
	username := c.Param("user")

	if !isAlphanumeric(username) {
		c.JSON(http.StatusBadRequest, "Invalid Username")
		return
	}

	if username != sessions.Default(c).Get("username") {
		c.JSON(http.StatusBadRequest, "Username does not match")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	fileExtension := getFilenameExtension(file.Filename)

	newName := generateRandomString(64)
	_, err = db.Exec("INSERT INTO images(id, username, imagepath, imagetype) VALUES (?, ?, ?, ?)",
		newName, username, "http://localhost:8000/images/"+newName+"."+fileExtension, fileExtension)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "Failed to upload image, try again")
		return
	}
	c.SaveUploadedFile(file, "./images/"+newName+"."+fileExtension)

	c.JSON(http.StatusOK, "Success")
}

// User adds a link to an image to the database
func postLink(c *gin.Context) {
	username := c.Param("user")

	if !isAlphanumeric(username) {
		c.JSON(http.StatusBadRequest, "Invalid Username")
		return
	}

	log.Println(sessions.Default(c).Get("username"))

	if username != sessions.Default(c).Get("username") {
		c.JSON(http.StatusBadRequest, "Username does not match")
		return
	}

	var link Imglink
	jsonData, _ := io.ReadAll(c.Request.Body)
	err := json.Unmarshal(jsonData, &link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if !isFileLink(link.Link) {
		c.JSON(http.StatusBadRequest, "Invalid link")
		return
	}

	_, err = db.Exec("INSERT INTO images (id, username, imagepath, imagetype) VALUES (?, ? , ? , ?)",
		generateRandomString(64), username, link.Link, "ext")

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "Unable to add image, try again")
		return
	}

	c.JSON(http.StatusOK, "Success")
}

// Gets all images posted by the user
func getImagesByUser(c *gin.Context) {
	username := c.Param("user")

	if !isAlphanumeric(username) {
		c.JSON(http.StatusBadRequest, "Invalid Username")
		return
	}

	//Get the rows where the image is equal to username
	rows, err := db.Query("SELECT * FROM images where username=?", username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var images []Image
	var imageInstance Image

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&imageInstance.Id, &imageInstance.Username, &imageInstance.Path, &imageInstance.Filename)
		images = append(images, imageInstance)
	}

	c.JSON(http.StatusOK, images)
}

// Gets image specified by user and image name
func getImagesByUserAndImageName(c *gin.Context) {
	log.Println("Hi")
	username := c.Param("User")
	imageName := c.Param("name")

	if !isAlphanumeric(username) || !isAlphanumeric(imageName) {
		c.JSON(http.StatusBadRequest, "Invalid request")
		return
	}

	//Find the element requested
	rows, err := db.Query("SELECT * FROM images where username=? AND filename=?", username, imageName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var imageInstance Image

	defer rows.Close()

	if rows.Next() {
		rows.Scan(&imageInstance.Id, &imageInstance.Username, &imageInstance.Path, &imageInstance.Filename)
	} else {
		c.JSON(http.StatusNotFound, "Image not found")
		return
	}

	c.JSON(http.StatusOK, imageInstance)
}

// Removes the image specified by user and image name
func removeImage(c *gin.Context) {
	username := c.Param("user")
	id := c.Param("name")

	if !isAlphanumeric(username) || !isAlphanumeric(id) {
		c.JSON(http.StatusBadRequest, "Invalid request")
		return
	}

	if username != sessions.Default(c).Get("username") {
		c.JSON(http.StatusBadRequest, "Username does not match")
		return
	}

	//Check to see if it exists
	rows, err := db.Query("SELECT * FROM images where username=? AND id=?", username, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var image Image

	defer rows.Close()
	if rows.Next() {
		rows.Scan(&image.Id, &image.Username, &image.Path, &image.Filename)
	} else {
		c.JSON(http.StatusNotFound, "Image was not found")
		return
	}

	_, err = db.Exec("DELETE FROM images where username=? AND id=?", username, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if image.Filename != "ext" {
		os.Remove("./Images/" + id + "." + image.Filename)
	}

	c.JSON(http.StatusOK, "Success")
}

// Returns the image file associated with the image name
func getImageFile(c *gin.Context) {
	filename := c.Param("name")
	c.File("./Images/" + filename)
	//c.Status(http.StatusOK)
	c.Set("content-type", "img/png")
}

func main() {

	cfg = mysql.Config{
		User:   os.Args[2],
		Passwd: os.Args[3],
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "imageshare",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		os.Exit(1)
	}
	fmt.Println("Connected to Database")

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.Use(sessions.Sessions("mysession", store))
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.POST("/logout", logout)
	router.POST("/:user/images", postImage)
	router.POST("/:user/link", postLink)
	router.GET("/:user/images", getImagesByUser)
	router.GET("/:user/images/:name", getImagesByUserAndImageName)
	router.GET("/images/:name", getImageFile)
	router.DELETE("/:user/images/:name", removeImage)
	router.Run("localhost:8000")

}

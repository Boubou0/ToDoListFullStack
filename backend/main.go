package main

import (
	"ToDoList/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"crypto/subtle"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/pbkdf2"
)

type Config struct {
	Port int `json:"port"`
}

func ChargerConfig(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	return corsConfig
}

func main() {
	config, err := ChargerConfig("config.json")
	if err != nil {
		log.Fatal("Erreur lors du chargement de la configuration:", err)
	}

	err = models.ConnectDatabase()
	checkErr(err)

	r := gin.Default()

	r.Use(cors.New(CORSConfig()))

	v1 := r.Group("/api/v1")
	{
		v1.GET("todolist", getTodoLists)
		v1.GET("todolist/:id", getTodoListById)
		v1.POST("todolist", addTodoList)
		v1.PUT("todolist/:id", updateTodoList)
		v1.DELETE("todolist/:id", deleteTodoList)

		v1.GET("tasks/:date", getTasksByDate)
		v1.GET("tasks/urgent", getAllUrgentTask)
		v1.POST("todolist/:id", addTask)
		v1.PUT("todolist/:id/task/:idTask", updateTask)
		v1.DELETE("todolist/:id/task/:idTask", deleteTask)

		v1.DELETE("task/:id", deleteTaskById)
		v1.PUT("task/:id", updateTaskById)

		v1.OPTIONS("todolist", options)
	}
	r.POST("/login", login)
	r.POST("/signup", insertUser)
	r.Use(verifyToken)

	r.GET("/info", getUserInfo)

	r.Run(":" + strconv.Itoa(config.Port))
}

func getUserInfo(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	c.JSON(http.StatusOK, user)
}

func verifyToken(c *gin.Context) {
	token, ok := getToken(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	id, username, err := validateToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	user, err := findUserByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}

	c.Set("id", id)
	c.Set("username", username)
	c.Set("user", user)

	c.Writer.Header().Set("Authorization", "Bearer "+token)
	c.Next()
}

func getToken(c *gin.Context) (string, bool) {
	authValue := c.GetHeader("Authorization")
	arr := strings.Split(authValue, " ")
	if len(arr) != 2 {
		return "", false
	}
	authType := strings.Trim(arr[0], "\n\r\t")
	if strings.ToLower(authType) != strings.ToLower("Bearer") {
		return "", false
	}
	return strings.Trim(arr[1], "\n\t\r"), true
}

func getSession(c *gin.Context) (uint, string, bool) {
	id, ok := c.Get("id")
	if !ok {
		return 0, "", false
	}
	username, ok := c.Get("username")
	if !ok {
		return 0, "", false
	}
	return id.(uint), username.(string), true
}

var jwtKey = []byte("FDr1VjVQiSiybYJrQZNt8Vfd7bFEsKP6vNX1brOSiWl0mAIVCxJiR4/T3zpAlBKc2/9Lw2ac4IwMElGZkssfj3dqwa7CQC7IIB+nVxiM1c9yfowAZw4WQJ86RCUTXaXvRX8JoNYlgXcRrK3BK0E/fKCOY1+izInW3abf0jEeN40HJLkXG6MZnYdhzLnPgLL/TnIFTTAbbItxqWBtkz6FkZTG+dkDSXN7xNUxlg==")

type authClaims struct {
	jwt.StandardClaims
	UserID uint `json:"userId"`
}

func generateToken(user User) (string, error) {
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, authClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: expiresAt,
		},
		UserID: user.ID,
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func validateToken(tokenString string) (uint, string, error) {
	var claims authClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return 0, "", err
	}
	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}
	id := claims.UserID
	username := claims.Subject
	return id, username, nil
}

type User struct {
	ID           uint
	Username     string
	PasswordSalt string
	PasswordHash string
}

const (
	saltSize  = 32
	iteration = 10000
	keyLen    = 64
)

func insertUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "incorrect parameters",
		})
		return
	}

	_, err := findUserByUsername(req.Username)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "username already taken",
		})
		return
	}

	salt := make([]byte, saltSize)
	_, err = rand.Read(salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error generating salt",
		})
		return
	}

	hashedPassword := pbkdf2.Key([]byte(req.Password), salt, iteration, keyLen, sha256.New)

	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashedPasswordBase64 := base64.StdEncoding.EncodeToString(hashedPassword)

	result, err := models.DB.Exec("INSERT INTO USER (username, password_salt, password_hash) VALUES (?, ?, ?)", req.Username, saltBase64, hashedPasswordBase64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error inserting user",
		})
		return
	}

	userID, _ := result.LastInsertId()
	user := &User{ID: uint(userID), Username: req.Username}

	c.JSON(http.StatusOK, gin.H{
		"message": "User inserted successfully",
		"user":    user,
	})

	return
}

func findUserByUsername(username string) (*User, error) {
	var user User
	err := models.DB.QueryRow("SELECT id, username, password_salt, password_hash FROM USER WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.PasswordSalt, &user.PasswordHash)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("user not found")
	case err != nil:
		return nil, err
	}

	return &user, nil
}

func findUserByID(id uint) (*User, error) {
	var user User
	err := models.DB.QueryRow("SELECT id, username, password_salt, password_hash FROM USER WHERE id = ?", id).
		Scan(&user.ID, &user.Username, &user.PasswordSalt, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "incorrect parameters",
		})
		return
	}

	user, err := findUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("user %s not found", req.Username),
		})
		return
	}

	if !verifyPassword(req.Password, user) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		return
	}

	token, err := generateToken(*user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"username": req.Username,
	})
}

func verifyPassword(inputPassword string, user *User) bool {
	salt, err := base64.StdEncoding.DecodeString(user.PasswordSalt)
	if err != nil {
		fmt.Println("Error decoding salt:", err)
		return false
	}

	hashedPasswordBytes, err := base64.StdEncoding.DecodeString(user.PasswordHash)
	if err != nil {
		fmt.Println("Error decoding hashed password:", err)
		return false
	}

	inputHashedPassword := pbkdf2.Key([]byte(inputPassword), salt, iteration, keyLen, sha256.New)
	return subtle.ConstantTimeCompare(hashedPasswordBytes, inputHashedPassword) == 1
}

func getTodoLists(c *gin.Context) {
	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	id, _, err := validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	user, err := findUserByID(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	todolists, err := models.GetTodoLists(user.ID)
	checkErr(err)

	if todolists == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": todolists})
	}
}

func getAllUrgentTask(c *gin.Context) {
	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	userID, _, err := validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	tasks, err := models.GetAllUrgentTasks(int(userID))
	checkErr(err)

	if tasks == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": tasks})
	}
}

func getTodoListById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	_, _, err = validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	taches, err := models.GetTodoListById(id)
	checkErr(err)

	if taches == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": taches})
	}
}

func getTasksByDate(c *gin.Context) {
	date := c.Param("date")
	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	id, _, err := validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	tasks, err := models.GetTasksByDate(int(id), date)
	checkErr(err)

	if tasks == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": tasks})
	}
}

func addTodoList(c *gin.Context) {
	var input struct {
		Nom string `json:"nom"`
	}
	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	id, _, err := validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	user, err := findUserByID(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	nouvelleListe := models.Todolist{
		Nom:    input.Nom,
		IdUser: int(user.ID),
	}

	err = models.AddTodoList(nouvelleListe, int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding todo list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo list added successfully"})
}

func addTask(c *gin.Context) {
	idTodolist, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todolist ID"})
		return
	}

	var input struct {
		Nom         string `json:"nom"`
		Description string `json:"description"`
		Etat        int    `json:"etat"`
		Priorite    int    `json:"priorite"`
		Deadline    string `json:"deadline"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	_, _, err = validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	nouvelleTache := models.Task{
		IdTodoList:  idTodolist,
		Nom:         input.Nom,
		Description: input.Description,
		Etat:        input.Etat,
		Priorite:    input.Priorite,
		Deadline:    input.Deadline,
	}

	err = models.AddTask(nouvelleTache)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task added successfully"})
}

func updateTask(c *gin.Context) {

	idTask, err := strconv.Atoi(c.Param("idTask"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}
	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	_, _, err = validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	err = models.UpdateTask(idTask, updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func updateTodoList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todolist ID"})
		return
	}

	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	_, _, err = validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var updatedTodoList models.Todolist
	if err := c.ShouldBindJSON(&updatedTodoList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	err = models.UpdateTodoList(id, updatedTodoList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating todo list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo list updated successfully"})
}

func deleteTodoList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todolist ID"})
		return
	}

	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	_, _, err = validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	err = models.DeleteTodoList(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting todo list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo list deleted successfully"})
}

func deleteTask(c *gin.Context) {
	idTodolist, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todolist ID"})
		return
	}

	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	_, _, err = validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	idTask, err := strconv.Atoi(c.Param("idTask"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	err = models.DeleteTask(idTodolist, idTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func deleteTaskById(c *gin.Context) {
	idTask, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	_, _, err = validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	err = models.DeleteTaskById(idTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
func updateTaskById(c *gin.Context) {
	idTask, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task ID"})
		return
	}

	token, ok := getToken(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	_, _, err = validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	err = models.UpdateTaskById(idTask, updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func options(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "options Called"})
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

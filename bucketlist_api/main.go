package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"fmt"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	// Postgres gorm driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DbConnect connects the app to the database
func dbConnect() (*gorm.DB, error) {
	godotenv.Load()
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_DATABASE")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	SSLMode := os.Getenv("SSL_MODE")
	dbType := os.Getenv("DB_TYPE")

	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s", dbHost, dbUser, dbName, SSLMode, dbPass)

	db, err := gorm.Open(dbType, connectionString)

	if err != nil {
		log.Fatalln("Error connecting to database", err)
		return nil, err
	}

	log.Println("DB Connection Successful")
	return db, nil
}

// Migrate function helps with the database migrations
func Migrate() {
	db, _ := dbConnect()
	defer db.Close()
	if err := db.AutoMigrate(&User{}).Error; err != nil {
		log.Fatalln("Error migrating the database ", err.Error())
	} else {
		log.Println("Migration successful...")
	}
}

// init is going to have the DB connections and any one-time tasks
func init() {
	Migrate()
}

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	// SetAllUserRoutes(router)
	router.HandleFunc("/signup", CreateEndPoint).Methods("POST")
	router.HandleFunc("/", Default).Methods("GET")
	router.HandleFunc("/del", DeleteEndPoint).Methods("POST")
	router.HandleFunc("/users", AllEndPoint).Methods("GET")
	router.HandleFunc("/search/{firstname}", SearchEndpoint).Methods("GET")
	router.HandleFunc("/users/{firstname}", EditEndPoint).Methods("PUT")
	router.HandleFunc("/users/{firstname}", DeleteEndPoint).Methods("DELETE")

	return router
}

type Message struct {
	Response     string
	StatusCode   uint
	ErrorMessage error
}

var message Message

// User struct is the user blueprint
type User struct {
	gorm.Model
	FirstName string `json:"first_name" gorm:"not null"`
	Surname   string `json:"surname" gorm:"not null"`
	UserEmail string `json:"email" gorm:"unique;not null"`
	Password  string `json:"password" gorm:"not null"`
}

var user User

// CreateEndPoint is a POST handler that posts a new user
var CreateEndPoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	db, err := dbConnect()
	if err != nil {
		panic("Connection to database failed")
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")

	password := []byte(r.FormValue("password"))
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	user = User{
		FirstName: r.FormValue("first_name"),
		Surname:   r.FormValue("surname"),
		UserEmail: r.FormValue("email"),
		Password:  string(hashedPassword),
	}

	feedback := db.Debug().Create(&user)
	if feedback.Error != nil {
		message.Response = "An error occured!"
		message.ErrorMessage = feedback.Error
		jsonmessage, _ := json.Marshal(message)
		w.Write([]byte(jsonmessage))
	} else {
		message.Response = "New user created successfully"
		message.StatusCode = 200
		message.ErrorMessage = nil
		jsonmessage, _ := json.Marshal(message)
		w.Write([]byte(jsonmessage))
	}

})

// EditEndPoint is a PUT handler that edits a database record
var EditEndPoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	db, err := dbConnect()
	if err != nil {
		panic("Connection to database failed")
	}
	defer db.Close()
	vars := mux.Vars(r)
	firstname := vars["firstname"]

	password := []byte(r.FormValue("password"))
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	feedback := db.Debug().Model(&user).Where("first_name = ?", firstname).Update(User{
		FirstName: r.FormValue("first_name"),
		Surname:   r.FormValue("surname"),
		UserEmail: r.FormValue("email"),
		Password:  string(hashedPassword),
	})
	if feedback.Error != nil {
		message.Response = "An error occured!"
		message.ErrorMessage = feedback.Error
		jsonmessage, _ := json.Marshal(message)
		w.Write([]byte(jsonmessage))
	} else {
		message.Response = "Database record Updated successfully"
		message.StatusCode = 200
		// message.ErrorMessage = nil
		jsonmessage, _ := json.Marshal(message)
		w.Write([]byte(jsonmessage))
	}

})

// AllEndPoint is a GET handler that fetches all users in the database
var AllEndPoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	db, err := dbConnect()
	if err != nil {
		panic("Connection to database failed")
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")

	var users []User
	feedback := db.Find(&users)
	if feedback.Error != nil {
		message.Response = "An error occured!"
		message.ErrorMessage = feedback.Error
		jsonmessage, _ := json.Marshal(message)
		w.Write([]byte(jsonmessage))
	} else {
		json, _ := json.Marshal(users)
		w.Write([]byte(json))
	}

})

// SearchEndpoint is a GET handler for searching for a specific user from the
// database using a first name as the unique parameter
var SearchEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	db, err := dbConnect()
	if err != nil {
		panic("Connection to database failed")
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	var fetchedUsers []User
	db.Where("first_name = ?", vars["firstname"]).First(&fetchedUsers)

	var message Message

	if len(fetchedUsers) == 0 {
		message.Response = fmt.Sprintf("the user with first name %s does not exist", vars["name"])
		message.StatusCode = 404
		// message.ErrorMessage = nil
		jsonmessage, _ := json.Marshal(message)
		w.Write([]byte(jsonmessage))
	} else {
		json, _ := json.Marshal(fetchedUsers)
		w.Write([]byte(json))
	}
})

// DeleteEndPoint handler deletes a user record using a given user name
var DeleteEndPoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	db, err := dbConnect()
	if err != nil {
		panic("Connection to database failed")
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	firstname := vars["firstname"]

	var user User
	newDB := db.Debug().Where("first_name = ?", firstname).Delete(&user)
	if newDB.Error != nil {
		panic(newDB.Error)
	}
	message.Response = "Database record deleted successfully"
	message.StatusCode = 200
	jsonmessage, _ := json.Marshal(message)
	w.Write([]byte(jsonmessage))

})

// Default endpoint
var Default = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	message.Response = "welcome to my GO App"
	message.StatusCode = 200
	message.ErrorMessage = nil
	jsonmessage, _ := json.Marshal(message)
	w.Write([]byte(jsonmessage))

})

// Define HTTP request routes
func main() {
	router := InitRoutes()
	log.Fatal(http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, router)))
}

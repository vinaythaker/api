// app.go

package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var config = "dev-config"

// App ...
type App struct {
	Router       *mux.Router
	db           *sql.DB
	api          *API
	server       *http.Server
	address      string
	readTimeOut  time.Duration
	writeTimeOut time.Duration
}

// Initialize ...
func (a *App) Initialize() {

	viper.SetConfigType("json")
	viper.SetConfigName(config)
	viper.AddConfigPath(".")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal(err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	a.address = viper.GetString("http-server.address")
	a.readTimeOut = viper.GetDuration("http-server.read-timeout") * time.Second
	a.writeTimeOut = viper.GetDuration("http-server.write-timeout") * time.Second

	a.Router, a.api = NewRouter()
	a.api.app = a
	a.api.petMap = make(map[int64]Pet)
	tag1 := Tag{1, "tag-dog-1"}
	tag2 := Tag{1, "tag-dog-2"}
	tag3 := Tag{1, "tag-cat-3"}
	tag4 := Tag{1, "tag-cat-4"}

	dogCat := Category{1, "dog"}
	catCat := Category{2, "cat"}
	dogTags := []Tag{tag1, tag2}
	dogUrls := []string{"dog-x", "dog-y", "dog-z"}

	a.api.petMap[1] = Pet{1, "woof", dogCat, "active", dogTags, dogUrls}

	catTags := []Tag{tag3, tag4}
	catUrls := []string{"cat-x", "cat-y", "cat-z"}

	a.api.petMap[2] = Pet{2, "meow", catCat, "active", catTags, catUrls}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s connect_timeout=%d",
		viper.GetString("db-server.host"),
		viper.GetString("db-server.port"),
		viper.GetString("db-server.user-id"),
		viper.GetString("db-server.secret"),
		viper.GetString("db-server.db-name"),
		viper.GetString("db-server.ssl-mode"),
		viper.GetInt("db-server.connect_timeout"),
	)

	log.Println("connecting to database")

	a.db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("creating db")

	err = a.createDB(a.db)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("creating http server")

	a.server = &http.Server{
		Handler:      a.Router,
		Addr:         a.address,
		ReadTimeout:  a.readTimeOut,
		WriteTimeout: a.writeTimeOut,
	}
}

// Run ...
func (a *App) Run() {
	// Start Server
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}

/*
func (a *App) Run() {
	// Start Server
	if err := a.server.ListenAndServe(); err != nil {
		log.Println(err)
	} else {
		log.Println(err)
	}
}
*/

func (a *App) createDB(db *sql.DB) error {

	const tableDropQuery = `DROP TABLE IF EXISTS pets`

	const tableCreationQuery = `
	CREATE TABLE IF NOT EXISTS pets
	(
	    id SERIAL PRIMARY KEY,
	    name VARCHAR(50) NOT NULL
	)`

	db.Exec(tableDropQuery)

	_, err := db.Exec(tableCreationQuery)
	return err
}

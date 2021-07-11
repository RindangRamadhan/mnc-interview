package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RindangRamadhan/mnc-interview/internal/db"
	"github.com/RindangRamadhan/mnc-interview/internal/router"
	"github.com/RindangRamadhan/mnc-interview/internal/tracing"
	"github.com/joho/godotenv"

	_ "github.com/RindangRamadhan/mnc-interview/internal/modules"
)

func Init() {
	db.Postgre.Init()
	db.Mongo.Init()
	tracing.Init()
}

func main() {
	var err error

	/**
	 * If want to change specific location file .env
	 * you can use this:
	 * `godotenv.Load(os.ExpandEnv("../PATH_TO_YOUR_FOLDER/.env"))`
	 */
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	/**
	 * Uncomment below if you want to connect to databse PostgreSQL
	 */
	// db.Conn, err = postgres.Get.Connect()
	// if err != nil {
	// 	log.Fatal("There was error when connecting to database.", err)
	// }

	fmt.Printf("Listing for requests at %s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	err = http.ListenAndServe(
		":"+os.Getenv("PORT"),
		router.Router,
	)

	if err != nil {
		log.Fatal("There was error when starting HTTP server", err)
	}
}

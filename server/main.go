package main

import (
	"clipboard-share-server/db"
	"clipboard-share-server/routes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	devFlag := flag.Bool("dev", false, "set dev mode")
	flag.Parse()

	if *devFlag {
		err := godotenv.Load()
		if err != nil {
			panic(fmt.Errorf("could not open .env in dev mode: %w", err))
		}
	}

	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		panic("ENV VAR DB_HOST NOT SET")
	}
	dbUser, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		panic("ENV VAR POSTGRES_USER NOT SET")
	}
	dbPass, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		panic("ENV VAR POSTGRES_PASSWORD NOT SET")
	}
	dbDB, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		panic("ENV VAR POSTGRES_DB NOT SET")
	}

	userDbLimitEnvVar, ok := os.LookupEnv("USER_DB_LIMIT")
	if !ok {
		db.PersonDbLimit = 1024 * 1024
	} else {
		lim, err := strconv.ParseInt(userDbLimitEnvVar, 10, 32)
		if err != nil {
			panic(fmt.Errorf("COULD NOT PARSE ENV VAR USER_DB_LIMIT: %w", err))
		}
		db.PersonDbLimit = int32(lim)
	}

	dbConn := fmt.Sprintf("postgres://%s:%s@%s/%s", dbUser, dbPass, dbHost, dbDB)

	err := db.Init(dbConn)
	if err != nil {
		panic(fmt.Errorf("error initializing database: %w", err))
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic("ENV VAR PORT NOT SET")
	}

	router := routes.NewRouter()

	server := http.Server{Addr: ":" + port, Handler: router}
	fmt.Printf("Serving on %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

package main

import (
	"os"

	"tunes-service/cache"
	"tunes-service/data"
	"tunes-service/server"
)

func main() {
	err := data.Migrate(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	db, _ := data.NewDB(os.Getenv("DATABASE_URL"))
	adb, _ := data.NewAuthDB(os.Getenv("AUTH_DATABASE_URL"))
	c := cache.NewCache()
	server.RunServer(db, adb, c)
}

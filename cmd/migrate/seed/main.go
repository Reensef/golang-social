package main

import (
	"log"

	"github.com/Reensef/golang-social/internal/db"
	"github.com/Reensef/golang-social/internal/env"
	"github.com/Reensef/golang-social/internal/seed"
	"github.com/Reensef/golang-social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/golang_social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStorage(conn)

	seed.GenerateFakeData(store)

	log.Default().Println("Seeding complete")
}

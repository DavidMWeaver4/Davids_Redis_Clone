package main

import (
	"log"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/server"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func main() {
	s := store.New()
	srv := server.New(":6379", s)

	log.Fatal(srv.ListenAndServe())
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/persistence"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/server"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func main() {
	s := store.New()
	defer s.Close()
	aof, err := persistence.NewAOF(persistence.AOFConfig{
		Path:        "appendonly.aof",
		FsyncPolicy: persistence.FsyncEverySec,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer aof.Close()

	srv := server.New(":6379", s, aof)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Printf("server error: %v", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Printf("server shutdown signal acknowledged")
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = srv.Shutdown(shutdown)
	if err != nil {
		log.Printf("server shutdown timed out or failed: %v", err)
		return
	}
	log.Println("server has been successfully shutdown")
}

package main

import (
	"context"
	"grace/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Connect to the database
	db := server.SetupDatabase("app")
	server.MigrateDatabase(db)

	httpServer := server.StartRESTServer(db)

	// Listen and serve in a goroutine, so main() does not block
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe() error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server exited properly")
}

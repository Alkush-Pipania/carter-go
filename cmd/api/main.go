package main

import (
	"context"
	"log"

	"github.com/Alkush-Pipania/carter-go/config"
	"github.com/Alkush-Pipania/carter-go/internal/app"
	"github.com/Alkush-Pipania/carter-go/internal/server"
	"github.com/Alkush-Pipania/carter-go/pkg/db"
)

func main() {

	cfg := config.LoadEnv()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initalize DB
	dbConn := db.Init(ctx, cfg.DbUrl)

	q := db.New(dbConn)

	container := app.NewContainer(ctx, q)

	// initialize the router
	router := app.NewRouter(container)

	// configure the server
	srv := server.New(router, cfg.Port)

	// start the server
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatalf("Server failed to start:", err)
	}
}

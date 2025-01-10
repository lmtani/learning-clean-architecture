package main

import (
	"database/sql"
	"fmt"

	"github.com/lmtani/learning-clean-architecture/configs"
	"github.com/lmtani/learning-clean-architecture/internal/infra/web/server"

	// postgres
	_ "github.com/lib/pq"
)

func main() {
	conf, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(conf.DBDriver, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	http := server.NewWebServer(conf.WebServerPort)
	httpOrderHandler := NewWebOrderHandler(db)
	http.AddHandler("/order", httpOrderHandler.Create)
	fmt.Println("Starting web server on port", conf.WebServerPort)
	http.Start()
}

package main

import (
	"database/sql"
	"fmt"

	"github.com/lmtani/learning-clean-architecture/configs"
	"github.com/lmtani/learning-clean-architecture/internal/infra/event/handler"
	"github.com/lmtani/learning-clean-architecture/internal/infra/web/server"
	"github.com/lmtani/learning-clean-architecture/pkg/events"
	"github.com/streadway/amqp"

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

	rabbitMQChannel := getRabbitMQChannel()
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	fmt.Println("Starting web server on port", conf.WebServerPort)
	http := server.NewWebServer(conf.WebServerPort)
	httpOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	http.AddHandler("/order", httpOrderHandler.Create)
	http.Start()

}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

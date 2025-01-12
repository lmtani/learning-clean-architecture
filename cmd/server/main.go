package main

import (
	"database/sql"
	"fmt"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/lmtani/learning-clean-architecture/configs"
	"github.com/lmtani/learning-clean-architecture/internal/infra/event/handler"
	"github.com/lmtani/learning-clean-architecture/internal/infra/graph"
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

	// Connect to database
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

	// Connect to RabbitMQ
	rabbitMQChannel := getRabbitMQChannel()
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	fmt.Println("Starting web server on port", conf.WebServerPort)
	webserver := server.NewWebServer(conf.WebServerPort)
	httpOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", httpOrderHandler.Create)
	go webserver.Start()

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", conf.GraphQLServerPort)
	http.ListenAndServe(":"+conf.GraphQLServerPort, nil)
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

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5"
	"github.com/lmtani/learning-clean-architecture/configs"
	"github.com/lmtani/learning-clean-architecture/internal/infra/database/psql"
	"github.com/lmtani/learning-clean-architecture/internal/infra/event/handler"
	"github.com/lmtani/learning-clean-architecture/internal/infra/graph"
	"github.com/lmtani/learning-clean-architecture/internal/infra/grpc/pb"
	"github.com/lmtani/learning-clean-architecture/internal/infra/grpc/service"
	"github.com/lmtani/learning-clean-architecture/internal/infra/web/server"
	"github.com/lmtani/learning-clean-architecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	conn := getPostgresConnection(ctx, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName))
	defer conn.Close(ctx)

	queries := psql.New(conn)

	// Connect to RabbitMQ
	rabbitMQChannel := getRabbitMQChannel(conf.RABBITMQHost)
	defer rabbitMQChannel.Close()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	// Start web server
	fmt.Println("Starting web server on port", conf.WebServerPort)
	webserver := server.NewWebServer(conf.WebServerPort)
	httpOrderHandler := NewWebOrderHandler(queries, eventDispatcher)
	webserver.AddHandler("POST /order", httpOrderHandler.Create)
	webserver.AddHandler("GET /order", httpOrderHandler.List)
	go webserver.Start()

	// Create use case
	createOrderUseCase := NewCreateOrderUseCase(queries, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(queries)

	// Start gRPC server
	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", conf.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	fmt.Println("Starting GraphQL server on port", conf.GraphQLServerPort)
	// Build your executable schema as before
	schema := graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			CreateOrderUseCase: *createOrderUseCase,
			ListOrdersUseCase:  *listOrdersUseCase,
		},
	})

	// Construct the handler directly
	srv := graphql_handler.New(schema)
	// Add the transports
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	// Add the extension
	srv.Use(extension.Introspection{})
	// Serve your GraphQL endpoint and playground
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	http.ListenAndServe(":"+conf.GraphQLServerPort, nil)
}

func getRabbitMQChannel(rabbitmqHost string) *amqp.Channel {
	var conn *amqp.Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitmqHost))
		if err == nil {
			break
		}
		fmt.Printf("Failed to connect to RabbitMQ, retrying in 2 seconds... (%d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

func getPostgresConnection(ctx context.Context, connStr string) *pgx.Conn {
	var conn *pgx.Conn
	var err error
	for i := 0; i < 5; i++ {
		conn, err = pgx.Connect(ctx, connStr)
		if err == nil {
			break
		}
		fmt.Printf("Failed to connect to PostgreSQL, retrying in 2 seconds... (%d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic(err)
	}
	return conn
}

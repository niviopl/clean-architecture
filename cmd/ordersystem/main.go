package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/niviopl/clean-architecture/internal/infra/config"
	"github.com/niviopl/clean-architecture/internal/infra/database"
	"github.com/niviopl/clean-architecture/internal/infra/graph"
	"github.com/niviopl/clean-architecture/internal/infra/graph/generated"
	"github.com/niviopl/clean-architecture/internal/infra/grpc/pb"
	grpcservice "github.com/niviopl/clean-architecture/internal/infra/grpc/service"
	"github.com/niviopl/clean-architecture/internal/infra/web/handlers"
	"github.com/niviopl/clean-architecture/internal/usecase"
)

func main() {
	cfg := config.Load()

	db := connectWithRetry(cfg)
	defer db.Close()

	runMigrations(db)

	orderRepository := database.NewOrderRepository(db)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	go startWebServer(cfg.WebServerPort, createOrderUseCase, listOrdersUseCase)
	go startGRPCServer(cfg.GRPCServerPort, createOrderUseCase, listOrdersUseCase)
	startGraphQLServer(cfg.GraphQLServerPort, createOrderUseCase, listOrdersUseCase)
}

// connectWithRetry waits for the database to be ready before proceeding,
// avoiding a startup race condition against the db container.
func connectWithRetry(cfg *config.Config) *sql.DB {
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?parseTime=true"

	db, err := sql.Open(cfg.DBDriver, dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	const maxAttempts = 30
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err = db.Ping(); err == nil {
			log.Println("database is ready")
			return db
		}
		log.Printf("waiting for database (attempt %d/%d): %v", attempt, maxAttempts, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("could not connect to database after %d attempts: %v", maxAttempts, err)
	return nil
}

func runMigrations(db *sql.DB) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	if err != nil {
		log.Fatalf("failed to create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("migrations applied successfully")
}

func startWebServer(port string, createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) {
	handler := handlers.NewOrderHandler(createOrderUseCase, listOrdersUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /order", handler.Create)
	mux.HandleFunc("GET /order", handler.List)

	log.Printf("web server running on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("web server error: %v", err)
	}
}

func startGRPCServer(port string, createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on grpc port: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderService := grpcservice.NewOrderService(createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	log.Printf("grpc server running on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc server error: %v", err)
	}
}

func startGraphQLServer(port string, createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}}))

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)

	log.Printf("graphql server running on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("graphql server error: %v", err)
	}
}

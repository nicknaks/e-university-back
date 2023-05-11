package main

import (
	"back/graph"
	"back/graph/generated"
	"back/internal/auth_service"
	"back/pkg/parser"
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq" // <------------ here
	"github.com/rs/cors"
	"net/http"
	"os"
)

const defaultPort = "8090"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3000"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	resolver := graph.NewResolver()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
		Directives: generated.DirectiveRoot{
			IsAuthenticated: resolver.Storage.IsAuth,
		},
	}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == "example.org"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Use(auth_service.InjectHTTPMiddleware())
	router.Handle("/", playground.Handler("Starwars", "/query"))
	router.Handle("/query", srv)

	data, err := parser.ParseFaculties(nil)
	if err != nil {
		panic(err)
	}

	data, err = parser.InitData(context.Background(), resolver.Storage, data)
	if err != nil {
		panic(err)
	}

	lessons, err := parser.ParseSchedule(nil, data)
	if err != nil {
		panic(err)
	}

	err = parser.ExtractTeachers(context.Background(), resolver.Storage, lessons)
	if err != nil {
		panic(err)
	}

	err = parser.InitUsers(context.Background(), resolver.Storage)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8090", router)
	if err != nil {
		panic(err)
	}
}

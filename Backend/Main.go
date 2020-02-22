package main

import (
	"Booker/Backend/db"
	"Booker/Backend/queue"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	log "github.com/sirupsen/logrus"

	"context"
	"net/http"
	"os"
)

func SetDatabaseContext(client *mongo.Client) (func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", client)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func init() {
	// Log as JSON
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout
	log.SetOutput(os.Stdout)

	// Only log warning severity or above
	log.SetLevel(log.InfoLevel)
}

func main() {
	r := chi.NewRouter()

	// Get a client for interacting with the database
	client, err := db.GetDBClient(MONGODB_CONN_INFO)
	if err != nil {
		log.Fatal("Could not connect to database")
	}
	log.Info("Connected to database")

	// Recovers from panics and returns an HTTP 500 status if possible
	r.Use(middleware.Recoverer)
	// Times out requests if they go on too long
	r.Use(middleware.Timeout(REQUEST_TIMEOUT))
	// Puts the database handle in request
	r.Use(SetDatabaseContext(client))

	// Mount the subrouters
	r.Mount(QUEUE_PATH, queue.QueueRouter())

	http.ListenAndServe(SERVER_PORT, r)
}
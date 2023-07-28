package server

import (
	"fmt"
	"grace/agent"
	"grace/database"
	"grace/handlers"
	"grace/repositories"
	"grace/response"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func StartRESTServer(db *gorm.DB) *http.Server {
	// Initialize repositories
	bookRepository := &repositories.BookRepository{DB: db}
	collectionRepository := &repositories.CollectionRepository{DB: db}
	collectionBookRepository := &repositories.CollectionBookRepository{DB: db}

	// Initialize handlers
	bookHandler := &handlers.BookHandler{Repository: bookRepository}
	collectionHandler := &handlers.CollectionHandler{
		Repository:         collectionRepository,
		BookRepository:     bookRepository,
		CollectionBookRepo: collectionBookRepository,
	}

	mux := mux.NewRouter()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = response.ResponseJson([]string{"/v1"}).Render(w)
	})

	var bookIdCmd = agent.APIEndpoint{
		Name:   "book with id",
		Path:   "book/{id}",
		Get:    agent.APIEndpointAction{Handler: bookHandler.GetBook},
		Put:    agent.APIEndpointAction{Handler: bookHandler.UpdateBook},
		Delete: agent.APIEndpointAction{Handler: bookHandler.DeleteBook},
		Patch:  agent.APIEndpointAction{Handler: bookHandler.PatchBook},
	}
	var bookCmd = agent.APIEndpoint{
		Name: "book",
		Path: "book",
		Get:  agent.APIEndpointAction{Handler: bookHandler.GetBooks},
		Post: agent.APIEndpointAction{Handler: bookHandler.CreateBook},
	}
	var collectionIdCmd = agent.APIEndpoint{
		Name:   "collection with id",
		Path:   "collection/{id}",
		Get:    agent.APIEndpointAction{Handler: collectionHandler.GetCollection},
		Put:    agent.APIEndpointAction{Handler: collectionHandler.UpdateCollection},
		Delete: agent.APIEndpointAction{Handler: collectionHandler.DeleteCollection},
	}
	var collectionCmd = agent.APIEndpoint{
		Name: "collection",
		Path: "collection",
		Post: agent.APIEndpointAction{Handler: collectionHandler.CreateCollection},
		Get:  agent.APIEndpointAction{Handler: collectionHandler.GetCollections},
	}
	var collectionBookCmd = agent.APIEndpoint{
		Name:   "collection book",
		Path:   "collection/{id}/book",
		Get:    agent.APIEndpointAction{Handler: collectionHandler.GetCollectionBooks},
		Post:   agent.APIEndpointAction{Handler: collectionHandler.AddBookToCollection},
		Delete: agent.APIEndpointAction{Handler: collectionHandler.RemoveBookFromCollection},
	}

	var api1 = []agent.APIEndpoint{
		bookIdCmd,
		bookCmd,
		collectionIdCmd,
		collectionCmd,
		collectionBookCmd,
	}

	for _, c := range api1 {
		agent.CreateCmd(mux, "v1", c)
	}
	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	return server
}

func SetupDatabase(db_name string) *gorm.DB {
	dsn := fmt.Sprintf("host=localhost user=app password=password dbname=%s port=5432 sslmode=disable", db_name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func MigrateDatabase(db *gorm.DB) {
	db.AutoMigrate(&database.Book{}, &database.Collection{}, &database.CollectionBook{})
}

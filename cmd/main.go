package main

import (
	"log"
	"os"

	"api-search/assert"
	"api-search/cmd/api"

	"github.com/joho/godotenv"
	"github.com/meilisearch/meilisearch-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	godotenv.Load()

	mongo := assert.Try(mongo.Connect(options.Client().ApplyURI(os.Getenv("MONGODB_URI"))))
	m := mongo.Database("api-search")

	sm := meilisearch.New(os.Getenv("MEILI_HOST"), meilisearch.WithAPIKey(os.Getenv("MEILI_KEY")))

	api := api.NewAPI(m, sm)

	log.Fatalln(api.Start(":8080"))
}

package apiserver

import (
	"context"
	"net/http"
	"order-service/internal/app/store/mongostore"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Start(config *Config) error {
	clientOptions := options.Client().ApplyURI(config.DatabaseURL)
	ctx := context.Background()
	defer ctx.Done()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	database := client.Database("restapi_dev")
	collection := database.Collection("order")

	store := mongostore.New(ctx, collection)
	//store := teststore.New()
	srv := newServer(store)

	return http.ListenAndServe(config.BindAddr, srv)
}

package config

import (
	"context"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

func getCredentials() option.ClientOption {
	return option.WithCredentialsFile("static/image-maesto-firebase-adminsdk-n7g64-6d04424972.json")
}

func initializeAppWithServiceAccount() *firebase.App {

	// Initialize app with background context
	app, err := firebase.NewApp(context.Background(), nil, getCredentials())
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing firebase app")
	}
	return app
}

func GetBucket(bucketName string) *storage.BucketHandle {
	client, err := storage.NewClient(context.Background(), getCredentials())

	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing firebase storage client")
	}

	bucket := client.Bucket(bucketName)

	return bucket
}

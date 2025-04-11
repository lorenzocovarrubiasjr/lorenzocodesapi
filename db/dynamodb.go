package db

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var DB *dynamodb.Client

func init() {
	// Explicitly set the region to match your DynamoDB table's region
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"), // Make sure this matches your table region
	)
	if err != nil {
		log.Printf("Warning: Failed to load AWS config: %v. DB will initialize on demand.", err)
		return
	}
	DB = dynamodb.NewFromConfig(cfg)
}

func GetClient() *dynamodb.Client {
	if DB == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("us-west-2"), // Make sure this matches your table region
		)
		if err != nil {
			log.Fatalf("Fatal: Cannot initialize DynamoDB client: %v", err)
		}
		DB = dynamodb.NewFromConfig(cfg)
	}
	return DB
}

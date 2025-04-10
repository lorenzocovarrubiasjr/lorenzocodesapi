package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var DB *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Failed to load AWS config: " + err.Error())
	}
	DB = dynamodb.NewFromConfig(cfg)
}

func GetClient() *dynamodb.Client {
	return DB
}

// func TableName() string {
// 	return "YourTableName" // Replace with your table name or make configurable
// }

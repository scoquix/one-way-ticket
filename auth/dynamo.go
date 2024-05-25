package auth

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"one-way-ticket/models"
	"os"
)

var (
	tableName = "sessions"
	svc       *dynamodb.DynamoDB
)

func init() {
	// Initialize a session that the SDK uses to load configuration,
	// credentials, and region from the environment
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := "us-east-1"
	token := ""

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, token),
	}))
	svc = dynamodb.New(sess)
}

// CreateSession creates a new session
func CreateSession(token string, ttl int64) error {
	sess := models.Session{
		Token: token,
		TTL:   ttl,
	}

	av, err := dynamodbattribute.MarshalMap(sess)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}
	return nil
}

// GetSessionForUser retrieves a session by sessionID
func GetSessionForUser(userID string) (*models.Session, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"token": {
				S: aws.String(userID),
			},
		},
	}

	result, err := svc.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	if result.Item == nil {
		return nil, nil // Session not found
	}

	var sess models.Session
	err = dynamodbattribute.UnmarshalMap(result.Item, &sess)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %v", err)
	}

	return &sess, nil
}

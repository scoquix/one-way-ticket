package routers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gin-gonic/gin"
	"one-way-ticket/auth"
	"one-way-ticket/users"
	"os"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	handler := auth.NewHandler(newDynamoClient())

	r.POST("/login", handler.Login)

	userRoutes := r.Group("/users")
	userRoutes.Use(auth.AuthenticateMiddleware())
	{
		userRoutes.GET("/", users.GetUsers)
		userRoutes.GET("/:id", users.GetUser)
		userRoutes.POST("/", users.CreateUser)
		userRoutes.PUT("/:id", users.UpdateUser)
		userRoutes.DELETE("/:id", users.DeleteUser)
	}

	return r
}

func newDynamoClient() dynamodbiface.DynamoDBAPI {
	// Initialize a session that the SDK uses to load configuration,
	// credentials, and region from the environment
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}
	token := ""

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, token),
	}))
	return dynamodb.New(sess)
}

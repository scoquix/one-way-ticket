package auth

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"one-way-ticket/mocks"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"one-way-ticket/models"
)

// Test CreateSession
func TestCreateSession(t *testing.T) {
	mockSvc := new(mocks.MockDynamoDBClient)
	mockSession := models.Session{
		Token: "test-token",
		TTL:   123456,
	}

	av, err := dynamodbattribute.MarshalMap(mockSession)
	assert.NoError(t, err)

	mockSvc.On("PutItem", &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}).Return(&dynamodb.PutItemOutput{}, nil)

	err = CreateSession(mockSvc, "test-token", 123456)
	assert.NoError(t, err)

	mockSvc.AssertExpectations(t)
}

// Test GetSessionForUser
func TestGetSessionForUser(t *testing.T) {
	mockSvc := new(mocks.MockDynamoDBClient)
	mockSession := models.Session{
		Token: "test-token",
		TTL:   123456,
	}

	av, err := dynamodbattribute.MarshalMap(mockSession)
	assert.NoError(t, err)

	mockSvc.On("GetItem", &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"token": {
				S: aws.String("test-token"),
			},
		},
	}).Return(&dynamodb.GetItemOutput{
		Item: av,
	}, nil)

	sess, err := GetSessionForUser(mockSvc, "test-token")
	assert.NoError(t, err)
	assert.Equal(t, &mockSession, sess)

	mockSvc.AssertExpectations(t)
}

// Test GetSessionForUser when session not found
func TestGetSessionForUser_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockDynamoDBClient)
	mockSvc.On("GetItem", &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"token": {
				S: aws.String("nonexistent-token"),
			},
		},
	}).Return(&dynamodb.GetItemOutput{
		Item: nil,
	}, nil)

	sess, err := GetSessionForUser(mockSvc, "nonexistent-token")
	assert.NoError(t, err)
	assert.Nil(t, sess)

	mockSvc.AssertExpectations(t)
}

// Test GetSessionForUser with error
func TestGetSessionForUser_Error(t *testing.T) {
	mockSvc := new(mocks.MockDynamoDBClient)
	mockSvc.On("GetItem", &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"token": {
				S: aws.String("error-token"),
			},
		},
	}).Return(&dynamodb.GetItemOutput{}, errors.New("dynamodb error"))

	sess, err := GetSessionForUser(mockSvc, "error-token")
	assert.Error(t, err)
	assert.Nil(t, sess)

	mockSvc.AssertExpectations(t)
}

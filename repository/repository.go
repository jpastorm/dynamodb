package repository

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Post struct {
	ID    int64  `json:"id,omitempty"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type dynamoDBRepo struct {
	tableName string
}

func NewDynamoDBRepository() PostRepository {
	return &dynamoDBRepo{
		tableName: "posts",
	}
}

func createDynamoDBClient() *dynamodb.DynamoDB {
	// Create AWS Session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	if err != nil {
		log.Println(err)
	}
	// Return DynamoDB client
	return dynamodb.New(sess)
}

func (repo *dynamoDBRepo) Save(post *Post) (*Post, error) {
	// Get a new DynamoDB client
	dynamoDBClient := createDynamoDBClient()

	// Transforms the post to map[string]*dynamodb.AttributeValue
	attributeValue, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		return nil, err
	}

	// Create the Item Input
	item := &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(repo.tableName),
	}

	// Save the Item into DynamoDB
	_, err = dynamoDBClient.PutItem(item)
	if err != nil {
		return nil, err
	}

	return post, err
}

func (repo *dynamoDBRepo) FindAll() ([]Post, error) {
	// Get a new DynamoDB client
	dynamoDBClient := createDynamoDBClient()

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName: aws.String(repo.tableName),
	}

	// Make the DynamoDB Query API call
	result, err := dynamoDBClient.Scan(params)
	if err != nil {
		return nil, err
	}
	var posts []Post = []Post{}
	for _, i := range result.Items {
		post := Post{}

		err = dynamodbattribute.UnmarshalMap(i, &post)

		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *dynamoDBRepo) FindByID(id string) (*Post, error) {
	// Get a new DynamoDB client
	dynamoDBClient := createDynamoDBClient()

	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	post := Post{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &post)
	if err != nil {
		panic(err)
	}
	return &post, nil
}

// Delete: TODO
func (repo *dynamoDBRepo) Delete(post *Post) error {
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	guuid "github.com/google/uuid"
	"github.com/guregu/dynamo"
)

var (
	region    = os.Getenv("aws_region")
	tableName = os.Getenv("dynamo_table")
	table     dynamo.Table
)

type SystemGCVM struct {
	UUID           guuid.UUID `json:"uuid"`
	UserId         string     `json:"userid"`
	VmName         string     `json:"vmname"`
	VmType         string     `json:"vmtype"`
	InstanceId     string     `json:"instanceid,omitempty" ,dynamo:",omitempty"`
	IpAddress      string     `json:"ipaddress,omitempty" ,dynamo:",omitempty"`
	SanCleanUp     string     `json:"sancleanup,omitempty" ,dynamo:",omitempty"`
	AgentInstalled string     `json:"agent,omitempty" ,dynamo:",omitempty"`
	SshKeyPath     string     `json:"keypath,omitempty" ,dynamo:",omitempty"`
	SanPath        string     `json:"sanpath,omitempty" ,dynamo:",omitempty"`
	CreationDate   time.Time  `json:"creation_date,omitempty" ,dynamo:",omitempty"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var vmList []SystemGCVM
	err1 := table.Scan().All(&vmList)
	if err1 != nil {
		fmt.Println("Error in fetching data from dynamoDB", err1.Error())
		return events.APIGatewayProxyResponse{Body: "Error in fetching data from dynamoDB", StatusCode: 500}, err1
	}
	jsonString, err := json.Marshal(vmList)
	if err != nil {
		fmt.Println("MarshalError", err.Error())
		return response(err.Error(), 400)
	}

	return response(string(jsonString), 200)
}

func response(body string, statusCode int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		}}, nil
}

func connectDb() {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	table = db.Table(tableName)
}

func main() {
	connectDb()
	lambda.Start(Handler)
}

package main

import (
	"fmt"
	"os"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	guuid "github.com/google/uuid"
	"github.com/guregu/dynamo"
	log "github.com/sirupsen/logrus"
)

var (
	region    = os.Getenv("aws_region")
	tableName = os.Getenv("dynamo_table")
	table     dynamo.Table
)

type SystemGCVM struct {
	UUID                     guuid.UUID `json:"uuid"`
	UserId                   string     `json:"userid"`
	VmName                   string     `json:"vmname"`
	VmType                   string     `json:"vmtype"`
	InstanceId               string     `json:"instanceid,omitempty" ,dynamo:",omitempty"`
	IpAddress                string     `json:"ipaddress,omitempty" ,dynamo:",omitempty"`
	SanCleanUp               string     `json:"sancleanup,omitempty" ,dynamo:",omitempty"`
	AgentInstalled           string     `json:"agent,omitempty" ,dynamo:",omitempty"`
	SshKeyPath               string     `json:"keypath,omitempty" ,dynamo:",omitempty"`
	CreationDate             time.Time  `json:"creation_date,omitempty" ,dynamo:",omitempty"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestVM SystemGCVM
	var retunData SystemGCVM
	var vmList []SystemGCVM
	err := json.Unmarshal([]byte(request.Body), &requestVM)
	if err != nil {
		log.Errorf("Unmarshal body %s: %s", request.Body, err)
		return response(err.Error(), 400)
	}
	err = table.Scan().All(&vmList)
	if err != nil {
		fmt.Println("Error in fetching data from dynamoDB", err.Error())
		return events.APIGatewayProxyResponse{Body: "Error in fetching data from dynamoDB", StatusCode: 500}, err
	}
	for _,vm := range vmList {
		if vm.InstanceId == requestVM.InstanceId{
			return response(string("Given InstanceId already exists in schedules."), 500)
		}
	}
	
	retunData, err = createEntry(requestVM)
	if err != nil {
		log.Errorf("CreateEntry %s: %s", requestVM.InstanceId, err)
		return response(err.Error(), 400)
	}
    jsonString, err := json.Marshal(retunData)
	if err != nil {
		fmt.Println("MarshalError", err.Error())
		return response(err.Error(), 400)
	}

    return response(string(jsonString), 200)
}

func response(body string, statusCode int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin": "*",
		}}, nil
}

func createEntry(data SystemGCVM) (SystemGCVM, error) {
	// set insert date
	data.CreationDate = time.Now()
	data.UUID = guuid.New()
	err := table.Put(data).Run()
	if err != nil {
		return data, err
	}
	return getEntry(data)
}

func getEntry(data SystemGCVM) (SystemGCVM, error) {
	var result SystemGCVM
	err := table.Get("UUID", data.UUID).One(&result)
	return result, err
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

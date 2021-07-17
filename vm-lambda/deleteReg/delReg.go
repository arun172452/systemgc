package main

import (
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	guuid "github.com/google/uuid"
	"github.com/guregu/dynamo"
	log "github.com/sirupsen/logrus"
)

var (
	region    = os.Getenv("aws_region")
	tableName = os.Getenv("dynamo_table")
	table     dynamo.Table
	snsTopic  = os.Getenv("sns_topic")
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
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	table = db.Table(tableName)
	key, err := guuid.Parse(request.PathParameters["uuid"])
	if err != nil {
		log.Fatalf("Got error while Parsing UUID: %s", err)
	}
	log.Printf("Deleting UUID:", key)
	data, err := getEntryForGivenUuid(key)
	if err != nil {
		log.Fatalf("Got error while fetching Data: %s", err)
	}
	if data.SshKeyPath != "" {
		snssvc := sns.New(sess)
		log.Printf("Topic of SNS:", &snsTopic)
		sjson := string(data.SshKeyPath)
		log.Printf("Sending delete file from s3 event for file: %s", data.SshKeyPath)
		_, err = snssvc.Publish(&sns.PublishInput{
			Message:  &sjson,
			TopicArn: &snsTopic,
		})
		if err != nil {
			log.Errorln(err.Error())
		}

	}
	err = table.Delete("UUID", key).Run()
	if err != nil {
		log.Fatalf("Got error calling DeleteItem: %s", err)
	}
	return response(string("Deleted UUID.."), 200)
}

func response(body string, statusCode int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		}}, nil
}
func getEntryForGivenUuid(key guuid.UUID) (SystemGCVM, error) {
	var result SystemGCVM
	err := table.Get("UUID", key).One(&result)
	return result, err
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	lambda.Start(Handler)
}

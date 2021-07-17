package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	guuid "github.com/google/uuid"
	"github.com/guregu/dynamo"
	log "github.com/sirupsen/logrus"
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
	CreationDate   time.Time  `json:"creation_date,omitempty" ,dynamo:",omitempty"`
}

var (
	tableName = os.Getenv("dynamo_table")
	table     dynamo.Table
	snsTopic  = os.Getenv("ec2_sns_topic")
	region    = os.Getenv("aws_region")
)

func main() {
	connectDb()
	lambda.Start(Handler)
}

func connectDb() {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	table = db.Table(tableName)
}

// Handler is the main lambda function
func Handler(ctx context.Context) {
	var vmList []SystemGCVM
	err := table.Scan().All(&vmList)
	if err != nil {
		log.Errorf("Error in fetching data from dynamoDB: %s", err)
		return
	}
	for _, vm := range vmList {
		if vm.VmType == "aws" {
			log.Infof("sending ec2 stop event for: %s via sns", vm.VmName)
			sendSNS(vm)
		}

	}
}

func sendSNS(data SystemGCVM) {
	json, err := json.Marshal(data)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	sjson := string(json)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sns.New(sess)

	_, err = svc.Publish(&sns.PublishInput{
		Message:  &sjson,
		TopicArn: &snsTopic,
	})
	if err != nil {
		log.Errorln(err.Error())
	}
}

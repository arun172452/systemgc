package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	guuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	region = os.Getenv("aws_region")
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

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	ec2svc := ec2.New(session.New(&aws.Config{
		Region: aws.String(region)}))
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{},
	}
	for _, record := range snsEvent.Records {
		var u SystemGCVM
		err := json.Unmarshal([]byte(record.SNS.Message), &u)
		if err != nil {
			log.Errorf("ERROR Unmarshal Body %s: %s", record.SNS.Message, err)
		}

		input.InstanceIds = append(input.InstanceIds, aws.String(u.InstanceId))
	}

	result, err := ec2svc.StopInstances(input)
	if err == nil {
		log.Println(result)

	} else {
		log.Println(err.Error())
	}
}

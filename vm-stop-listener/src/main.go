package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	guuid "github.com/google/uuid"
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

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, snsEvent events.SNSEvent) {

	for _, record := range snsEvent.Records {
		var u SystemGCVM
		err := json.Unmarshal([]byte(record.SNS.Message), &u)
		if err != nil {
			log.Errorf("ERROR Unmarshal Body %s: %s", record.SNS.Message, err)
		}

		log.Println("got sns message")
		log.Printf("%+v\n", u)
		//to implement API call to stop onpremise VM In Future
	}
}

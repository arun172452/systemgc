package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	guuid "github.com/google/uuid"
	"github.com/guregu/dynamo"
)

var (
	region  = os.Getenv("region")
	vmTable = os.Getenv("table_vm")
	dbTable = os.Getenv("table_db")
	table   dynamo.Table
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

type DbRegistration struct {
	UUID         guuid.UUID `json:"uuid"`
	UserId       string     `json:"userid"`
	DbType       string     `json:"dbtype"`
	DbUrl        string     `json:"dburl"`
	DbUser       string     `json:"dbuser"`
	DbScript     string     `json:"dbscript"`
	TableRegex   string     `json:"tableregex,omitempty" ,dynamo:",omitempty"`
	ScriptName   string     `json:"filename,omitempty" ,dynamo:",omitempty"`
	CreationDate time.Time  `json:"creation_date,omitempty" ,dynamo:",omitempty"`
}

func Handler(ctx context.Context, S3Event events.S3Event) {
	fmt.Println("Bucket Name", S3Event.Records[0].S3.Bucket.Name)
	fmt.Println("File Name", S3Event.Records[0].S3.Object.Key)
	svc := s3.New(session.New())
	input := &s3.HeadObjectInput{
		Bucket: aws.String(S3Event.Records[0].S3.Bucket.Name),
		Key:    aws.String(S3Event.Records[0].S3.Object.Key),
	}
	result, err := svc.HeadObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	uploadedFor := DerefString(result.Metadata["Type"])
	uuid := DerefString(result.Metadata["Uuid"])
	if uploadedFor == "db" {
		connectDb(dbTable)
		//Get entry
		var cps []DbRegistration
		table.Scan().Filter("$ = ?", "UUID", uuid).All(&cps)
		if len(cps) != 0 {
			var p DbRegistration
			err := table.Get("UUID", cps[0].UUID).One(&p)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			p.ScriptName = S3Event.Records[0].S3.Object.Key

			err = table.Put(p).Run()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	if uploadedFor == "vm" {
		connectDb(vmTable)
		//Get entry
		var cps []SystemGCVM
		table.Scan().Filter("$ = ?", "UUID", uuid).All(&cps)
		if len(cps) != 0 {
			var p SystemGCVM
			err := table.Get("UUID", cps[0].UUID).One(&p)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			p.SshKeyPath = S3Event.Records[0].S3.Object.Key

			err = table.Put(p).Run()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func connectDb(tableName string) {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	table = db.Table(tableName)
}

func DerefString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func main() {
	lambda.Start(Handler)
}

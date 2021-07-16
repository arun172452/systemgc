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
	"github.com/aws/aws-sdk-go/service/ec2"
	guuid "github.com/google/uuid"
	"github.com/guregu/dynamo"
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
	CreationDate            time.Time  `json:"creation_date,omitempty" ,dynamo:",omitempty"`
}

type instanceInfo struct {
    InstanceId       string `json:"instanceId"`
	Name             string  `json:"name"`
}
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var vmList []SystemGCVM
	var instanceList []string
	err1 := table.Scan().All(&vmList)
	if err1 != nil {
		fmt.Println("Error in fetching data from dynamoDB", err1.Error())
		return events.APIGatewayProxyResponse{Body: "Error in fetching data from dynamoDB", StatusCode: 500}, err1
	}
	for _,vm := range vmList {
		instanceList = append(instanceList,vm.InstanceId)
	}
	//Type of vm - aws/onpremise
	environment := request.PathParameters["name"]
	var instanceSummary []instanceInfo
	if(environment == "onpremise"){
		//To call onpremise API with appropriate security key provided by respective team
		vmone := instanceInfo{InstanceId:"192.168.204.31",Name:"ukvm01"}
		vmtwo := instanceInfo{InstanceId:"192.168.204.219",Name:"ukvm219"}
		vmthree := instanceInfo{InstanceId:"192.168.204.217",Name:"ukvm217"}
		instanceSummary = append(instanceSummary,vmone)
		instanceSummary = append(instanceSummary,vmtwo)
		instanceSummary = append(instanceSummary,vmthree)
	} else {
		ec2svc := ec2.New(session.New(&aws.Config{
			Region: aws.String(region)}))
		filters := &ec2.DescribeInstancesInput{}
		resp, err := ec2svc.DescribeInstances(filters)
		if err != nil {
			fmt.Println("there was an error listing instances in", err.Error())
			return events.APIGatewayProxyResponse{Body: "there was an error listing instances", StatusCode: 500}, err
		}
		for idx := range resp.Reservations {
			nonStandard := "no"
			tagName := "NoName"
			for _, inst := range resp.Reservations[idx].Instances {
				if(*inst.State.Code == 48){
					continue
				}
				if(contains(instanceList,*inst.InstanceId)==true){
					continue
				}
				for _, tag := range inst.Tags {
					if (*tag.Key == "Schedule"){
						nonStandard = "yes"
					}
					if (*tag.Key == "Name") || (*tag.Key == "name") {
						tagName = *tag.Value
					}
				}
				if(nonStandard == "no"){
					instanceDe := instanceInfo{InstanceId:*inst.InstanceId,Name:tagName}
					instanceSummary = append(instanceSummary,instanceDe)
				}
			
			}
		}
	}
    jsonString, err := json.Marshal(instanceSummary)
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

func connectDb() {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	table = db.Table(tableName)
}

func contains(instanceList []string, instanceId string) bool {
    for _, a := range instanceList {
        if a == instanceId {
            return true
        }
    }
    return false
}

func main() {
	connectDb()
	lambda.Start(Handler)
}

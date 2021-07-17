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
	"github.com/aws/aws-sdk-go/service/s3"
	guuid "github.com/google/uuid"
	"github.com/guregu/dynamo"
	log "github.com/sirupsen/logrus"
)

var (
	region     = os.Getenv("aws_region")
	tableName  = os.Getenv("dynamo_table")
	bucketName = os.Getenv("s3_bucket_name")
	table      dynamo.Table
)

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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var dbDetails DbRegistration
	var retunData DbRegistration
	var dbList []DbRegistration
	err := json.Unmarshal([]byte(request.Body), &dbDetails)
	if err != nil {
		log.Errorf("Unmarshal body %s: %s", request.Body, err)
		return response(err.Error(), 400)
	}
	err = table.Scan().All(&dbList)
	if err != nil {
		fmt.Println("Error in fetching data from dynamoDB", err.Error())
		return events.APIGatewayProxyResponse{Body: "Error in fetching data from dynamoDB", StatusCode: 500}, err
	}
	for _, db := range dbList {
		if db.DbUrl == dbDetails.DbUrl {
			return response(string("Given DB Url already exists in schedules."), 500)
		}
	}

	retunData, err = createEntry(dbDetails)
	if err != nil {
		log.Errorf("CreateEntry %s: %s", dbDetails.DbUrl, err)
		return response(err.Error(), 400)
	}

	if retunData.DbScript == "yes" {
		url, err := getSignedUrlForSshKeyUpload(retunData)
		if err != nil {
			log.Errorf("Error Creating Signed Url %s: %s", dbDetails.DbUrl, err)
		} else {
			retunData.ScriptName = url
		}
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
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		}}, nil
}

func getSignedUrlForSshKeyUpload(data DbRegistration) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		fmt.Println("MarshalError", err.Error())
	}
	svc := s3.New(sess)
	fileNAme := data.UUID.String() + ".sql"
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileNAme),
	})
	query := req.HTTPRequest.URL.Query()
	query.Add("x-amz-meta-type", "db")
	query.Add("x-amz-meta-uuid", data.UUID.String())
	req.HTTPRequest.URL.RawQuery = query.Encode()
	str, err := req.Presign(15 * time.Minute)
	return str, err
}
func createEntry(data DbRegistration) (DbRegistration, error) {
	// set insert date
	data.CreationDate = time.Now()
	data.UUID = guuid.New()
	err := table.Put(data).Run()
	if err != nil {
		return data, err
	}
	return getEntry(data)
}

func getEntry(data DbRegistration) (DbRegistration, error) {
	var result DbRegistration
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

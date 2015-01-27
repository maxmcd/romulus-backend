package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/stripe/aws-go/aws"
	"github.com/stripe/aws-go/gen/lambda"
	"github.com/stripe/aws-go/gen/s3"
	"io/ioutil"
	"net/http"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	accessKey := "AKIAJNQEC3YKVDYKKXXQ"
	secretKey := "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"
	parseApplicationID := "TfJHzuJVZYU97rJc02JrJ8jy8JtsDNe1tbqACmJh"
	parseRestApiKey := "1myLlwt5YOWBWs3uNGQnIn71BymgzaPxmFxH1bIm"

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.parse.com/1/users/me", nil)
	handle(err)
	req.Header.Add("X-Parse-Application-Id", parseApplicationID)
	req.Header.Add("X-Parse-REST-API-Key", parseRestApiKey)
	req.Header.Add("X-Parse-Session-Token", "5a7FLdW20cjmUFV64Nijbf0yG")
	response, err := client.Do(req)
	handle(err)
	contents, err := ioutil.ReadAll(response.Body)
	handle(err)
	response.Body.Close()
	fmt.Println(contents)

	// http://stackoverflow.com/questions/11066946/partly-json-unmarshal-into-a-map-in-go
	var objmap map[string]*json.RawMessage
	json.Unmarshal(contents, &objmap)
	var username string
	err = json.Unmarshal(*objmap["username"], &username)

	fmt.Println(username)

	creds := aws.Creds(accessKey, secretKey, "")

	testLambda(creds)

	s3Cli := s3.New(creds, "us-west-2", nil)
	s3.PutObjectRequest{nil, "hi"}
}

// {
//     "app_id":"TfJHzuJVZYU97rJc02JrJ8jy8JtsDNe1tbqACmJh",
//     "createdAt":"2015-01-13T20:31:59.131Z",
//     "js_key":"PXp8V500uu6cMrUeAbuMBXuT1Kev872A7XvET1uG",
//     "objectId":"co6deoN38h",
//     "sessionToken":"5a7FLdW20cjmUFV64Nijbf0yG",
//     "updatedAt":"2015-01-13T21:23:14.626Z",
//     "username":"romulus"
// }

func testLambda(creds aws.CredentialsProvider) {
	lambda_cli := lambda.New(creds, "us-west-2", nil)
	request := &lambda.ListFunctionsRequest{nil, nil}
	resp, err := lambda_cli.ListFunctions(request)
	if err != nil {
		// fmt.Println(err)
	}
	fmt.Println(len(resp.Functions))
	fmt.Println(*resp.Functions[0].Description)
}

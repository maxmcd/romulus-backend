package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	fmt.Println("Serving on port 8080")
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/lambda/new", lambdaNewHandler)
	http.HandleFunc("/lambda/trigger", lambdaTriggerHandler)
	http.ListenAndServe(":8080", nil)

}

func lambdaTriggerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case "POST":
		body, header := lambdaTriggerPostReponse(r)
		w.WriteHeader(header)
		w.Write(body)
	default:
		body, header := defaultResponse(r)
		w.WriteHeader(header)
		w.Write(body)
	}
}

func lambdaNewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case "POST":
		body, header := lambdaNewPostReponse(r)
		w.WriteHeader(header)
		w.Write(body)
	default:
		body, header := defaultResponse(r)
		w.WriteHeader(header)
		w.Write(body)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case "POST":
		body, header := uploadPostReponse(r)
		w.WriteHeader(header)
		w.Write(body)
	default:
		body, header := defaultResponse(r)
		w.WriteHeader(header)
		w.Write(body)
	}
}

func lambdaTriggerPostReponse(r *http.Request) (body []byte, header int) {
	header = http.StatusInternalServerError
	body = []byte{}

	// get some kind of uid for the lambda function
	// pass it to amazon to trigger the function
	// pass along any parameters as well
	// work for parameters should be on the users end
	return
}

func lambdaNewPostReponse(r *http.Request) (body []byte, header int) {
	header = http.StatusInternalServerError
	body = []byte{}

	reader, err := r.MultipartReader()
	if err != nil {
		body = []byte(err.Error())
		return
	}

	form, err := reader.ReadForm(20000000) //20mb max allocated memory
	if err != nil {
		body = []byte(err.Error())
		return
	}
	file := form.File["body"][0]
	_ = file
	// create zip, including any dependencies
	// upload to lambda
	return
}

func defaultResponse(r *http.Request) (body []byte, header int) {
	header = http.StatusOK
	body = []byte("Deadly when I play a dope melody\nAnything less than the best is a felony")
	return
}
func uploadPostReponse(r *http.Request) (body []byte, header int) {
	header = http.StatusInternalServerError
	body = []byte{}
	// body, err := ioutil.ReadAll(r.Body)
	// handle(err)

	reader, err := r.MultipartReader()
	if err != nil {
		body = []byte(err.Error())
		return
	}

	form, err := reader.ReadForm(20000000) //20mb max allocated memory
	if err != nil {
		body = []byte(err.Error())
		return
	}
	// keys := [5]string{"bucket", "body", "key", "contentType", "sessionToken"}
	filebody, err := form.File["body"][0].Open()
	if err != nil {
		body = []byte(err.Error())
		return
	}
	contents, err := ioutil.ReadAll(filebody)
	if err != nil {
		body = []byte(err.Error())
		return
	}
	key := form.Value["key"][0]
	contentType := form.Value["contentType"][0]
	sessionToken := form.Value["sessionToken"][0]

	username, err := getParseUsernameFromSession(sessionToken)
	if err != nil {
		body = []byte(err.Error())
		return
	}
	// this should probably return an error

	if username != "" {
		fmt.Println(username)
		err := uploadS3File(username, contents, contentType, key)
		if err != nil {
			body = []byte(err.Error())
			return
		}
		header = http.StatusOK
	} else {
		header = http.StatusForbidden
	}
	return
}

func getParseUsernameFromSession(session string) (username string, err error) {
	parseApplicationID := "TfJHzuJVZYU97rJc02JrJ8jy8JtsDNe1tbqACmJh"
	parseRestApiKey := "1myLlwt5YOWBWs3uNGQnIn71BymgzaPxmFxH1bIm"

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.parse.com/1/users/me", nil)
	if err != nil {
		return
	}
	req.Header.Add("X-Parse-Application-Id", parseApplicationID)
	req.Header.Add("X-Parse-REST-API-Key", parseRestApiKey)
	req.Header.Add("X-Parse-Session-Token", session)
	response, err := client.Do(req)
	if err != nil {
		return
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	response.Body.Close()
	fmt.Println(string(contents))

	var objmap map[string]*json.RawMessage
	json.Unmarshal(contents, &objmap)

	if objmap["username"] != nil {
		err = json.Unmarshal(*objmap["username"], &username)
		return
	} else {
		return
	}

}

func uploadS3File(username string, contents []byte, fileType string, key string) (err error) {
	// grabs auth values from env variables
	auth, err := aws.EnvAuth()
	if err != nil {
		return
	}
	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket("romulus-host")
	err = bucket.Put(username+"/"+key, contents, fileType, s3.PublicRead)
	return
}

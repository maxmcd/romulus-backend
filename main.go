package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"io/ioutil"
	"net/http"
	"os"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Serving on port 8080")
	http.HandleFunc("/", uploadHandler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	switch r.Method {
	case "POST":
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

		testS3(username)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func testS3(username string) {
	// grabs auth values from env variables
	auth, err := aws.EnvAuth()
	handle(err)
	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket("romulus-host")
	err = bucket.Put(username+"/test.hi", []byte("content"), "text/html", s3.PublicRead)
	handle(err)
	fmt.Println(err)
}

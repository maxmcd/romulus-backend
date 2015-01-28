package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"io/ioutil"
	"net/http"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Serving on port 8080")
	http.HandleFunc("/", uploadHandler)
	http.ListenAndServe(":8080", nil)

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case "POST":
		// body, err := ioutil.ReadAll(r.Body)
		// handle(err)

		reader, err := r.MultipartReader()
		handle(err)

		form, err := reader.ReadForm(20000000) //20mb max allocated memory
		handle(err)

		// keys := [5]string{"bucket", "body", "key", "contentType", "sessionToken"}

		body, err := form.File["body"][0].Open()
		handle(err)
		contents, err := ioutil.ReadAll(body)
		handle(err)

		key := form.Value["key"][0]
		contentType := form.Value["contentType"][0]
		sessionToken := form.Value["sessionToken"][0]

		username := getParseUsernameFromSession(sessionToken)
		// this should probably return an error

		if username != "" {
			fmt.Println(username)
			uploadS3File(username, contents, contentType, key)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	default:
		w.Write([]byte("Deadly when I play a dope melody\nAnything less than the best is a felony"))
	}
}

func getParseUsernameFromSession(session string) (username string) {
	parseApplicationID := "TfJHzuJVZYU97rJc02JrJ8jy8JtsDNe1tbqACmJh"
	parseRestApiKey := "1myLlwt5YOWBWs3uNGQnIn71BymgzaPxmFxH1bIm"

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.parse.com/1/users/me", nil)
	handle(err)
	req.Header.Add("X-Parse-Application-Id", parseApplicationID)
	req.Header.Add("X-Parse-REST-API-Key", parseRestApiKey)
	req.Header.Add("X-Parse-Session-Token", session)
	response, err := client.Do(req)
	handle(err)
	contents, err := ioutil.ReadAll(response.Body)
	handle(err)
	response.Body.Close()
	fmt.Println(string(contents))

	var objmap map[string]*json.RawMessage
	json.Unmarshal(contents, &objmap)

	if objmap["username"] != nil {
		var username string
		err = json.Unmarshal(*objmap["username"], &username)

		return username
	} else {
		return ""
	}

}

func uploadS3File(username string, contents []byte, fileType string, key string) {
	// grabs auth values from env variables
	auth, err := aws.EnvAuth()
	handle(err)
	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket("romulus-host")
	err = bucket.Put(username+"/"+key, contents, fileType, s3.PublicRead)
	handle(err)
	fmt.Println(err)
}

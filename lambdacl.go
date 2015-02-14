package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"time"
)

func main() {

	// code := `
	//    console.log('Loading event');
	//    exports.handler = function(event, context) {
	//      console.log("value1 = " + event.key1);
	//      console.log("value2 = " + event.key2);
	//      console.log("value3 = " + event.key3);
	//      context.done(null, "Hello World");  // SUCCESS with message
	//    }`

	// err := create("helloworld3", code)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	err := run("helloworld3")
	if err != nil {
		panic(err)
	}
	getStats("helloworld3")
}

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz1234567890-")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// {
//    "key1":"value1",
//    "key2":"value2",
//    "key3":"value3"
// }

func create(functionName, code string) (err error) {
	filename, err := createZip(code)
	if err != nil {
		return
	}
	cmd := exec.Command(
		"aws",
		"lambda",
		"upload-function",
		"--region", "us-east-1",
		"--function-name", functionName,
		"--function-zip", filename,
		"--role", "arn:aws:iam::651778473396:role/lambda_s3_role",
		"--mode", "event",
		"--handler", "index.handler",
		"--runtime", "nodejs",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// cmd.Run()

	// cmd.Wait()
	return
}

func createZip(code string) (filename string, err error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	var files = []struct {
		Name, Body string
	}{
		{"index.js", code},
	}

	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			return filename, err
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			return filename, err
		}
	}

	err = w.Close()
	if err != nil {
		return filename, err
	}

	filename = randSeq(10) + ".zip"
	err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		return
	}

	return
}

func run(functionName string) (err error) {
	cmd := exec.Command(
		"aws",
		"lambda",
		"invoke-async",
		"--region", "us-east-1",
		"--function-name", functionName,
		"--invoke-args", "inputfile.txt",
		"--debug",
	)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// cmd.Stdin = os.Stdin
	re := regexp.MustCompile("'x-amzn-requestid': '(.*?)'")

	var b bytes.Buffer
	var e bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &e
	err = cmd.Run()
	debugOutput := e.Bytes()
	if err != nil {
		return
	}
	resultString := re.FindStringSubmatch(string(debugOutput))
	if len(resultString) > 1 {
		fmt.Println(resultString[1])
	} else {
		err = fmt.Errorf("Request id not found")
	}
	fmt.Println(string(b.Bytes()))

	cmd.Run()

	cmd.Wait()
	return
}

func getStats(functionName string) {
	cmd := exec.Command(
		"aws",
		"logs",
		"describe-log-streams",
		"--region", "us-east-1",
		"--log-group-name", "/aws/lambda/"+functionName,
	)

	time.Sleep(time.Second * 3)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// cmd.Stdin = os.Stdin
	jsonOutput, err := cmd.CombinedOutput()
	// fmt.Println(string(jsonOutput))
	// cmd.Run()

	// cmd.Wait()
	f := make(map[string][]map[string]interface{})
	err = json.Unmarshal(jsonOutput, &f)

	var largestTimestamp float64
	largestTimestamp = 0
	largestIndex := 0
	for index, value := range f["logStreams"] {
		if value["lastEventTimestamp"].(float64) > largestTimestamp {
			largestTimestamp = value["lastEventTimestamp"].(float64)
			largestIndex = index
		}
	}
	logStreamName := f["logStreams"][largestIndex]["logStreamName"].(string)

	// panic(err)
	if err != nil {
		panic(err)
	}

	// aws logs get-log-events --log-group-name /aws/lambda/helloworld3 --log-stream-name 0aee4e0328f341599460d14a2a1e6b53 --region us-east-1
	cmd = exec.Command(
		"aws",
		"logs",
		"get-log-events",
		"--region", "us-east-1",
		"--log-stream-name", logStreamName,
		"--log-group-name", "/aws/lambda/"+functionName,
	)
	newOut, err := cmd.Output()
	fmt.Println(string(newOut))
}

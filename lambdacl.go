package main

import (
	"os"
	"os/exec"
)

func main() {

	code := `
    console.log('Loading event');
    exports.handler = function(event, context) {
      console.log("value1 = " + event.key1);
      console.log("value2 = " + event.key2);
      console.log("value3 = " + event.key3);
      context.done(null, "Hello World");  // SUCCESS with message
    }`

	create("helloworld2", code)

	run()
}

// {
//    "key1":"value1",
//    "key2":"value2",
//    "key3":"value3"
// }

func create(functionName, code string) (err error) {
	filename := createZip(code)
	cmd := exec.Command(
		"aws",
		"lambda",
		"upload-function",
		"--region", "us-east-1",
		"--function-name", functionName,
		"--function-zip", "file.zip",
		"--role", "arn:aws:iam::651778473396:role/lambda_s3_role",
		"--mode", "event",
		"--handler", "helloworld.handler",
		"--runtime", "nodejs",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()

	cmd.Wait()
}

func createZip(code string) (filename string) {

}

func run() {
	cmd := exec.Command(
		"aws",
		"lambda",
		"invoke-async",
		"--region", "us-east-1",
		"--function-name", "helloworld",
		"--invoke-args", "inputfile.txt",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()

	cmd.Wait()
}

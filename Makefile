default: package 

package:
	GOOS=linux go build -o main
	zip lambda.zip main

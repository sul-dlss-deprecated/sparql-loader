default: package 

package:
	GOOS=linux go build -o neptune cmd/neptune/main.go
	zip lambda.zip neptune

deps:
	go get
run:
	go run main.go comm.go node.go token.go log.go helpers.go
build:
	go build -o bin/main main.go comm.go node.go token.go log.go helpers.go
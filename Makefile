.PHONY: build clean_lambda

build:
	env GOOS=linux go build -ldflags="-s -w" -o dist/lambda/share src/lambda/share/main.go
	env GOOS=linux go build -ldflags="-s -w" -o dist/lambda/create src/lambda/create/main.go

clean_lambda:
	rm -rf ./dist/lambda
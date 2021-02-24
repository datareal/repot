.PHONY: build clean_lambda deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o dist/lambda/share src/lambda/share/main.go

clean_lambda:
	rm -rf ./dist/lambda

deploy: clean_lambda build
	sls deploy --verbose
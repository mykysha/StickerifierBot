gomod-download:
	go get -u github.com/aws/aws-lambda-go
	go get -u github.com/disintegration/imaging
	go get -u gopkg.in/telebot.v3

gomod-tidy:
	go mod tidy -go=1.19 -compat=1.19

gomod-update:
	make gomod-download
	make gomod-tidy

build:
	make gomod-update
	env GOOS=linux go build -o bin/webhook cmd/main.go

aws-deploy:
	make build
	serverless deploy --verbose

aws-redeploy:
	serverless remove
	make aws-deploy
# docker run --rm -it amazon/aws-cli command

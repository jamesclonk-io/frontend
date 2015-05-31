TAG?=latest

all: jcio-frontend
	docker build -t jamesclonk/jcio-frontend:${TAG} .
	rm jcio-frontend

jcio-frontend: main.go
	GOARCH=amd64 GOOS=linux go build -o jcio-frontend

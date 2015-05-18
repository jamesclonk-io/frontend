TAG?=latest

all: frontend
	docker build -t jamesclonk/jcio-frontend:${TAG} .
	rm frontend

frontend: main.go
	GOARCH=amd64 GOOS=linux go build -o frontend

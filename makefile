SHELL=bash

build:
	go build main.go

test:
	go test -v ./...
cover:
	./cover ./...

gen:
	@printf "\e[92mGen mocks \e[0m\n"
	sudo docker run -v "${PWD}":/src -w /src vektra/mockery:v1.1.2 -all

deploy:
	sudo docker-compose down
	sudo docker-compose up -d --build

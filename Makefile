SHELL := /bin/zsh

define PROJECT_HELP_MSG
Usage:
  make help:\t\t show this message
  make test:\t\t run unit tests
  make lint:\t\t run go linter
  make up:\t\t run in docker compose
  make upb:\t\t rebuild and run in docker compose
  make down:\t\t shutdown docker compose
endef
export PROJECT_HELP_MSG

help:
	echo -e $$PROJECT_HELP_MSG

test:
	go test -race -tags=integration ./...

lint:
	golangci-lint run

up:
	docker-compose up

upb:
	docker-compose build; docker-compose up

down:
	docker-compose down

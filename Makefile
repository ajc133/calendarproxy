.PHONY: build build-prod deploy
build:
	env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ./bin/darwin/arm64/calendar-proxy ./cmd/server/main.go

run: build
	./bin/darwin/arm64/calendar-proxy

build-prod:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o ./devops/calendar-proxy ./cmd/server/main.go

deploy: build-prod
	ansible-playbook devops/deploy.yaml

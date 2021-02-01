goversion=1.15.6
short_sha=$(shell git rev-parse --short HEAD || echo latest)
version?=$(short_sha)
img=eu.gcr.io/matrix-varins-1556713043069/vitorarins/todoer:$(version)
vols=-v `pwd`:/app -w /app
run_go=docker run --rm $(vols) golang:$(goversion)
cov=coverage.out
covhtml=coverage.html


.PHONY: all
all: test

.PHONY: test
test:
	$(run_go) go test -coverprofile=$(cov) -race ./...

.PHONY: coverage
coverage: test
	@$(run_go) go tool cover -html=$(cov) -o=$(covhtml)
	@open $(covhtml) || xdg-open $(covhtml)

.PHONY: image
image:
	docker build -t $(img) --build-arg GOVERSION=$(goversion) --build-arg VERSION=$(version) .

.PHONY: run
run: image
	docker run -p 8080:8080 $(img)

.PHONY: publish
publish: image
	docker push $(img)

.PHONY: build
build: 
	go build -o ./cmd/todoer/todoer -ldflags "-X main.VersionString=$(version)" ./cmd/todoer/todoer.go

.PHONY: deploy
deploy: publish
	kubectl apply -k deploy

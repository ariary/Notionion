before.build:
	go mod tidy && go mod download

build.notionion:
	@echo "build in ${PWD}";go build -o notionion notionion.go
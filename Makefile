BUILD_ATTRIBUTE ?= GOOS=linux CGO_ENABLED=1

build:${SOURCE_FILES}
	@for f in $^; \
	do \
		cmdPath=$$(echo $$f | sed -e "s/.go//g");\
		buildScript="$(BUILD_ATTRIBUTE) go build -tags lambda.norpc -ldflags \" -s -w\" -o build/$$cmdPath $$cmdPath.go";\
		echo $$buildScript;\
		eval $$buildScript;\
	done;

docker-build:
	docker build --build-arg SOURCE_FILES=./cmd/main.go -t product-recommendation .

mockery:
	mockery --all --dir=pkg/repository --output=./mocks --outpkg=mocks --with-expecter
	mockery --all --dir=pkg/database --output=./mocks --outpkg=mocks --with-expecter
test:
	go test -v -cover ./pkg/usecase/...
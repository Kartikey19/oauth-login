all:
	make test
	go run *.go

lint:
	golint db
	golint websocket
	golint .

build:	build_linux_amd64	build_linux_i386

build_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o release/linux/amd64/recro

build_linux_i386:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -v -a -o release/linux/i386/recro

docker:
	docker build -t kartikey/recro  .

testr:
	GOCACHE=off make test
test:
	while IFS= read -r tbl; do `echo tbl` ; done < .env
	echo "Added .env file environments variable to system."
	go test . ./postgres -cover -v

run:	test
	go run *.go
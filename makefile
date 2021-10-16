

windows:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o responses.exe

linux: 
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o responses
all: windows linux
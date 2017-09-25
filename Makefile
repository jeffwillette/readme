build:
	@GOOS=linux go build -o readme-linux
	@GOOS=windows go build -o readme.exe
	@go build -o readme-osx

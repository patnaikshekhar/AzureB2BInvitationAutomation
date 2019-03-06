$env:GOOS = "linux"
go build
docker-compose build
docker-compose up
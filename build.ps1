$env:GOOS = "linux"
go build main.go
docker-compose build
docker-compose up
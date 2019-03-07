$env:GOOS = "linux"
echo "Building"
go build

echo "Building Image"
docker build -t patnaikshekhar/b2binvitationdemo:$env:BUILD_VERSION .

echo "Deploying Image"
docker push patnaikshekhar/b2binvitationdemo:$env:BUILD_VERSION
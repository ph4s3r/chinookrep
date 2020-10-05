

//running with docker

docker pull redis:5.0

docker run -d -p 6379:6379 --name=redis redis:5.0

docker build . -t "peet/chinookrep:latest"

docker run -it -p 8000:8000 --link redis peet/chinookrep





//running without docker (needs go 1.13 and build-essentials)

sudo apt-get install build-essential
sudo apt-get install redis-server

git clone https://github.com/ph4s3r/chinookrep.git

export GO111MODULE=on

go mod init main

redis-server&

go build . && go run .





//test

curl -v 'http://localhost:8000/doReport?toDate=20200926&fromDate=20100901'



//running with docker (tested on win&linux)

docker pull redis:5.0

docker run -d -p 6379:6379 --name=redis redis:5.0

docker build . -t "peet/chinookrep:latest"

cd chinookrep

docker run -it -p 8000:8000 --link redis peet/chinookrep





//running without docker (runs only on linux + needs go 1.13)

sudo apt-get install build-essential redis-server && git clone https://github.com/ph4s3r/chinookrep.git && cd chinookrep && export GO111MODULE=on && go mod init main && redis-server& && go build . && go run .





//test

i.e. curl -v 'http://localhost:8000/doReport?toDate=20200926&fromDate=20100901'

watch program / docker container stdout for results



//running with docker (tested on win&linux)

docker pull redis:5.0

docker run -d -p 6379:6379 --name=redis redis:5.0

docker build . -t "peet/chinookrep:latest"

cd chinookrep

docker run -it -p 8000:8000 --link redis peet/chinookrep





//test

i.e. curl -v 'http://localhost:8000/doReport?toDate=20200926&fromDate=20100901'

watch docker container stdout / logs for results

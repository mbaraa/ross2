fetch:
	go get && cd ./client && npm install

frontend:
	cd ./client && npm install && npm run build

backend:
	go build -ldflags="-w -s"

all:
	go build -ldflags="-w -s" && cd ./client && npm install && npm run build

run:
	go run main.go & cd ./client && npm start

runb:
	go run main.go

runf: 
	cd ./client && npm start	

clean:
	rm ./ross2 && cd ./client && rm -rf ./build

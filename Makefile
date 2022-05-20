fetch:
	cd ./server && go get && cd ../client && npm install

frontend:
	cd ./client && npm install && npm run build

backend:
	cd ./server && go build -ldflags="-w -s"

all:
	cd ./server && go build -ldflags="-w -s" && cd ../client && npm install && npm run build

run:
	cd ./server && go run main.go & cd ../client && npm start

runb:
	cd ./server && go run main.go

runf: 
	cd ./client && npm start	

clean:
	rm ./server/ross2 && cd ../client && rm -rf ./build

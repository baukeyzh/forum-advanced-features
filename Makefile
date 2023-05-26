run:
	go run ./cmd
build:
	docker image build -t forum:alpha .
docker-run:
	docker run -d -p 8082:8081 --name myforum forum:alpha 
docker-stop:
	docker stop $$(docker ps -aq)
docker-delete:
	docker rm $$(docker ps -aq)
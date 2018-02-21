docker-build:
	docker build -t worker-rq .

docker-run:
	docker run -it --rm --name worker-rq worker-rq


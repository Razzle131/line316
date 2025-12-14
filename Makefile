container_runtime := $(shell which podman || which docker)

$(info using ${container_runtime})

up: down
	${container_runtime} compose up --build -d

down:
	${container_runtime} compose down

run:
	docker build -t tp_model:latest -f Dockerfile.tp_model ./
	docker run --rm -it --net=host tp_model:latest nc -lkv 0.0.0.0 8080
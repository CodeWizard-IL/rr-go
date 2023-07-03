clean-forward:
	cd forward && make clean

clean-reverse:
	cd reverse && make clean

build-forward:
	cd forward && make build

build-reverse:
	cd reverse && make build

run-forward:
	cd forward && make run

run-reverse:
	cd reverse && make run

clean: clean-forward clean-reverse
build: build-forward build-reverse
run: run-forward run-reverse

docker-build-forward:
	docker build -t rr-http-proxy-forward --build-arg PACKAGE=forward .

docker-build-reverse:
	docker build -t rr-http-proxy-reverse --build-arg PACKAGE=reverse .

docker-build: docker-build-forward docker-build-reverse
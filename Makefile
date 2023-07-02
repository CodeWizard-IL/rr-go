clean-http-proxy:
	cd rr-http-proxy && make clean

build-http-proxy:
	cd rr-http-proxy && make build

run-http-proxy:
	cd rr-http-proxy && make run

clean: clean-http-proxy
build: build-http-proxy
run: run-http-proxy
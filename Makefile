build:
	docker build --tag=goscout .
run:
	docker run -p 3000:3000 goscout
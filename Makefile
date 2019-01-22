goscout:
	docker build --tag=goscout .
	docker run -p 3000:3000 goscout
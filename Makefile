all: run
.PHONY: all

run:
	docker build -t kirigaikabuto/my-city-api:latest .
	docker-compose up --build
git:
	git add .
	git commit -m "feat:add update"
	git push origin master
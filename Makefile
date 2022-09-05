all: run
.PHONY: all

run:
	docker build -t kirigaikabuto/my-city-api:latest .
	docker-compose --env-file ./config/local.env up --build
prod:
	docker build -t kirigaikabuto/my-city-api:latest .
	docker-compose --env-file ./config/prod.env up --build
git:
	git add .
	git commit -m "feat:add update"
	git push origin master
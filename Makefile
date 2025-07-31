all: build

build:
	@docker compose -f ./srcs/config/docker-compose.yml --env-file srcs/config/config.env up -d --build

down:
	@docker compose -f ./srcs/config/docker-compose.yml --env-file srcs/config/config.env down

re: down clean docs build

docs:
	cd ./srcs/requirements/app && swag init --dir ./cmd,./internal/handler --parseDependency --parseInternal --output ./docs
clean: down
	@docker system prune -a --force

fclean:
	@docker stop $$(docker ps -qa)
	@docker system prune --all --force --volumes


.PHONY	: all build down re clean fclean

# чтобы запустить тесты
# docker exec -it app_con /bin/sh
# go test ./...

# utils
# nohup ./app > server.log 2>&1 &
# con_to_db:
# 	psql -h localhost -U user_db -d wallets_db

#docker exec -it db_con psql -U user_db -d wallets_db
#\dt
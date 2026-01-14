ENV_FILE=.env
include $(ENV_FILE)

API_MIGRATE_IMAGE=api_migrate
API_TEST_IMAGE=api_test
DUMPS_PATH=/srv/dumps
image_exists = $(shell docker images -q $(1))


.PHONY: api postgres postgres_it client postgres_dump
api:
	cd api && ENV_FILE=../${ENV_FILE} docker compose --env-file=../${ENV_FILE} up --build -d

postgres:
	cd postgres && ENV_FILE=../${ENV_FILE} docker compose --env-file=../${ENV_FILE} up --build -d

postgres_it:
	docker exec -it postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

postgres_dump:
	docker exec -it postgres pg_dump -U $(POSTGRES_USER) -d $(POSTGRES_DB) > \
		${DUMPS_PATH}/pg_dump_$(shell date +"%Y-%m-%dT%H-%M-%S").sql

client:
	cd ${HOME_PATH}/client && ENV_FILE=../${ENV_FILE} docker compose --env-file=../${ENV_FILE} up --build -d

.PHONY: migrate_build migrate
migrate_build:
ifeq ($(call image_exists,$(API_MIGRATE_IMAGE)),)
	cd api && docker build -f Dockerfile.migrate -t $(API_MIGRATE_IMAGE) .
endif

migrate: postgres_dump migrate_build
	docker run --rm \
		--network db-network \
		--env-file $(ENV_FILE) \
		-v ${LOGS_PATH}:/logs/ \
		-v ${HOME_PATH}/api/pkg/migrations:/app/pkg/migrations:ro \
		$(API_MIGRATE_IMAGE) || true
	docker rmi $(API_MIGRATE_IMAGE)
	

.PHONY: test_build test
test_build:
ifeq ($(call image_exists,$(API_TEST_IMAGE)),)
	cd api && docker build -f Dockerfile.test -t $(API_TEST_IMAGE) .
endif

test: test_build
	docker run --rm \
		--network db-network \
		--env-file $(ENV_FILE) \
		-v ${LOGS_PATH}:/logs/ \
		$(API_TEST_IMAGE) || true
	docker rmi $(API_TEST_IMAGE)

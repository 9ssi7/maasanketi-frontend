.PHONY: backend
backend:
	@$(MAKE) -C modules/backend up

.PHONY: web
web:
	@$(MAKE) -C modules/web up

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

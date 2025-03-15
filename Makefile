.PHONY: backend
backend:
	@$(MAKE) -C modules/backend up

.PHONY: web
web:
	@$(MAKE) -C modules/web up
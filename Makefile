build:
	$(MAKE) -C src/cos build
	$(MAKE) -C src/centrifugo build

run:
	@rm -rf config.json
	@docker run -v $(PWD)/:/centrifugo -t centrifugo/centrifugo centrifugo genconfig
	docker-compose up



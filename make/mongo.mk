OS_NAME := $(shell uname -s | tr A-Z a-z)
TIMESTAMP ?= $$(date -u +'%Y%m%d%H%M%S')
MONGO_VERSION := 4.4.0

install-mongo:
ifeq ($(findstring darwin,$(OS_NAME)),darwin)
	@echo Installing Mongo $(MONGO_VERSION)
	@rm -rf mongodb-macos-x86_64-$(MONGO_VERSION) mongodb-macos-x86_64-$(MONGO_VERSION).tgz
	@wget https://fastdl.mongodb.org/osx/mongodb-macos-x86_64-$(MONGO_VERSION).tgz
	@tar -zxvf mongodb-macos-x86_64-$(MONGO_VERSION).tgz
	@sudo cp ./mongodb-macos-x86_64-$(MONGO_VERSION)/bin/* /usr/local/bin/
	@rm -rf mongodb-macos-x86_64-$(MONGO_VERSION) mongodb-macos-x86_64-$(MONGO_VERSION).tgz
endif

run-mongo:
	@mkdir -p ./mongodb/data && mkdir -p ./mongodb/log
	@echo Running Mongo as Forked service
	@mongod --config ./mongo.conf --dbpath ./mongodb/data --logpath ./mongodb/log/mongo.log

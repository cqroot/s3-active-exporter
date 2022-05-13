EXPORTER_NAME=s3-active-exporter

.PHONY: build
build:
	@BuildVersion=$$(git describe --tags --abbrev=0); \
		echo "Build Version: $${BuildVersion}"; \
		sed -i "s/BuildVersion string = \"[^\"]*\"/BuildVersion string = \"$${BuildVersion}\"/" internal/version.go
	@mkdir -p $(EXPORTER_NAME)/bin
	@CGO_ENABLED=0 go build -o $(EXPORTER_NAME)/bin/$(EXPORTER_NAME) main.go

.PHONY: pack
pack: build
	cp -r conf systemd $(EXPORTER_NAME)/

.PHONY: run
run: pack
	@$(EXPORTER_NAME)/bin/$(EXPORTER_NAME) --config $(EXPORTER_NAME)/conf/$(EXPORTER_NAME).yml

.PHONY: clean
clean:
	go clean
	rm -rf ./bin/$(EXPORTER_NAME) ./$(EXPORTER_NAME)

.PHONY: tar
tar: pack
	BuildVersion=$$(git describe --tags --abbrev=0); \
		tar cvf $(EXPORTER_NAME)-$${BuildVersion}-$$(arch).tar $(EXPORTER_NAME)

.PHONY: docker-build
docker-build: pack
	@BuildVersion=$$(git describe --tags --abbrev=0); \
		docker build --force-rm -t $(EXPORTER_NAME):$${BuildVersion:1} .

.PHONY: docker-run
docker-run:
	docker run \
		-itd -P \
		--hostname $(EXPORTER_NAME) \
		--name $(EXPORTER_NAME) \
		$(EXPORTER_NAME) --log.debug

.PHONY: docker-exec
docker-exec:
	docker exec -it $(EXPORTER_NAME) /bin/sh

.PHONY: docker-clean
docker-clean: clean
	docker rm -f $(EXPORTER_NAME); docker rmi $(EXPORTER_NAME)

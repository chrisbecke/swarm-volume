PLUGIN_IMAGE = docker-volume-plugin:nfs
PLUGIN = chrisbecke/swarm-volume:nfs
TARGET = nfs-swarm-volume

# PLUGIN_IMAGE = docker-volume-plugin:efs
# PLUGIN = chrisbecke/swarm-volume:efs
# TARGET = efs-swarm-volume

# PLUGIN_IMAGE = docker-volume-plugin:gfs
# PLUGIN = chrisbecke/swarm-volume:glusterfs
# TARGET = glusterfs-swarm-volume

TMP_ID = $(shell docker create $(PLUGIN_IMAGE))

EFS_ID = fs-03a6f6df47f67f7d1
GFS_SERVERS=lab717.mgsops.net,lab718.mgsops.net,lab719.mgsops.net
GFS_VOLUME=gv0

build: build-plugin

build-image:
	@echo [MAKE] Building image $(PLUGIN_IMAGE)
#	docker buildx build --progress plain . --tag $(PLUGIN_IMAGE)
	@docker compose build $(TARGET)

build-plugin: plugin/$(TARGET)/rootfs plugin/$(TARGET)/config.json
	@echo [MAKE] Creating plugin $(PLUGIN) from plugin/$(TARGET)
	@docker --context default plugin disable --force $(PLUGIN) ; true
	@docker --context default plugin rm --force $(PLUGIN) ; true
	@sudo docker plugin create $(PLUGIN) ./plugin/$(TARGET)/

plugin:

plugin/$(TARGET)/rootfs:
	@$(MAKE) build-image
	@echo [MAKE] Exporting plugin to $@
	@mkdir -p $@
	@-docker export "$(TMP_ID)" | sudo tar -x -C $@

plugin/$(TARGET)/config.json: $(TARGET)/config.json
	@echo [MAKE] Generating $@
	@cp $< $@

clean:
	@echo [MAKE] Deleting plugin/$(TARGET)
	@sudo rm -rf plugin/$(TARGET)

test:
	@-docker plugin disable $(PLUGIN) --force
	@docker plugin set $(PLUGIN) GFS_VOLUME=$(GFS_VOLUME) GFS_SERVERS=$(GFS_SERVERS)
#	@docker plugin set $(PLUGIN) EFS_ID=$(EFS_ID)
	@docker plugin enable $(PLUGIN)

push:
	@docker plugin push $(PLUGIN)
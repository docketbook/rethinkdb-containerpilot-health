DOCKER_STD_ARGS=-it --rm -v $(shell pwd):/go/src/${IMPORT_PATH} -w /go/src/${IMPORT_PATH}
IMPORT_PATH := github.com/docketbook/rethinkdb-health
ROOT := $(shell pwd)
#github.com/hashicorp/consul/api

DOCKERRUN := docker run -it --rm \
	-v ${ROOT}/vendor:/go/src \
	-v ${ROOT}:/go/src/${IMPORT_PATH} \
	-w /go/src/${IMPORT_PATH} \
	rethinkdb_health_work_image

cmd: build/rethinkdb_health_work_image
	$(DOCKERRUN)

release: build/rethinkdb_health_work_image
	$(DOCKERRUN) sh -c "godep restore && go install && mkdir -p ./bin && cp /go/bin/rethinkdb-health ./bin/rethinkdb-health"
	
# builds the builder container
build/rethinkdb_health_work_image:
	mkdir -p ${ROOT}/build
	docker rmi -f rethinkdb_health_work_image > /dev/null 2>&1 || true
	docker build -t rethinkdb_health_work_image ${ROOT}
	docker inspect -f "{{ .ID }}" rethinkdb_health_work_image > build/rethinkdb_health_work_image

SYSTEM=$(shell go env GOOS)
XLSX_PATH=./configure/xlsx
PROTO_PATH=./configure/proto
DATA_PATH=./configure/data
CFG_GO_PATH=./common/config/repository/
REDIS_GO_PATH=./common/dao/repository/redis
PB_GO_PATH=./common/pb
OUTPUT=./output


.PHONY: config pb pbtool docker_stop docker_run

############################生成代码选项##############################
config:
	@echo "gen config code..."
	@rm -rf ${CFG_GO_PATH}
	@go run ./tools/cfgtool/main.go -xlsx=${XLSX_PATH} -proto=${PROTO_PATH} -code=${CFG_GO_PATH} -text=${DATA_PATH} -pb=${PB_GO_PATH}
	make pb
	make pbtool

pb:
	@echo "Building pb"
	-rm -rf ${PB_GO_PATH}/*.pb.go && mkdir -p ${PB_GO_PATH}
ifeq (${SYSTEM}, windows)
	protoc.exe -I${PROTO_PATH} ${PROTO_PATH}/*.proto --go_out=..
else # linux darwin(mac)
	protoc -I${PROTO_PATH} ${PROTO_PATH}/*.proto --go_out=..
endif 

pbtool:
	@echo "gen redis code..."
	@go run ./tools/pbtool/main.go -pb=${PB_GO_PATH} -redis=${REDIS_GO_PATH}



docker_stop:
	@echo "停止docker环境"
	docker-compose -f ./configure/env/local/docker_compose.yaml down

docker_run:
	@echo "启动docker环境"
	docker-compose -f ./configure/env/local/docker_compose.yaml up -d

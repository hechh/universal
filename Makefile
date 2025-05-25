
SYSTEM=$(shell go env GOOS)
XLSX_PATH=./configure/xlsx
PROTO_PATH=./configure/proto
DATA_PATH=./configure/data
CFG_GO_PATH=./common/config/repository/
REDIS_GO_PATH=./common/dao/repository/redis
PB_GO_PATH=./common/pb
HTTP_KIT_GO_PATH=./server/client/httpkit
OUTPUT=./output
SERVER_PATH=./server


## 需要编译的服务
TARGET=client gate
LINUX=$(TARGET:%=%_linux)
BUILD=$(TARGET:%=%_build)
START=$(TARGET:%=%_start)
STOP=$(TARGET:%=%_stop)


.PHONY: ${TARGET} config pb pbtool docker_stop docker_run


all: clean
	make ${BUILD}

linux: clean
	make ${LINUX}

build: clean ${BUILD}

clean:	
	-rm -rf ${OUTPUT}
	@mkdir -p ${OUTPUT}/log

# 交叉编译(编译linux执行文件)
$(LINUX): %_linux: %
	@echo "Building $*"
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build ${GCFLAGS} -o ${OUTPUT}/ ${SERVER_PATH}/$*

# 随系统编译(编译当前系统执行文件)
$(BUILD): %_build: % 
	@echo "Building $*"
	go build ${GCFLAGS} -o ${OUTPUT}/$* ${SERVER_PATH}/$*


#------------------------生成代码选项-----------------------------
config:
	@echo "gen config code..."
	@rm -rf ${CFG_GO_PATH}
	@go run ./tools/cfgtool/main.go -xlsx=${XLSX_PATH} -proto=${PROTO_PATH} -code=${CFG_GO_PATH} -text=${DATA_PATH} -pb=${PB_GO_PATH} -client=${HTTP_KIT_GO_PATH}
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


#------------------------docker环境选项-----------------------------
docker_stop:
	@echo "停止docker环境"
	docker-compose -f ./configure/env/local/docker_compose.yaml down

docker_run:
	@echo "启动docker环境"
	docker-compose -f ./configure/env/local/docker_compose.yaml up -d

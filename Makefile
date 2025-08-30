
SYSTEM=$(shell go env GOOS)

XLSX_PATH=./configure/xlsx
GEN_DATA_PATH=./configure/data
GEN_PROTO_PATH=./configure/proto
GEN_PB_GO_PATH=./common/pb
GEN_CFG_PATH=./common/config/repository/
OUTPUT=./output
GEN_REDIS_GO_PATH=./common/redis/repository/


.PHONY: ${TARGET} cfgtool dbtool pb startall stopall docker_stop docker_run

TARGET=gate
LINUX=$(TARGET:%=%_linux)
BUILD=$(TARGET:%=%_build)


all: clean
	@cp -rf ./configure/env/local/* ./configure/data ${OUTPUT}
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
cfgtool:
	@echo "gen config code..."
	@rm -rf ${GEN_CFG_PATH}
#	@go run ./tool/cfgtool/main.go -xlsx=${XLSX_PATH} -proto=${GEN_PROTO_PATH}  -text=${GEN_DATA_PATH}
	@go run ./tool/cfgtool/main.go -xlsx=${XLSX_PATH} -proto=${GEN_PROTO_PATH}  -text=${GEN_DATA_PATH} -pb=${GEN_PB_GO_PATH} -code=${GEN_CFG_PATH}
	make pb

dbtool: pb
	@echo "gen redis code..."
	@rm -rf ${GEN_REDIS_GO_PATH}
	@go run ./tool/dbtool/main.go -pb=${GEN_PB_GO_PATH} -redis=${GEN_REDIS_GO_PATH}

pb:
	@echo "Building pb"
	-rm -rf ${GEN_PB_GO_PATH}/*.pb.go && mkdir -p ${GEN_PB_GO_PATH}
ifeq (${SYSTEM}, windows)
	protoc.exe -I${GEN_PROTO_PATH} ${GEN_PROTO_PATH}/*.proto --go_out=..
else # linux darwin(mac)
	protoc -I${GEN_PB_GO_PATH} ${GEN_PROTO_PATH}/*.proto --go_out=..
endif 
	@go run ./tool/pbtool/main.go -pb=${GEN_PB_GO_PATH} 

#------------------------docker环境选项-----------------------------
startall: ${START}

stopall: $(STOP) 

$(START): %_start: %
	@echo "running $*"
	-cd ./output && nohup ./$* -config=./local.yaml -id=1 >./log/$*_monitor.log 2>&1 &

$(STOP): %_stop: %
	@echo "stopping $*"
	-kill -9 $$(ps -ef | grep $* | grep -v grep | awk '{print $$2}')

docker_stop:
	@echo "停止docker环境"
	docker-compose -f ./configure/env/docker_compose.yaml down

docker_run:
	@echo "启动docker环境"
	docker-compose -f ./configure/env/docker_compose.yaml up -d


SYSTEM=$(shell go env GOOS)
GCFLAGS=-gcflags "all=-N -l"
PROTO_PATH=./configure/proto
TABLE_PATH=./configure/table
JSON_PATH=./configure/json
PROTO_PATH=./configure/proto
CFG_PATH=./common/config/cfg
CODE_PATH=./common/config/repository
PB_PATH=./common/pb
OUTPUT=./output


.PHONY: protoc tool proto

############################生成代码选项##############################
protoc:
	-mkdir -p ${PB_PATH} && rm -rf ${PB_PATH}/*
ifeq (${SYSTEM}, windows)
	protoc.exe -I${PROTO_PATH} ${PROTO_PATH}/*.proto --go_opt paths=source_relative --go_out=${PB_PATH}
else # linux darwin(mac)
	protoc -I${PROTO_PATH} ${PROTO_PATH}/*.proto --go_opt paths=source_relative --go_out=${PB_PATH}
endif 


##########################client工具代码自动生成#######################
proto: 
	go run ./tools/gomaker/main.go -action=proto -xlsx=${TABLE_PATH} -dst=${PROTO_PATH} -src=${PROTO_PATH}

# 配置转换工具
xlsx:
	go run ./tools/cfgtool/main.go -xlsx=${TABLE_PATH} -data=${JSON_PATH} -pb=${CFG_PATH} -code=${CODE_PATH} -proto=${PROTO_PATH}


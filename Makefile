
SYSTEM=$(shell go env GOOS)
GCFLAGS=-gcflags "all=-N -l"
PROTO_PATH=./protocol
GEN_GO_PATH=./common/pb
OUTPUT=./output


.PHONY: protoc

############################生成代码选项##############################
protoc:
	-mkdir -p ${GEN_GO_PATH} && rm -rf ${GEN_GO_PATH}/*
ifeq (${SYSTEM}, windows)
	protoc.exe -I${PROTO_PATH} ${PROTO_PATH}/*.proto --go_opt paths=source_relative --go_out=${GEN_GO_PATH}
else # linux darwin(mac)
	protoc -I${PROTO_PATH} ${PROTO_PATH}/*.proto --go_opt paths=source_relative --go_out=${GEN_GO_PATH}
endif 



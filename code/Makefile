
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


##########################client工具代码自动生成#######################
tool: protoc
#	go install ./tool/gomaker
	go run ./tools/gomaker/main.go -action="httpkit.tpl" -src="./common/pb" -dst="./tools/client/internal/httpkit/Init.gen.go" -tpl="./tools/gomaker/templates/"
	go run ./tools/gomaker/main.go -action=pbclass.tpl -src="common/pb" -dst="tools/client/internal/httpkit/pbclass.gen.go" -tpl="tools/gomaker/templates/"
	go run ./tools/gomaker/main.go -action="proto.tpl" -src="common/pb" -dst="tools/client/internal/httpkit/json.gen.go" -tpl="tools/gomaker/templates/"
#	go install ./tools/client

cfgtool:
	go run ./tools/cfgtool/main.go -src="./share/table/" -dst="./output/json/"

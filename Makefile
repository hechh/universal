
SYSTEM=$(shell go env GOOS)
GCFLAGS=-gcflags "all=-N -l"
PROTO_PATH=./configure/proto
PB_PATH=./common/pb
OUTPUT=./output


.PHONY: protoc tool pb

############################生成代码选项##############################
protoc:
	-mkdir -p ${PB_PATH} && rm -rf ${PB_PATH}/*
ifeq (${SYSTEM}, windows)
	protoc.exe -I${PROTO_PATH} ${PROTO_PATH}/*.proto --go_opt paths=source_relative --go_out=${PB_PATH}
else # linux darwin(mac)
	protoc -I${PROTO_PATH} ${PROTO_PATH}/*.proto --go_opt paths=source_relative --go_out=${PB_PATH}
endif 


##########################client工具代码自动生成#######################
tool: 
#	go install ./tool/gomaker
#	go run ./tools/gomaker/main.go -action=client -src="./tools/gomaker/test" -dst="./tools/client/internal/httpkit" -tpl="./tools/gomaker/templates/"
	go run ./tools/gomaker/main.go -action=pb -src="./configure/table" -dst="./configure/proto" -tpl="./tools/gomaker/templates/"
#	go install ./tools/client

pb: 
	-rm -rf ${PROTO_PATH}/*.gen.proto ${PB_PATH}/*.pb.go
	go run ./tools/gomaker/main.go -action=pb -src="./configure/table" -dst="./configure/proto" -tpl="./tools/gomaker/templates/"
	make protoc

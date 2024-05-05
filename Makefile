SYSTEM=$(shell go env GOOS)
GCFLAGS=-gcflags "all=-N -l"
RACE=
OUTPUT=./output
PROTO_PATH=./proto
GEN_GO_PATH=./common/pb

TARGET=gate 
BUILD=$(TARGET:%=%_build)

.PHONY: race build clean all gen gen_go yaml copy

############################编译选项##############################
#--------设置target变量---------
race:RACE=-race
#---------程序编译选项-----------
all: clean copy $(TARGET)

clean:
	-rm -rf ${OUTPUT}

$(TARGET): gen
	-mkdir -p ${OUTPUT}/bin
	-mkdir -p ${OUTPUT}/yaml && cp -rf ./env/${@}.yaml ${OUTPUT}/yaml/
ifeq (${SYSTEM}, windows)
	go build ${GCFLAGS} ${RACE} -o ${OUTPUT}/bin/$@.exe ./cmd/$@/...
else
ifeq (${SYSTEM}, linux)
	CGO_ENABLED=0 GOOS=linux go build ${GCFLAGS} ${RACE} -o ${OUTPUT}/bin/ ./cmd/$@/...
else
	CGO_ENABLED=0 GOOS=darwin go build ${GCFLAGS} ${RACE} -o ${OUTPUT}/bin/ ./cmd/$@/...
endif
endif

############################生成代码选项##############################
copy:
	-mkdir -p ${OUTPUT}/ && cp -rf ./env/*.sh ${OUTPUT}/

gen:
	-mkdir -p ${GEN_GO_PATH} && rm -rf ${GEN_GO_PATH}/*
ifeq (${SYSTEM}, windows)
	protoc.exe -I${PROTO_PATH} ${PROTO_PATH}/*.proto --go_opt paths=source_relative --go_out=${GEN_GO_PATH}
else # linux darwin(mac)
	protoc -I${PROTO_PATH} ${PROTO_PATH}/*.proto --go_opt paths=source_relative --go_out=${GEN_GO_PATH}
endif 
	gomaker -action=uerrors -src="common/pb/*pb.go" -dst="common/uerrors/" -tpl="tools/gomaker/templates"


stop:
	./output/run.sh stop gate
start:
	./output/run.sh start gate

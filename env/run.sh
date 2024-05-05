#!/bin/bash

stop(){
   pids=$(ps -ef | grep -v grep | grep $1 | awk '{print $2}') 
    for pid in $pids
    do
        kill -9 ${pid}
        echo "kill -f ${pid}"
    done
}

start(){
   ./output/bin/$1 -id 1 -yaml output/yaml/${1}.yaml
}


#*****************处理****************
case $1 in
"stop")
    stop $2
    ;;
"start")
    start $2
    ;;
*)  # 默认情况，如果没有匹配的情况
    echo "Unknown"
    ;;
esac
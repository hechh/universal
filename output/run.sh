#!/bin/bash

stop(){
   pids=$(ps -ef | grep -v grep | grep $1 | awk '{print $2}') 
    for pid in $pids
    do
        echo "kill -9 ${pid}"
        kill -9 ${pid}
    done
}

start(){
#    nohup ./output/bin/$1 -id 1 -yaml output/bin/yaml/${1}.yaml -log output/log/${1}.log &
    ./output/bin/$1 -id 1 -yaml output/bin/yaml/${1}.yaml
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
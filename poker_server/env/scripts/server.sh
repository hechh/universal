#!/bin/bash

hchall(){
    ./start.sh room 1
    ./start.sh room 2
    ./start.sh db 1
    ./start.sh db 2
    ./start.sh builder 1
    ./start.sh builder 2
    ./start.sh match 1
    ./start.sh match 2
    ./start.sh game 1
    ./start.sh game 2
    ./start.sh gate 1
    ./start.sh gm 1
#    ./start.sh client 1
}

startall(){
    ./start.sh room 1
    ./start.sh db 1
    ./start.sh builder 1
    ./start.sh match 1
    ./start.sh game 1
    ./start.sh gate 1
    ./start.sh gm 1
    ./start.sh client 1
}

stopall(){
    ./stop.sh client
    ./stop.sh gate
    ./stop.sh room
    ./stop.sh match
    ./stop.sh builder 
    ./stop.sh game
    ./stop.sh gm
    ./stop.sh db
    return 0
}

case $1 in
# 测试使用
hch) 
    hchall
    ;;
startall)
    startall
    ;;
stopall)
    stopall
    ;;
restart)
    stopall
    startall
    ;;
esac

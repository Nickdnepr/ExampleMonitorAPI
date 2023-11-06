#!/bin/bash

start() {
	cd /opt/monitor_api/
	if [[ ! -f "lab2.db" ]]; then
		./lab2Run --createdb
	fi
	./lab2Run --start
}

stop() {
	kill -s SIGKILL `ps -ef | grep -i "lab2Run" | awk '{print $2;}'`
}

case $1 in
	start|stop) "$1" ;;
esac
echo ""

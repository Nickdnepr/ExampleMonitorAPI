#!/bin/bash

if [[ "$EUID" -ne 0 ]]
then
	echo "Root acces is required!"
	exit 0
fi

if [[ ! -f /usr/local/go/bin/go ]]
then
	echo "Golang is required to install this service!"
	exit 0
fi

/usr/local/go/bin/go build -o lab2Run . 
mkdir -p /opt/monitor_api/pub
cp monitor_api.sh /opt/monitor_api/monitor_api.sh
cp lab2Run /opt/monitor_api/lab2Run
cp monitor_api.service /etc/systemd/system/monitor_api.service
cp monitors /opt/monitor_api/monitors

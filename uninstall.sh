#!/bin/bash

if [[ "$EUID" -ne 0 ]]
then
	echo "Root acces is required!"
	exit 0
fi

rm -rf /opt/monitor_api/
rm /etc/systemd/system/monitor_api.service

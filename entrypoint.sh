#!/bin/sh
set -e

if [ ! -f "/etc/taskmate/conf/config.yml" ];then
    cp /etc/taskmate/config.yml /etc/taskmate/conf/config.yml
fi

exec "$@"

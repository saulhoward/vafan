#!/bin/sh

### BEGIN INIT INFO
# Provides:          vafan
# Required-Start:    $all
# Required-Stop:     $all
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: starts the vafan server
# Description:       starts vafan using start-stop-daemon
### END INIT INFO

PATH=/sbin:/bin:/usr/sbin:/usr/bin

BIN=/usr/sbin/vafan-server
PIDFILE=/var/run/vafan.pid
USER=vafan
GROUP=vafan

test -f $BIN || exit 0
set -e
case "$1" in
  start)
    echo -n "Starting vafan server: "
    start-stop-daemon --start --chuid $USER:$GROUP \
        --make-pidfile --background --pidfile $PIDFILE \
        --exec $BIN
    echo "vafan."
    ;;
  stop)
    echo -n "Stopping vafan server: "
    start-stop-daemon --stop --quiet --pidfile $PIDFILE --exec $BIN
    rm -f $PIDFILE
    echo "vafan."
    ;;
  restart)
    echo -n "Restarting vafan server: "
    $0 stop
    sleep 1
    $0 start
    echo "vafan."
    ;;
  *)
    echo "Usage: $0 {start|stop|restart}" >&2
    exit 1
    ;;
esac
exit 0

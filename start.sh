#!/bin/sh
process="up2c"


PID=$(ps -o pid,comm | grep "$process" | grep -v grep | awk '{print $1}')

case "$1" in
start)
  if [ `echo ${PID} | awk -v tem=0 '{print($1>tem)? "1":"0"}'` -eq 0 ]; then
    nohup /etc/up2c/up2c > /dev/null 2>&1 &
    echo "The "$process" is start..."
  else
    echo "The "$process" is running..."
  fi
  ;;
restart)
  if [ `echo ${PID} | awk -v tem=0 '{print($1>tem)? "1":"0"}'` -eq 0 ]; then
    echo "The "$process" not running..."
  else
    kill -9 $PID
    nohup /etc/up2c/up2c > /dev/null 2>&1 &
    echo "The "$process" is restart..."
  fi
  ;;
stop)
  if [ `echo ${PID} | awk -v tem=0 '{print($1>tem)? "1":"0"}'` -eq 0 ]; then
    echo "The not running..."
  else
    kill -9 "$PID"
    echo "The "$process" is stop..."
  fi
  ;;
status)
  if [ `echo ${PID} | awk -v tem=0 '{print($1>tem)? "1":"0"}'` -eq 0 ]; then
    echo "The "$process" not running..."
  else
    echo "The "$process" is running"
  fi
  ;;
*)
  echo "Usage:{start|restart|stop|status}"
  ;;

esac
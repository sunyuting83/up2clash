# 进程名
process="up2c"

# 获取进程ID
PID=$(ps -C "$process" --no-header |wc -l)

case "$1" in
start)
  if [ "$PID" -eq 0 ]; then
    nohup /etc/up2c/up2c > /dev/null 2>&1 &
    echo "The "$process" is start..."
  else
    echo "The "$process" is running..."
  fi
  ;;
restart)
  if [ "$PID" -eq 0 ]; then
    echo "The "$process" not running..."
  else
    kill -9 $PID
    nohup /etc/up2c/up2c > /dev/null 2>&1 &
    echo "The "$process" is restart..."
  fi
  ;;
stop)
  if [ "$PID" -eq 0 ]; then
    echo "The not running..."
  else
    kill -9 "$PID"
    echo "The "$process" is stop..."
  fi
  ;;
status)
  if [ "$PID" -eq 0 ]; then
    echo "The "$process" not running..."
  else
    echo "The "$process" is running"
  fi
  ;;
*)
  echo "Usage:{start|restart|stop|status}"
  ;;

esac
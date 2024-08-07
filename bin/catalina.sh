#!/bin/bash
WORK_DIR=$(cd "$(dirname "$0")/.." && pwd)
cd "${WORK_DIR}"
# 应用名
APP_NAME="${WORK_DIR##*/}"
# 添加全路径是为了进程检测
PROCESS_CHECK="${WORK_DIR}/bin/${APP_NAME}"
# 配置目录
CONF_DIR=./conf
# 日志目录
LOG_DIR=./logs
# 优雅关机最长等待时间(秒)
SHUTDOWN_SECONDS=30

# 扩展脚本
if [ -f ${CONF_DIR}/setenv.sh ]; then
  source ${CONF_DIR}/setenv.sh
fi

start() {
    echo "Starting ${APP_NAME} ..."
    echo "WORK_DIR=${WORK_DIR},APP_NAME=${APP_NAME}"
    if [ ! -d logs ]; then
        mkdir logs
    fi
    RUN_CMD="${WORK_DIR}/bin/$APP_NAME"
    if [ "${IN_CONTAINER}" == true ]; then
        $RUN_CMD > ${LOG_DIR}/catalina.out 2>&1
    else
        nohup $RUN_CMD > ${LOG_DIR}/catalina.out 2>&1 &
    fi
    echo "Waiting and checking process 10s..."
    for i in {1..10}; do
      echo -n "-"
      sleep 1
    done
    echo
    pid=$(getPid)
    if [ -z "${pid}" ]; then
      echoRed "${APP_NAME} start failed"
      exit 1
    fi
    if [ "${IN_CONTAINER}" != true ]; then
        startSupervisor
    fi
    echoGreen "Started ${APP_NAME},pid is ${pid}"
}

stop() {
  if [ "${IN_CONTAINER}" != true ]; then
     uninstallSupervisor
  fi
  pid=$(getPid)
  if [ -z "${pid}" ]; then
    echoRed "${APP_NAME} not running"
    return
  fi
  echoRed "Stopping ${APP_NAME}, pid: ${pid}"
  kill "$pid"
  sleep 1
  if [ -n "$(getPid)" ]; then
    echoRed "Graceful stop Waiting ${SHUTDOWN_SECONDS} s..."
  fi
  while [ $SHUTDOWN_SECONDS -gt 0 ]; do
    if [ -n "$(getPid)" ]; then
      echo '${SHUTDOWN_SECONDS}'
      sleep 1
      SHUTDOWN_SECONDS=$(($SHUTDOWN_SECONDS - 1))
    else
      break
    fi
  done
  pid=$(getPid)
  if [ -n "${pid}" ]; then
    kill -9 "$pid"
    echoRed "Stopped ${APP_NAME} force"
  else
    echoGreen "Stopped ${APP_NAME}"
  fi
}


restart() {
  stop
  sleep 1
  start
}

startSupervisor() {
  cmd="${WORK_DIR}/bin/catalina.sh supervise"
  exists=$(crontab -l | grep "${cmd}" | grep -v grep | wc -l)
  if [ ${exists} -lt 1 ]; then
    crontab <<EOF
$(crontab -l)
* * * * * ${cmd} > /dev/null 2>&1 &
EOF
  fi
}
# 关闭自动拉起
uninstallSupervisor() {
  cmd="${WORK_DIR}/bin/catalina.sh supervise"
    crontab <<EOF
$(crontab -l | grep -v "${cmd}" | grep -v grep)
EOF
echoRed "Uninstalled supervisor"
}

supervise() {
  log=${LOG_DIR}/supervisor.log
  pid=$(getPid)
  if [ -z "${pid}" ]; then
    time=$(date +"%Y-%m-%d %H:%M:%S")
    echo "=============== ${time} ===============" >>${log}
    echo "${APP_NAME} not running， now start a new process." >>${log}
    start
  fi
}

health() {
  pid=$(getPid)
  if [ -z "${pid}" ]; then
    echoRed "${APP_NAME} not running"
    exit 1
  else
    echoGreen "${APP_NAME} is running，pid is ${pid}"
  fi
}

getPid() {
  pid=$(ps -ef | grep "${PROCESS_CHECK}" | grep -v grep | awk '{print $2}' | head -n 1)
  echo "${pid}"
}

echoRed() {
  echo -e "\033[31m$*\033[0m"
}
echoGreen() {
  echo -e "\033[32m$*\033[0m"
}
echoYellow() {
  echo -e "\033[33m$*\033[0m"
}

usage() {
  echoRed "usage:"
  echoYellow "catalina.sh [start|stop|restart|supervise|health|help]"
}

case $1 in
start)
  start
  ;;
stop)
  stop
  ;;
restart)
  restart
  ;;
health)
  health
  ;;
supervise)
  supervise
  ;;
*)
  usage
  exit 1
  ;;
esac
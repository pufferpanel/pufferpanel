#!/bin/bash
# pufferd Installation Script

pufferdVersion={{ pufferdVersion }}

export DEBIAN_FRONTEND=noninteractive
downloadUrl="https://dl.pufferpanel.com/pufferd/${pufferdVersion}/pufferd"

RED=$(tput setf 4)
GREEN=$(tput setf 2)
NORMAL=$(tput sgr0)
BOLD=$(tput bold)

function checkResponseCode() {
    if [ $? -ne 0 ]; then
        echo -e "${RED}An error occurred while installing, halting...${NORMAL}"
        exit 1
    fi
}

if [ "$(id -u)" != "0" ]; then
    echo "This script must be run as root!" 1>&2
    exit 1
fi

if [ "$SUDO_USER" == "" ]; then
    SUDO_USER="root"
fi

if type apt-get &> /dev/null; then
    if [[ -f /etc/debian_version ]]; then
        echo -e "System detected as some variant of Ubuntu or Debian."
        OS_INSTALL_CMD="apt"
    else
        echo -e "${RED}This OS does not appear to be supported by this program!${NORMAL}"
        exit 1
    fi
elif type yum &> /dev/null; then
    echo -e "System detected as CentOS variant."
    OS_INSTALL_CMD="yum"
else
    echo -e "${RED}This OS does not appear to be supported by this program, or apt-get/yum is not installed!${NORMAL}"
    exit 1
fi

# Install Other Dependencies
echo "Installing some dependiencies."
if [ $OS_INSTALL_CMD == 'apt' ]; then
    apt-get install -y openssl curl openjdk-8-jdk tar python lib32gcc1 lib32tinfo5 lib32z1 lib32stdc++6
else
    yum -y install openssl curl java-1.8.0-openjdk-devel tar python glibc.i686 libstdc++.i686
fi

# Ensure /srv exists
mkdir -p /srv/pufferd

cd /srv/pufferd
curl -L -o pufferd $downloadUrl
checkResponseCode

mkdir /var/lib/pufferd /etc/pufferd

chmod +x pufferd
./pufferd -install -auth {{ settings.master_url }} -token {{ node.daemon_secret }} -config /etc/pufferd/config.json
checkResponseCode

chown -R pufferd:pufferd /srv/pufferd /var/lib/pufferd /etc/pufferd
checkResponseCode

initScript="
#!/bin/sh
### BEGIN INIT INFO
# Provides:          pufferd
# Required-Start:    $local_fs $network $named $time $syslog
# Required-Stop:     $local_fs $network $named $time $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Description:       pufferd daemon service
### END INIT INFO

SCRIPT=/srv/pufferd/pufferd\ --config=/etc/pufferd/config.json
RUNAS=pufferd

PIDFILE=/var/run/pufferd.pid
LOGFILE=/var/log/pufferd.log

start() {
  if [ -f /var/run/$PIDNAME ] && kill -0 $(cat /var/run/$PIDNAME); then
    echo 'Service already running' >&2
    return 1
  fi
  echo 'Starting service…' >&2
  local CMD="$SCRIPT &> \"$LOGFILE\" & echo \$!"
  su -c "$CMD" $RUNAS > "$PIDFILE"
  echo 'Service started' >&2
}

stop() {
  if [ ! -f "$PIDFILE" ] || ! kill -0 $(cat "$PIDFILE"); then
    echo 'Service not running' >&2
    return 1
  fi
  echo 'Stopping service…' >&2
  kill -15 $(cat "$PIDFILE") && rm -f "$PIDFILE"
  echo 'Service stopped' >&2
}

uninstall() {
  echo -n "Are you really sure you want to uninstall this service? That cannot be undone. [yes|No] "
  local SURE
  read SURE
  if [ "$SURE" = "yes" ]; then
    stop
    rm -f "$PIDFILE"
    echo "Notice: log file is not be removed: '$LOGFILE'" >&2
    update-rc.d -f pufferd remove
    rm -fv "$0"
  fi
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  uninstall)
    uninstall
    ;;
  restart)
    stop
    start
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|uninstall}"
esac
"

if [ command -v systemctl >/dev/null ]; then
  systemctl start pufferd
  systemctl enable pufferd
else
  echo "systemd not installed, installing init.d script"
  echo "${initScript}" > /etc/init.d/pufferd
  chmod +x "/etc/init.d/pufferd"
  touch /var/log/pufferd.log
  chown pufferd:pufferd /var/log/pufferd.log
  update-rc.d pufferd defaults
  service pufferd start
fi

echo "Successfully installed the daemon"

exit 0

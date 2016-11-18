#!/bin/bash
# pufferd Installation Script

pufferdVersion=nightly

export DEBIAN_FRONTEND=noninteractive
downloadUrl="https://dl.pufferpanel.com/pufferd/${pufferdVersion}/pufferd"

RED=$(tput setf 4)
GREEN=$(tput setf 2)
NORMAL=$(tput sgr0)
BOLD=$(tput bold)

function checkResponseCode() {
    if [ $? -ne 0 ]; then
        echo -e "${RED}An error occured while installing, halting...${NORMAL}"
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
    apt-get install -y openssl curl openjdk-7-jdk tar python lib32gcc1 lib32tinfo5 lib32z1 lib32stdc++6
else
    yum -y install openssl curl java-1.8.0-openjdk-devel tar python glibc.i686 libstdc++.i686
fi
checkResponseCode

# Ensure /srv exists
mkdir -p /srv/pufferd
checkResponseCode

cd /srv/pufferd
curl -L -o pufferd $downloadUrl
checkResponseCode

chmod +x pufferd
./pufferd -install -auth {{ settings.master_url }} -token {{ node.daemon_secret }}
checkResponseCode

echo "Successfully Installed Scales"
exit 0

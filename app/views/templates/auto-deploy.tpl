#!/bin/bash
# pufferd Installation Script

pufferdVersion={{ pufferdVersion }}

pufferdRepo="pufferd"

if [ "${pufferdVersion}" == "nightly" ]; then
    pufferdRepo="pufferd-test"
fi

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

function writeServiceFile() {
    echo -e "Installing service file"
    echo "[Unit]
Description=pufferd daemon service

[Service]
Type=simple
WorkingDirectory=/var/lib/pufferd
ExecStart=${pufferdLocation}pufferd --run
ExecStop=${pufferdLocation}pufferd --shutdown $MAINPID
User=pufferd
Group=pufferd
TimeoutStopSec=2m
SendSIGKILL=no

[Install]
WantedBy=multi-user.target" > /lib/systemd/system/pufferd.service
}

if [ "$(id -u)" != "0" ]; then
    echo "This script must be run as root!" 1>&2
    exit 1
fi

if [ "$SUDO_USER" == "" ]; then
    SUDO_USER="root"
fi

# Determine distro type
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
elif type pacman &> /dev/null; then
    echo -e "System detected as Arch."
    OS_INSTALL_CMD="pacman"
else
    echo -e "${RED}This OS does not appear to be supported by this program, or apt-get/yum is not installed!${NORMAL}"
    exit 1
fi

# Uninstall if already installed
if [ -f /srv/pufferd/pufferd ] || [ -f /usr/sbin/pufferd ]; then
    echo "${red}WARNING: pufferd is already installed, continuing will DELETE the current pufferd installation and ALL SERVER FILES${normal}"
    echo "It is highly recommended that you back up all data in ${bold}/var/lib/pufferd${normal} prior to reinstalling"
    shopt -s nocasematch
    echo -n "Are you sure you wish DELETE ALL SERVER FILES and reinstall pufferd? [y/N]: "
    read installOverride
    if [[ "${installOverride}" != "y" ]]; then
        exit
    fi
    if [ -f /srv/pufferd/pufferd ]; then
        /srv/pufferd/pufferd --uninstall
    elif [ -f /usr/sbin/pufferd ]; then
        if [ $OS_INSTALL_CMD == 'apt' ]; then
          apt-get remove --purge pufferd
        elif [ $OS_INSTALL_CMD == 'yum' ]; then
          yum remove pufferd
          rm -rf /etc/pufferd
        fi
    fi
fi

# Install Other Dependencies
echo "Installing some dependencies"
if [ $OS_INSTALL_CMD == 'apt' ]; then
    if [ $(lsb_release -sc) == 'jessie' ]; then
        sudo echo "deb http://http.debian.net/debian jessie-backports main" > /etc/apt/sources.list.d/backports.list
        dpkg --add-architecture i386
        apt-get update
        apt-get install -y -t jessie-backports openjdk-8-jdk-headless
        apt-get install -y openssl curl git tar lib32gcc1 lib32tinfo5 lib32z1 lib32stdc++6 libcurl3-gnutls:i386
    elif [ $(lsb_release -sc) == 'trusty' ]; then
        sudo add-apt-repository -y ppa:openjdk-r/ppa
        dpkg --add-architecture i386
        apt-get update
        apt-get install -y openssl curl git openjdk-8-jdk-headless tar lib32gcc1 lib32tinfo5 lib32z1 lib32stdc++6 libcurl3-gnutls:i386
    else
        dpkg --add-architecture i386
        apt-get update
        apt-get install -y openssl curl git openjdk-8-jdk-headless tar lib32gcc1 lib32tinfo5 lib32z1 lib32stdc++6 libcurl3-gnutls:i386
    fi
    curl -s https://packagecloud.io/install/repositories/pufferpanel/${pufferdRepo}/script.deb.sh | bash
elif [ $OS_INSTALL_CMD == 'yum' ]; then
    yum -y install openssl curl git java-1.8.0-openjdk-devel tar glibc.i686 libstdc++.i686 libcurl.i686
    curl -s https://packagecloud.io/install/repositories/pufferpanel/${pufferdRepo}/script.rpm.sh | bash
elif [ $OS_INSTALL_CMD == 'pacman' ]; then
    grep -e "^\[multilib\]$" /etc/pacman.conf &> /dev/null
    if [ $? -eq 0 ]; then
        pacman -S openssl curl git jdk8-openjdk tar lib32-glibc lib32-gcc-libs --noconfirm --needed
    else
        echo -e "Please enable [multilib] in /etc/pacman.conf for lib32 libraries"
    fi
fi

mkdir -p /var/lib/pufferd /var/log/pufferd /etc/pufferd

echo -e "Installing pufferd using package manager"
pufferdLocation="/srv/pufferd"
installed=0
if [ $OS_INSTALL_CMD == 'apt' ]; then
    apt-get update
    apt-get install pufferd
    pufferdLocation="/usr/sbin/"
elif [ $OS_INSTALL_CMD == 'yum' ]; then
    yum install -y pufferd
    pufferdLocation="/usr/sbin/"
fi

if [ -f "${pufferdLocation}/pufferd" ]; then
    echo "Detected installation via package successful"
else
    echo -e "Failed to install using package manager, manually installing"
    echo -e "Downloading pufferd from $downloadUrl"
    pufferdLocation="/srv/pufferd/"
    mkdir -p /srv/pufferd
    curl -L -o /srv/pufferd/pufferd $downloadUrl
    checkResponseCode
    chmod +x /srv/pufferd/pufferd
    checkResponseCode
    writeServiceFile
    checkResponseCode
    useradd --system --home /var/lib/pufferd --user-group pufferd
fi

if type systemctl &> /dev/null; then
    echo "Stopping service to prepare for installation"
    systemctl stop pufferd
elif type service &> /dev/null; then
    echo "Stopping service to prepare for installation"
    service pufferd stop
fi

cd $pufferdLocation
echo -e "Executing pufferd installation"
./pufferd --install --auth {{ settings.master_url }} --token {{ node.daemon_secret }} --config /etc/pufferd/config.json
checkResponseCode

chown -R pufferd:pufferd /var/lib/pufferd /etc/pufferd /var/log/pufferd
if [ -f /srv/pufferd ]; then
  chown -R pufferd:pufferd /srv/pufferd
fi

echo "Preparing for docker containers if enabled"
groupadd --force --system docker
usermod -a -G docker pufferd

if type systemctl &> /dev/null; then
    echo "Starting pufferd service"
    systemctl start pufferd
    systemctl enable pufferd
elif type service &> /dev/null; then
    echo "Starting pufferd service"
    service pufferd start
fi

echo "Successfully installed the daemon"

exit 0

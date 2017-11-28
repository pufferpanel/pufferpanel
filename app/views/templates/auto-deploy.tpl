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

if [ "$(id -u)" != "0" ]; then
    echo "This script must be run as root!" 1>&2
    exit 1
fi

if [ "$SUDO_USER" == "" ]; then
    SUDO_USER="root"
fi

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
        /srv/pufferd/pufferd -uninstall
    else
        /usr/sbin/pufferd -uninstall
    fi
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
elif type pacman &> /dev/null; then
    echo -e "System detected as Arch."
    OS_INSTALL_CMD="pacman"
else
    echo -e "${RED}This OS does not appear to be supported by this program, or apt-get/yum is not installed!${NORMAL}"
    exit 1
fi

# Install Other Dependencies
echo "Installing some dependiencies."
if [ $OS_INSTALL_CMD == 'apt' ]; then
    curl -s https://packagecloud.io/install/repositories/pufferpanel/${pufferdRepo}/script.deb.sh | bash
    if [ $(lsb_release -sc) == 'jessie' ]; then
        sudo echo "deb http://http.debian.net/debian jessie-backports main" > /etc/apt/sources.list.d/backports.list
        apt-get update
        apt-get install -y -t jessie-backports openjdk-8-jdk-headless
        apt-get install -y openssl curl git tar lib32gcc1 lib32tinfo5 lib32z1 lib32stdc++6
    elif [ $(lsb_release -sc) == 'trusty' ]; then
        sudo add-apt-repository -y ppa:openjdk-r/ppa
        apt-get update
        apt-get install -y openssl curl git openjdk-8-jdk-headless tar lib32gcc1 lib32tinfo5 lib32z1 lib32stdc++6
    else
        apt-get update
        apt-get install -y openssl curl git openjdk-8-jdk-headless tar lib32gcc1 lib32tinfo5 lib32z1 lib32stdc++6
    fi
elif [ $OS_INSTALL_CMD == 'yum' ]; then
    curl -s https://packagecloud.io/install/repositories/pufferpanel/${pufferdRepo}/script.rpm.sh | bash
    yum -y install openssl curl git java-1.8.0-openjdk-devel tar glibc.i686 libstdc++.i686
elif [ $OS_INSTALL_CMD == 'pacman' ]; then
    grep -e "^\[multilib\]$" /etc/pacman.conf &> /dev/null
    if [ $? -eq 0 ]; then
        pacman -S openssl curl git jdk8-openjdk tar lib32-glibc lib32-gcc-libs --noconfirm --needed
    else
        echo -e "Please enable [multilib] in /etc/pacman.conf for lib32 libraries"
    fi
fi

mkdir /var/lib/pufferd /var/log/pufferd

echo -e "Installing pufferd using package manager"
pufferdLocation="/srv/pufferd"
if [ $OS_INSTALL_CMD == 'apt' ]; then
    apt-get update
    apt-get install pufferd
    pufferdLocation="/usr/sbin/"
elif [ $OS_INSTALL_CMD == 'yum' ]; then
    yum install -y pufferd
    pufferdLocation="/usr/sbin/"
else
    echo -e "Downloading pufferd from $downloadUrl"
    mkdir -p /srv/pufferd
    curl -L -o /srv/pufferd/pufferd $downloadUrl
    checkResponseCode
fi

cd $pufferdLocation
echo -e "Executing pufferd installation"
chmod +x pufferd
./pufferd --install --auth {{ settings.master_url }} --token {{ node.daemon_secret }} --config /etc/pufferd/config.json
checkResponseCode

chown -R pufferd:pufferd /var/lib/pufferd /etc/pufferd /var/log/pufferd
checkResponseCode

echo "Preparing for docker containers if enabled"
groupadd --force --system docker
usermod -a -G docker pufferd

if [ -f /srv/pufferd ]; then
    chown -R pufferd:pufferd /srv/pufferd
fi

if type systemctl &> /dev/null; then
    echo "Starting service"
  systemctl start pufferd
  systemctl enable pufferd
fi

echo "Successfully installed the daemon"

exit 0
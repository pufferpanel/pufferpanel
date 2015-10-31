#!/bin/bash
# PufferPanel Installer Script

if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root!" 1>&2
   exit 1
fi

if [ $SUDO_USER == "" ]; then
	SUDO_USER="root"
fi

DEBIAN_FRONTEND=noninteractive
RED=$(tput setf 4)
NORMAL=$(tput sgr0)
SSL_COUNTRY="US"
SSL_STATE="New-York"
SSL_LOCALITY="New-York"
SSL_ORG="PufferPanel"
SSL_ORG_NAME="SSL"
SSL_EMAIL="auto-generate@ssl.example.com"
SSL_COMMON="{{ node.fqdn }}"
SSL_PASSWORD=""

function checkResponseCode() {
    if [ $? -ne 0 ]; then
        echo "${RED}An error occured while installing, removing docker and halting...${NORMAL}"
        apt-get autoremove --purge docker-engine
        rm -rf /var/lib/docker
        exit 1
    fi
}

function spinner() {
    local pid=$1
    local delay=0.1
    local spinstr='|/-\'
    while [ "$(ps a | awk '{print $1}' | grep $pid)" ]; do
        local temp=${spinstr#?}
        printf "%c" "$spinstr"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
        printf "\b\b\b\b\b\b"
    done
    printf "    \b\b\b\b"
}

# Install NodeJS Dependencies
curl -sL https://deb.nodesource.com/setup_4.x | bash -
checkResponseCode

# Install Docker
curl -sSL https://get.docker.com/ | sh
checkResponseCode

# Install Other Dependencies
echo "Installing some dependiencies. Please Wait... "
(apt-get install -y openssl curl git make gcc g++ nodejs openjdk-7-jdk tar > /dev/null 2>&1 || true) &
spinner $!
checkResponseCode

# Add your user to the docker group
echo "Configuring Docker for:" $SUDO_USER
usermod -aG docker $SUDO_USER
checkResponseCode

# Add the Scales User Group
addgroup --system scalesuser
checkResponseCode

# Change the SFTP System
sed -i '/Subsystem sftp/c\Subsystem sftp internal-sftp' /etc/ssh/sshd_config
checkResponseCode

# Add Match Group to the End of the File
echo -e "Match group scalesuser
    ChrootDirectory %h
    X11Forwarding no
    AllowTcpForwarding no
    ForceCommand internal-sftp" >> /etc/ssh/sshd_config
checkResponseCode

# Restart SSHD
service ssh restart
checkResponseCode

# Ensure /srv exists
mkdir -p /srv
checkResponseCode

# Clone the repository
git clone https://github.com/PufferPanel/Scales /srv/scales
checkResponseCode

cd /srv/scales
checkResponseCode

# Checkout the Latest Version of Scales
git checkout tags/$(git describe --abbrev=0 --tags)
checkResponseCode

# Install the dependiencies for Scales to run.
# This process may take a few minutes to complete.
npm install
checkResponseCode

# Generate SSL Certificates
openssl req -x509 -days 365 -newkey rsa:4096 -keyout https.key -out https.pem -nodes -passin pass:$SSL_PASSWORD \
    -subj "/C=$SSL_COUNTRY/ST=$SSL_STATE/L=$SSL_LOCALITY/O=$SSL_ORG/OU=$SSL_ORG_NAME/CN=$SSL_COMMON/emailAddress=$SSL_EMAIL"
checkResponseCode

echo '{
	"listen": {
		"sftp": {{ node.daemon_sftp }},
		"rest": {{ node.daemon_listen }},
		"socket": {{ node.daemon_console }},
		"uploads": {{ node.daemon_upload }}
	},
	"urls": {
		"repo": "{{ settings.master_url }}auth/remote/pack",
		"download": "{{ settings.master_url }}auth/remote/download",
		"install": "{{ settings.master_url }}auth/remote/install-progress"
	},
	"ssl": {
		"key": "https.key",
		"cert": "https.pem"
	},
	"basepath": "{{ node.daemon_base_dir }}",
	"keys": [
		"{{ node.daemon_secret }}"
	],
	"upload_maxfilesize": 100000000
}' > config.json
checkResponseCode

npm start
checkResponseCode

echo "Successfully Installed Scales"
exit 0

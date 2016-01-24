#!/bin/bash
# PufferPanel Installer Script

export DEBIAN_FRONTEND=noninteractive
scalesApt=http://ci.pufferpanel.com:8080/browse/SC-DE/latestSuccessful/artifact/shared/Debian-Bundle/scales.tar.gz
scalesYum=http://ci.pufferpanel.com:8080/browse/SC-C7/latestSuccessful/artifact/shared/CentOS-Bundle/scales.tar.gz

if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root!" 1>&2
   exit 1
fi

if [ "$SUDO_USER" == "" ]; then
	SUDO_USER="root"
fi

RED=$(tput setf 4)
GREEN=$(tput setf 2)
NORMAL=$(tput sgr0)
BOLD=$(tput bold)
KERNEL=$(uname -r)

{% if node.docker == 1 %}
echo -e "${RED}${BOLD}[!!] STOP - READ THIS BEFORE CONTINUING [!!]${NORMAL}"
echo -e "THIS SOFTWARE DOES NOT AND WILL NOT RUN PROPERLY ON NON-STANDARD KERNELS. PLEASE ENSURE THAT THE OUTPUT BELOW IS VALID."
echo -e ""
echo -e "${BOLD}${GREEN}Kernel Version: ${KERNEL}${NORMAL}"
echo -e ""
echo -e "If this looks anything like '-grsec-xxxx-grs-ipv6-64' then it is probably a non-standard kernel. Standard kernels appear as '3.13.0-37-generic' or similar and should be at least version 3.10 or higher. Documentation for updating your kernel can be found at: http://scales.pufferpanel.com/docs/switching-ovh-kernels"
echo
read -r -p "I have read the above and confirmed that my kernel is a.) standard and b.) of a high enough version [y/N]: " response
case $response in
    [yY][eE][sS]|[yY])
        ;;
    *)
        exit 0
        ;;
esac
{% endif %}

if type apt-get &> /dev/null; then
    if [[ -f /etc/debian_version ]]; then
        echo -e "System detected as some variant of Ubuntu or Debian."
        OS_INSTALL_CMD="apt-get"
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
        echo -e "${RED}An error occured while installing, halting...${NORMAL}"
        exit 1
    fi
}

# Install Other Dependencies
echo "Installing some dependiencies."
if [ $OS_INSTALL_CMD == 'apt-get' ]; then
    apt-get install -y openssl curl git make gcc g++ nodejs openjdk-7-jdk tar python lib32gcc1 lib32tinfo5 lib32z1 lib32stdc++6
else
    yum -y install openssl curl git make gcc-c++ nodejs java-1.8.0-openjdk-devel tar python glibc.i686 libstdc++.i686
fi
checkResponseCode

{% if node.docker == 1 %}

# Install Docker
curl -sSL https://get.docker.com/ | sh
checkResponseCode

# Add your user to the docker group
echo "Configuring Docker for:" $SUDO_USER
usermod -aG docker $SUDO_USER
checkResponseCode

{% endif %}

# Add the Scales User Group
groupadd --system scalesuser

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
if [ $OS_INSTALL_CMD == 'apt-get' ]; then
    service ssh restart
else
    service sshd restart
fi
checkResponseCode

# Ensure /srv exists
mkdir -p /srv
checkResponseCode

cd /srv/
if [ $OS_INSTALL_CMD == 'apt-get' ]; then
    curl -o scales.tar.gz $scalesApt
else
    curl -o scales.tar.gz $scalesYum
fi
checkResponseCode

tar -xf scales.tar.gz
checkResponseCode
rm -f scales.tar.gz

cd /srv/scales
# Generate SSL Certificates
openssl req -x509 -days 365 -newkey rsa:4096 -keyout https.key -out https.pem -nodes -passin pass:$SSL_PASSWORD \
    -subj "/C=$SSL_COUNTRY/ST=$SSL_STATE/L=$SSL_LOCALITY/O=$SSL_ORG/OU=$SSL_ORG_NAME/CN=$SSL_COMMON/emailAddress=$SSL_EMAIL"
checkResponseCode

echo '{
	"listen": {
		"sftp": {{ node.daemon_sftp }},
		"rest": {{ node.daemon_listen }}
	},
	"urls": {
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
	"upload_maxfilesize": 100000000,
        "docker": {% if node.docker == 1 %} true {% else %} false {% endif %}
        
}' > config.json
checkResponseCode

chmod +x scales
./scales start
checkResponseCode

echo "Successfully Installed Scales"
exit 0

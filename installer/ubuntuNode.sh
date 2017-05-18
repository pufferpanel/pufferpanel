#!/bin/sh

#############################################################
#                                                           #
#                       PufferPanel                         #
#                                                           #
#               Ubuntu installer by Chomiciak               #
#                                                           #
# Use on the same license as Panel's (and on your own risk!)#
#                                                           #
#############################################################
#                                                           #
#                  This is Node installer                   #
#                                                           #
#############################################################

# All specified?
if [ -n "$0" ] && [ -n "$1" ] ; then
# Good. Let's get to work!

# Create folders!
mkdir /srv/modpacks/
mkdir /srv/scripts/
mkdir /srv/servers/

# Modpack directory chmod
chmod 755 /srv/modpacks/

# Curl install
apt-get -y install curl

# OpenSSL install
apt-get -y install openssl

# git install
apt-get -y install git

# RSSH install
apt-get -y install rssh

# quota install
apt-get -y install quota

# gcc install
apt-get -y install gcc

# make install
apt-get -y install make

# NodeJS install
apt-get -y install nodejs

# CPUlimit install
git clone https://github.com/DaneEveritt/cpulimit.git
cd cpulimit/
make
cp cpulimit /usr/bin

# Add user and user's password. Authorization mode = password, not key
useradd $0
passwd $0
$1

# Add "some_username ALL=(ALL) NOPASSWD: ALL" line to visudo file
echo $0" ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# RSSH setup
groupadd rsshusers

# add something to /etc/ssh/sshd_config
echo "Match Group rsshusers" >> /etc/ssh/sshd_config
echo " ChrootDirectory /srv/servers/%u" >> /etc/ssh/sshd_config
echo " X11Forwarding no" >> /etc/ssh/sshd_config
echo " AllowTcpForwarding no" >> /etc/ssh/sshd_config
echo " ForceCommand internal-sftp" >> /etc/ssh/sshd_config

#Change Subsystem sftp /usr/libexec/openssh/sftp-server to Subsystem sftp internal-sftp 
sed -i 's/Subsystem sftp \/usr\/libexec\/openssh\/sftp-server/Subsystem sftp internal-sftp/' /etc/ssh/sshd_config

#The default action for rssh to lock down everything. Granting access sftp open to the RSSH.
echo "allowsftp" >> /etc/rssh.conf

# GSD installing
cd /srv/ && git clone https://github.com/gametainers/gsd.git
cd /srv/gsd
npm install

# Restart SSHD
service sshd restart

# End data
echo "Copy scripts folder content to /srv/scripts/"
echo "and then..."
echo "type 'cd /srv/gsd && npm start' to start GSD!"

else

# Notification when params not specified
echo "Please specify User and Password:"
echo "Add following params: (after this filename, eg. ./ubuntuNode.sh USER PASSWORD)"
echo "USER PASSWORD"
echo ""
echo "And remember: you must start this as ROOT!!!"

fi

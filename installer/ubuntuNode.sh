#!/bin/sh

#############################################################
#                                                           #
#                       PufferPanel                         #
#                                                           #
#               Ubuntu installer by Chomiciak               #
#                                                           #
#            Use on the same license as Panel's             #
#                                                           #
#############################################################
#                                                           #
#                  This is Node installer                   #
#                                                           #
#############################################################

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
gmake
make
cp cpulimit /usr/bin

# Add user and user's password. Authorization mode = password, not key
useradd $1
passwd $1
$2

# Add "some_username ALL=(ALL) NOPASSWD: ALL" line to visudo file
# how?

# RSSH setup
groupadd rsshusers

# add something to /etc/ssh/sshd_config
echo "Match Group rsshusers" >> /etc/ssh/sshd_config
echo " ChrootDirectory /srv/servers/%u" >> /etc/ssh/sshd_config
echo " X11Forwarding no" >> /etc/ssh/sshd_config
echo " AllowTcpForwarding no" >> /etc/ssh/sshd_config
echo " ForceCommand internal-sftp" >> /etc/ssh/sshd_config

#Change Subsystem sftp /usr/libexec/openssh/sftp-server to Subsystem sftp internal-sftp 
# (how???)

#The default action for rssh to lock down everything. Granting access sftp open to the RSSH.
echo "allowsftp" >> /etc/rssh.conf

# GSD installing
cd /srv/ && git clone https://github.com/gametainers/gsd.git
cd /srv/gsd
npm install
echo "type 'npm start' to start GSD!"


#restart sshd
service sshd restart

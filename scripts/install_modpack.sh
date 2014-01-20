#!/bin/bash
#
#   PufferPanel - A Minecraft Server Management Panel
#   Copyright (c) 2013 Dane Everitt
#
#   This program is free software: you can redistribute it and/or modify
#   it under the terms of the GNU General Public License as published by
#   the Free Software Foundation, either version 3 of the License, or
#   (at your option) any later version.
#
#   This program is distributed in the hope that it will be useful,
#   but WITHOUT ANY WARRANTY; without even the implied warranty of
#   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#   GNU General Public License for more details.
#
#   You should have received a copy of the GNU General Public License
#   along with this program.  If not, see http://www.gnu.org/licenses/.
#
# PufferPanel Plugin Install Script
# Parameters: 
#	$1 = user
#	$2 = download
#	$3 = modpack zip

# Download the Mod
cd /tmp
curl -o "$3" "$2"

# Move into Server Directory
cd /srv/servers/$1/server

# Remove Possible Old Installer Files & Unzip into Directory
rm -rf pufferpanel_modpack_installer
unzip -o /tmp/$3 -x __MACOSX/* -d pufferpanel_modpack_installer

# Move into Directory & Apply Ownership
cd pufferpanel_modpack_installer
chown -R $1:$1 *

# Move all Files Out
rsync -rgop * ../
cd /srv/servers/$1/server

# Cleanup
rm -rf pufferpanel_modpack_installer
rm -rf /tmp/$3
#!/bin/bash
#
# PufferPanel - A Minecraft Server Management Panel
# Copyright (c) 2014 PufferPanel
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program.  If not, see http://www.gnu.org/licenses/.

version=0.2.0

divider="-----"
installPath=".."

echo "PufferPanel Language Builder - v${version}"
echo $divider

while getopts "h?p:" opt; do
    case "$opt" in
    p)  
        installPath=$OPTARG
        ;;
    h)
        echo "PufferPanel Language Builder - v${version}"
        echo "Optional parameters: "
        echo "-p        | Path to PufferPanel installation (i.e PufferPanel)"
        exit;
        ;;
    esac
done

shift $((OPTIND-1))

[ "$1" = "--" ] && shift

echo "Cleaning up old language files"
rawPath="${installPath}/app/languages/raw/"
outputPath="${installPath}/app/languages/"
rm -f ${outputPath}*.json

echo "Generating language files"

cd $rawPath
for f in *.txt
do
    echo -n "    "
    filename=$(basename $f .txt)
    echo "Building ${filename}"
    # Yes, this looks bad, but only way for bash to work with it
    php -r '$content = file("$argv[1].txt"); $json = array();
        foreach($content as $line => $string){
            list($id, $lang) = explode(",", $string, 2);
            $json = array_merge($json, array(strtolower($id) => trim($lang)));
        }
        $fp = fopen("../$argv[1].json", "w+");
        fwrite($fp, json_encode($json, JSON_UNESCAPED_SLASHES | JSON_PRETTY_PRINT));
        fclose($fp); ' -- ${filename}
done

echo "Building complete"
cd $(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)
<?php
/*
    PufferPanel - A Minecraft Server Management Panel
    Copyright (c) 2013 Dane Everitt

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see http://www.gnu.org/licenses/.
 */

if(!isset($_GET['l']))
	die('No language input passed.');

if(!file_exists('../'.$_GET['l'].'.txt'))
	die('Language file does not exist.');

$content = file('../'.$_GET['l'].'.txt');
$json = array();

foreach($content as $line => $string){

	list($id, $lang) = explode(",", $string, 2);
	$json = array_merge($json, array(strtolower(str_replace(".", "_", $id)) => trim($lang)));

}

$fp = fopen('../../../src/core/lang/'.$_GET['l'].'.json', 'w+');
fwrite($fp, json_encode($json, JSON_UNESCAPED_SLASHES | JSON_PRETTY_PRINT));
fclose($fp);

echo 'Language written to file.';

 ?>

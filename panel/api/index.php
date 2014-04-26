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

header('Content-Type: application/json');

//%7B%22auth%22%3A%7B%22key%22%3A%22random_key%22%7D%2C%22function%22%3A%22add%22%2C%22data%22%3A%7B%22server_name%22%3A%22api_test%22%2C%22node%22%3A1%2C%22modpack%22%3A%22default%22%2C%22email%22%3A%22dane%40daneeveritt.com%22%2C%22server_ip%22%3A%22127.0.0.1%22%2C%22server_port%22%3A25570%2C%22alloc_mem%22%3A128%2C%22alloc_disk%22%3A1024%2C%22sftp_pass%22%3A%22password%22%2C%22sftp_pass_2%22%3A%22password%22%2C%22cpu_limit%22%3A0%7D%7D

/*
array(
	"auth" => array(
		"key" => "random_key"
	),
	"function" => "add",
	"data" => array(
		"server_name" => "api_test",
		"node" => 1,
		"modpack" => "default",
		"email" => "dane@daneeveritt.com",
		"server_ip" => "127.0.0.1",
		"server_port" => 25570,
		"alloc_mem" => 128,
		"alloc_disk" => 1024,
		"sftp_pass" => "password",
		"sftp_pass_2" => "password",
		"cpu_limit" => 0,
	)
)
*/

require_once('functions/api.core.php');

$api->init();
?>
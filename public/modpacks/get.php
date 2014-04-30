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
require_once('../core/framework/framework.database.connect.php');

/*
 * MySQL PDO Connection Engine
 */
$mysql = Database\database::connect();

function pdo_exception_handler($exception) {
    if ($exception instanceof PDOException) {
        
        error_log($exception);
        
        die(json_encode(array('error' => 'A MySQL error was encountered with this request.', 'e_code' => $exception->getCode(), 'e_line' => $exception->getLine(), 'e_time' => date('d-M-Y H:i:s', time()))));
        
    } else {
    
    	die('Exception handler from unknown source.');
    
    }
}
set_exception_handler('pdo_exception_handler');

if(!isset($_GET['pack']))
	exit('no parameters');

list($encrypted, $iv) = explode('.', rawurldecode($_GET['pack']));

$dlHash = openssl_decrypt($encrypted, 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($iv));


$query = $mysql->prepare("SELECT * FROM `modpacks` WHERE `download_hash` = :dlhash");
$query->execute(array(
	':dlhash' => $dlHash
));

$getSettings = $mysql->prepare("SELECT `setting_val` FROM `acp_settings` WHERE `setting_ref` = 'modpack_dir'");
$getSettings->execute();

function download($filename) { 
		
	$chunksize = 1*(1024*1024);
	$buffer = '';
	$handle = fopen($filename, 'rb');
	
		if ($handle === false){
		
			return false;
			
		}

		while (!feof($handle)){
		
			$buffer = fread($handle, $chunksize);
			print $buffer;
			
		}
		
	return fclose($handle);
	
}

if($query->rowCount() == 1 && $getSettings->rowCount() == 1){

	$row = $query->fetch();
	$setting = $getSettings->fetch();
	
	header("Pragma: private");
	header("Expires: 0");
	header("Cache-Control: must-revalidate, post-check=0, pre-check=0");
	header("Content-Type: application/force-download");
	header("Content-Description: File Transfer");
	header('Content-Disposition: response; filename="'.$row['hash'].'.zip"');
	header("Content-Transfer-Encoding: binary");
	header('Accept-Ranges: bytes');
	header("Content-Length: ".filesize($setting['setting_val'].$row['hash'].'.zip'));
	    
	download($setting['setting_val'].$row['hash'].'.zip');
	
}else{

	exit('no pack exists');

}
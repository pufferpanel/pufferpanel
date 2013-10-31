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
#For Hashing Node & Server FTP/SSH passwords
$eHash = file_get_contents('HASHFILE.txt');

$iv = mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_CAST_256, MCRYPT_MODE_CBC), MCRYPT_RAND);

echo 'encryption_iv: '.base64_encode($iv).'<br />';
echo 'encrypted_pass: '.openssl_encrypt($_GET['openssl'], 'AES-256-CBC', $eHash, false, $iv);

#For Creating User Account Passwords
#CHANGE THE SALT THISISMYSALT_CHANGE_ME! It can contain any character except $. Do not remove the $ already there.
$salt = crypt($_GET['ripemd'], '$6$rounds=5000$THISISMYSALT_CHANGE_ME$');
echo '<br />Userpass: '.hash('ripemd320', $salt);


$keyset  = "abcdefghijklmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
$randkey = "";

for ($i=0; $i<42; $i++)
	$randkey .= substr($keyset, rand(0, strlen($keyset)-1), 1);

echo '<br /><br />Server Hash: '.$randkey;

?>
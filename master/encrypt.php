<?php

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
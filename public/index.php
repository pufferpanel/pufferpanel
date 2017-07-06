<?php
$baseUrl = !file_exists(dirname(__FILE__).'/pufferpanel') ? dirname($_SERVER['PHP_SELF']).'/' : '/';
header('Location: '.$baseUrl.'index');

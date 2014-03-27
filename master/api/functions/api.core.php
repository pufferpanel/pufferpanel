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
require_once('../../core/framework/framework.database.connect.php');
require_once('../../core/framework/framework.auth.php');
require_once('../../core/framework/framework.ssh2.php');
require_once('../../core/framework/framework.query.php');
require_once('../../core/framework/framework.settings.php');
require_once('api.init.php');

$api = new apiInitializationClass(); 
$api->auth = new apiAuthenticationClass();
$api->process = new apiProcessingClass();
$api->run = new apiModuleRunClass();

?>
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
require_once('../core/framework/framework.core.php');
require_once('api.init.php');
require_once('modules/module.add.php');

$api = new \API\apiInitializationClass();
$api->module = new stdClass();

$api->module->add = new \Modules\Add\apiModuleAddServer();
#$api->auth = new apiAuthenticationClass();
#$api->process = new apiProcessingClass();
#$api->run = new apiModuleRunClass();

?>
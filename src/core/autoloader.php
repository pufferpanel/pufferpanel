<?php
/*
	PufferPanel - A Minecraft Server Management Panel
	Copyright (c) 2015 PufferPanel

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

/*
 * Include composer and configuration files
 */
require_once(BASE_DIR.'vendor/autoload.php');

/*
* Include Required Global Component Files
*/
require_once(SRC_DIR.'core/components/authentication.php');
require_once(SRC_DIR.'core/components/errorhandler.php');
require_once(SRC_DIR.'core/components/functions.php');
require_once(SRC_DIR.'core/components/gsd.php');
require_once(SRC_DIR.'core/components/page.php');

/*
 * Include Required Global Class Files
 */
require_once(SRC_DIR.'core/authentication.php');
require_once(SRC_DIR.'core/email.php');
require_once(SRC_DIR.'core/files.php');
require_once(SRC_DIR.'core/language.php');
require_once(SRC_DIR.'core/user.php');
require_once(SRC_DIR.'core/log.php');
require_once(SRC_DIR.'core/query.php');
require_once(SRC_DIR.'core/server.php');
require_once(SRC_DIR.'core/settings.php');
require_once(SRC_DIR.'core/routes.php');
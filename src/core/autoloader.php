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
require_once(BASE_DIR . 'vendor/autoload.php');

\Unirest\Request::verifyPeer(false);
\Unirest\Request::timeout(5);

/*
 * Include Required Global Component Files
 */
require_once(SRC_DIR . 'core/components/authentication.php');
require_once(SRC_DIR . 'core/components/errorhandler.php');
require_once(SRC_DIR . 'core/components/functions.php');
require_once(SRC_DIR . 'core/components/daemon.php');
require_once(SRC_DIR . 'core/components/page.php');
require_once(SRC_DIR . 'core/components/url.php');

/**
 * Load config types
 */
require_once(SRC_DIR . 'core/config/configinterface.php');
require_once(SRC_DIR . 'core/config/jsonconfig.php');
require_once(SRC_DIR . 'core/config/databaseconfig.php');
require_once(SRC_DIR . 'core/config/inmemdatabaseconfig.php');

/*
 * Include Required Global Class Files
 */
require_once(SRC_DIR . 'core/config.php');
require_once(SRC_DIR . 'core/databasemanager.php');
require_once(SRC_DIR . 'core/authentication.php');
require_once(SRC_DIR . 'core/email.php');
require_once(SRC_DIR . 'core/files.php');
require_once(SRC_DIR . 'core/user.php');
require_once(SRC_DIR . 'core/language.php');
require_once(SRC_DIR . 'core/log.php');
require_once(SRC_DIR . 'core/permissions.php');
require_once(SRC_DIR . 'core/server.php');
require_once(SRC_DIR . 'core/daemon.php');
require_once(SRC_DIR . 'core/settings.php');
require_once(SRC_DIR . 'core/routes.php');
require_once(SRC_DIR . 'core/version.php');
require_once(SRC_DIR . 'core/oauth2.php');

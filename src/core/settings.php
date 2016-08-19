<?php

/*
  PufferPanel - A Game Server Management Panel
  Copyright (c) 2015 Dane Everitt

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

namespace PufferPanel\Core;

use PufferPanel\Core\Config\InMemoryDatabaseConfig;

class Settings {

    private static $settings;

    /**
     * Gets a specific key from the global config, or the entire config if key is null.
     *
     * @param string $base Key
     */
    public static function config($base = null) {

        self::checkStatus();
        return self::$settings->config($base);
    }

    /**
     * Returns the internal Config instance of the global config.
     *
     * @return Core\ConfigInterface
     */
    public static function getInstance() {

        self::checkStatus();
        return self::$settings;
    }

    /**
     * Checks to make sure the config is loaded.
     */
    protected static function checkStatus() {

        if (is_null(self::$settings)) {
            self::$settings = new InMemoryDatabaseConfig('acp_settings', 'setting_ref', 'setting_val');
        }
    }

}

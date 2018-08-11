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

namespace PufferPanel\Core\Config;

use \Exception;

class JsonConfig implements ConfigInterface {

    private $config;

    /**
     * Constructor class for implementing configuration files from JSON.
     *
     * @param string $path Configuration file relative to BASE_DIR
     * @param bool $array
     */
    public function __construct($path, $array = false) {
        if (!file_exists(BASE_DIR . $path)) {
            throw new Exception("The config file " . $path . " does not exist.");
        }

        $this->config = json_decode(file_get_contents(BASE_DIR . $path), $array);

        if (json_last_error() != "JSON_ERROR_NONE") {
            $errMsg = 'Unknown error';

            switch (json_last_error()) {
                case JSON_ERROR_DEPTH:
                    $errMsg = 'Maximum stack depth exceeded';
                    break;
                case JSON_ERROR_STATE_MISMATCH:
                    $errMsg = 'Underflow or the modes mismatch';
                    break;
                case JSON_ERROR_CTRL_CHAR:
                    $errMsg = 'Unexpected control character found';
                    break;
                case JSON_ERROR_SYNTAX:
                    $errMsg = 'Syntax error, malformed JSON';
                    break;
                case JSON_ERROR_UTF8:
                    $errMsg = 'Malformed UTF-8 characters, possibly incorrectly encoded';
                    break;
            }


            throw new Exception("An error occurred when trying decode " . $path . ". " . $errMsg);
        }
    }

    public function config($base = null) {
        return (is_null($base)) ? $this->config : property_exists($this->config, $base) ? $this->config->{$base} : null;
    }

}

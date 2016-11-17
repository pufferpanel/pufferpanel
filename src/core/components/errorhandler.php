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

namespace PufferPanel\Core\Components;

/**
 * PufferPanel Core Components Trait
 */
trait Error_Handler {

    protected static $_error;

    protected final static function _setError($error) {

        self::$_error = $error;
    }

    public function retrieveLastError($wrapped = true) {

        return (!$wrapped) ? self::$_error : '<div class="alert alert-danger">' . self::$_error . '</div>';
    }

}

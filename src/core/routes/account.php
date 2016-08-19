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

namespace PufferPanel\Core\Router;

use \ORM,
    PufferPanel\Core\Settings;

class Account extends \PufferPanel\Core\Email {

    protected $_user;

    /**
     * Constructor class that requires the user class be passed as a variable. This is done mostly because
     * I can't just extend \PufferPanel\Core\User for some reason, and I am too lazy to fix that right now.
     *
     * @param object $user
     */
    public function __construct($user) {

        $this->_user = $user;
    }

    /**
     * Handles updating a user password for the router. Can also be used outside of the router if needed.
     *
     * @param string $new
     * @return bool
     */
    public function updatePassword($new) {

        $new_password = $this->hash($new);

        $this->account = ORM::forTable('users')->findOne($this->_user->getData('id'));

        if (!$this->account) {
            return false;
        }

        $this->account->password = $new_password;
        $this->account->session_id = null;
        $this->account->session_ip = null;
        $this->account->save();

        $this->buildEmail('password_changed', array(
            'IP_ADDRESS' => $_SERVER['REMOTE_ADDR'],
            'GETHOSTBY_IP_ADDRESS' => gethostbyaddr($_SERVER['REMOTE_ADDR'])
        ))->dispatch($this->_user->getData('email'), Settings::config()->company_name . ' - Password Change Notification');

        return true;
    }

}

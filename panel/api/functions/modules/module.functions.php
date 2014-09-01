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

trait functions {

	use \Functions\general;

	public static function throwResponse($text, $success = false){

		exit(json_encode(
			array(
				'success' => $success,
				'response' => $text
			)
		));

	}

	public static function getStoredData() {

		if(!isset($_POST['request']))
			self::throwResponse("No data was sent in the POST request.");
		else
			return json_decode($_POST['request'], true);

	}

    private function validateUserbyEmail() {

        $this->email = self::getStoredData()['data']['email'];

        $this->select = $this->mysql->prepare("SELECT `id` FROM `users` WHERE `email` = :email");
        $this->select->execute(array(
            ':email' => $this->email
        ));

            if($this->select->rowCount() != 1)
                self::throwResponse('No user is associated with that email.', false);
            else
                $this->oid = $this->select->fetch()['id'];

    }

    private function validateNode() {

        $this->node = self::getStoredData()['data']['node'];

        $this->select = $this->mysql->prepare("SELECT `id` FROM `nodes` WHERE `id` = :id");
        $this->select->execute(array(
            ':id' => $this->node
        ));

            if($this->select->rowCount() != 1)
                self::throwResponse('No node is associated with that ID.', false);

    }

}

?>

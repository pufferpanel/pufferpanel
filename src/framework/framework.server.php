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

/**
 * PufferPanel Core Server management class.
 */
class server extends user {

	use Page\components;

    /**
     * Constructor class for building server data.
     *
     * @param string $hash The server hash.
     * @param int $userid The ID of the user who is requesting the server information.
     * @param int $isroot The root administrator status of the user requesting the server information.
     * @return array|bool Returns an array on success or false on failure.
     */
	public function __construct($hash = null, $userid = null, $isroot = null){

		$this->mysql = self::connect();

		if(!is_null($userid) && $userid === true && !is_null($hash))
			$this->_buildData($hash, $userid, $isroot);
        else if(!is_null($userid) && is_null($hash) && is_null($isroot))
            $this->_rebuildData($userid);
        else
            $this->_s = false;

	}

    /**
     * Re-runs the __construct() class with a defined ID for the admin control panel.
     *
     * @param int $id This value should be the ID of the server you are getting information for.
     * @return void
     */
    public function rebuildData($id){

            $this->__construct(null, $id);

    }

    /**
     * Provides the corresponding value for the id provided from the MySQL Database.
     *
     * @param string $id The column value for the data you need (e.g. email).
     * @return string|bool A string is returned on success, if the command fails 'false' is returned.
     */
	public function getData($id){

		if($this->_s === true && array_key_exists($id, $this->_data))
			return $this->_data[$id];
		else
			return false;

	}

    /**
     * Returns data about the node in which the server selected is running.
     *
     * @param string $id The column value for the data you need (e.g. sftp_ip).
     * @return string|bool A string is returned on success, if the command fails 'false' is returned.
     */
    public function nodeData($id) {

        if($this->_n === true && array_key_exists($id, $this->_ndata))
			return $this->_ndata[$id];
		else
			return null;

    }

    /**
     * Handles incoming requests to access a server and redirects to the correct location and sets a cookie.
     *
     * @param string $id The column value for the data you need (e.g. sftp_ip).
     * @return void
     */
    public function nodeRedirect($hash, $userid, $rootAdmin) {

		if($rootAdmin == 1){

			$query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `hash` = ? AND `active` = '1'");
			$query->execute(array(
				$hash
			));

		}else{

			$query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `owner_id` = :ownerid AND `hash` = :hash AND `active` = '1'");
			$query->execute(array(
				':ownerid' => $userid,
				':hash' => $hash
			));

		}

			if($query->rowCount() == 1){

				$row = $query->fetch();

					setcookie('pp_server_hash', $row['hash'], 0, '/');

					$this->redirect('node/index.php');

			}else
				$this->redirect('servers.php?error=error');

	}

    /**
     * Rebuilds server data using a specified ID. Useful for Admin CP applications.
     *
     * @param int $userid The server ID.
     * @return array|bool Returns an array on success or false on failure.
     */
    private function _rebuildData($userid){

        $this->query->server = $this->mysql->prepare("SELECT * FROM `servers` WHERE `id` = :value");
        $this->query->server->execute(array(
            ':value' => $userid
        ));

        if($this->query->rowCount() == 1){

            $this->row = $this->query->fetch();

            foreach($this->row as $this->id => $this->val)
                $this->_data = array_merge($this->_data, array($this->id => $this->val));

        }else
            $this->_s = false;

        /*
         * Grab Node Information
         */
        if(isset($this->_data['node']) && $this->_data['node'] !== false){

            $this->_ndata = array();
            $this->_n = true;

            $this->query->node = $this->mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :node LIMIT 1");
            $this->query->node->execute(array(
                ':node' => $this->_data['node']
            ));

            if($this->query->node->rowCount() == 1){

                $this->node = $this->query->node->fetch();

                foreach($this->node as $this->nid => $this->nval)
                $this->_ndata = array_merge($this->_ndata, array($this->nid => $this->nval));

            }else
                $this->_n = false;

        }else
            $this->_n = false;

    }

    /**
     * Builds server data using a specified ID, Hash, and Root Administrator Status.
     *
     * @param string $hash The server hash.
     * @param int $userid The ID of the user who is requesting the server information.
     * @param int $isroot The root administrator status of the user requesting the server information.
     * @return array|bool Returns an array on success or false on failure.
     */
    private function _buildData($hash, $userid, $isroot){

        $this->_data = array();
        $this->_s = true;

            if($isroot == '1'){

                $this->query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `hash` = :hash AND `active` = 1");
                $this->query->execute(array(
                    ':hash' => $hash
                ));

            }else{

                $this->query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `hash` = :hash AND `owner_id` = :ownerid AND `active` = 1");
                $this->query->execute(array(
                    ':hash' => $hash,
                    ':ownerid' => $userid
                ));

            }

            if($this->query->rowCount() == 1){

                $this->row = $this->query->fetch();

                foreach($this->row as $this->id => $this->val)
                    $this->_data = array_merge($this->_data, array($this->id => $this->val));

            }else
                $this->_s = false;

        /*
         * Grab Node Information
         */
        if(isset($this->_data['node']) && $this->_data['node'] !== false){

            $this->_ndata = array();
            $this->_n = true;

            $this->query->node = $this->mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :node LIMIT 1");
            $this->query->node->execute(array(
                ':node' => $this->_data['node']
            ));

            if($this->query->node->rowCount() == 1){

                $this->node = $this->query->node->fetch();

                foreach($this->node as $this->nid => $this->nval)
                $this->_ndata = array_merge($this->_ndata, array($this->nid => $this->nval));

            }else
                $this->_n = false;

        }else
            $this->_n = false;

    }

}

?>

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

trait add {

	use \Database\database, functions;

	private function validateData() {

		/*
		 * @TODO: Validate Incoming Data
		 */
		$this->data = self::getStoredData()['data'];

		/*
		 * Validate Server Name
		 */
		if(!preg_match('/^[\w-]{4,35}$/', $this->data['server_name']))
			self::throwResponse('Invalid server name was provided.', false);

		/*
		 * Validate Email
		 */
		if(!filter_var($this->data['email'], FILTER_VALIDATE_EMAIL))
			self::throwResponse('Invalid email was provided.', false);

        /*
         * Validate Node
         */
        self::validateNode();

		/*
		 * Validate Disk & Memory
		 */
		if(!is_numeric($this->data['alloc_mem']) || !is_numeric($this->data['alloc_disk']))
			self::throwResponse('Invalid memory or disk allocation.', false);

		/*
		 * Validate CPU Limit
		 */
		if(!is_numeric($this->data['cpu_limit']))
			self::throwResponse('Invalid CPU limit.', false);

        /*
         * Validate Server IP & Port
         */
        /*
         * Determine if Node (IP & Port) is Avaliable
         */
        $this->validateServerIP();

		/*
		 * Validate User Exists
		 */
		self::validateUserbyEmail();

	}

	private function runRequest() {

        /*
         * Run Functions
         */
        $this->addToDatabase();
        $this->handlePorts();
        $this->handleModpack();
        $this->addToGSD();

        /*
         * Send User Email
         */
        $this->email = new tplMail();
        $this->settings = new settings();

        $this->email->buildEmail('admin_new_server', array(
                'NAME' => $this->data['server_name'],
                'CONNECT' => $this->node['ip'].':21',
                'USER' => $this->ftpUser.'-'.$this->gsdContent['id'],
                'PASS' => $this->rawPassword
        ))->dispatch($this->data['email'], $this->settings->get('company_name').' - New Server Added');

        self::throwResponse('added', true);

	}

    private function addToDatabase() {

        /*
         * Generate Variables
         */
        $this->ftpUser = self::generateFTPUsername($this->data['server_name']);
        $this->serverHash = self::gen_UUID();

        /*
         * Password Handling
         */
        $this->iv = $this->generate_iv();

        if(is_null($this->data['ftp_pass']))
            $this->rawPassword = self::keygen(14);
        else
            $this->rawPassword = $this->data['ftp_pass'];

        $this->password = $this->encrypt($this->rawPassword, $this->iv);

        /*
         * Add Server to Database
         */
        $this->add = $this->mysql->prepare("INSERT INTO `servers` VALUES(NULL, NULL, NULL, :hash, :gsd_secret, :e_iv, :node, :sname, :modpack, :sjar, 1, :oid, :ram, :disk, :cpu, :date, :sip, :sport, :ftpuser, :ftppass)");

        $this->add->execute(array(
            ':hash' => $this->serverHash,
            ':gsd_secret' => self::gen_UUID(),
            ':e_iv' => $this->iv,
            ':node' => $this->data['node'],
            ':sname' => $this->data['server_name'],
            ':modpack' => '0000-0000-0000-0',
            ':sjar' => "server.jar",
            ':oid' => $this->oid,
            ':ram' => $this->data['alloc_mem'],
            ':disk' => $this->data['alloc_disk'],
            ':cpu' => $this->data['cpu_limit'],
            ':date' => time(),
            ':sip' => $this->data['server_ip'],
            ':sport' => $this->data['server_port'],
            ':ftpuser' => $this->ftpUser,
            ':ftppass' => $this->password
        ));

    }

	private function addToGSD() {

		/*
		 * Add User to GSD
         */
        $data = json_encode(array(
            "name" => $this->serverHash,
            "user" => $this->ftpUser,
            "overide_command_line" => "",
            "path" => "/home/".$this->ftpUser,
            "variables" => array(
                "-jar" => "server.jar",
                "-Xmx" => $this->data['alloc_mem']."M"
            ),
            "gameport" => $this->data['server_port'],
            "gamehost" => "",
            "plugin" => "minecraft",
            "autoon" => false
        ));

        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, 'http://'.$this->node['ip'].':8003/gameservers');
        curl_setopt($ch, CURLOPT_HTTPHEADER, array(
            'X-Access-Token: '.$this->node['gsd_secret']
        ));
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 5);
        curl_setopt($ch, CURLOPT_TIMEOUT, 5);
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_POSTFIELDS, "settings=".$data);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        $this->gsdContent = json_decode(curl_exec($ch), true);
        curl_close($ch);

        /*
         * Update MySQL
         */
        $this->update = $this->mysql->prepare("UPDATE `servers` SET `gsd_id` = :gsdid WHERE `hash` = :hash");
        $this->update->execute(array(
            ':gsdid' => $this->gsdContent['id'],
            ':hash' => $this->serverHash
        ));

	}

	private function handlePorts() {

        /*
         * Update Port Information
         */
        $this->ips[$this->data['server_ip']]['ports_free']--;
        $this->ports[$this->data['server_ip']][$this->data['server_port']]--;

        $this->mysql->prepare("UPDATE `nodes` SET `ips` = :ips, `ports` = :ports WHERE `id` = :id")->execute(array(
            ':ips' => json_encode($this->ips),
            ':ports' => json_encode($this->ports),
            ':id' => $this->data['node']
        ));

	}

	private function handleModpack() {

        /*
         * Handle Modpack Information
         * @TODO: Await GSD implementation
         */

	}

    private function validateServerIP() {

        /*
         * Validates Selected IP and Port
         */
        $this->select = $this->mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :id");
        $this->select->execute(array(
            ':id' => $this->data['node']
        ));

        if($this->select->rowCount() == 1)
            $this->node = $this->select->fetch();
        else
            self::throwResponse('Unable to query database for the node provided.', false);

            /*
             * Validate IP & Port
             */
            $this->ips = json_decode($this->node['ips'], true);
            $this->ports = json_decode($this->node['ports'], true);

            if(!array_key_exists($this->data['server_ip'], $this->ips))
                self::throwResponse('Invalid server ip provided.', false);

            if(!array_key_exists($this->data['server_port'], $this->ports[$this->data['server_ip']]))
                self::throwResponse('Invalid server port provided.', false);

            if($this->ports[$this->data['server_ip']][$this->data['server_port']] == 0)
                self::throwResponse('The specified port is currently in use.', false);

    }

}

class apiModuleAddServer {

    use add;

    public function __construct() {

        $this->mysql = self::connect();

        $this->validateData();
        $this->runRequest();

    }

}

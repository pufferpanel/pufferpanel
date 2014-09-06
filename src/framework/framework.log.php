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

class log extends user {
    
    //use Database\database;
    private $url;
    
    public function __construct($uid)
        {
        
            $this->uid = ($uid !== false) ? $uid : null;
            $this->mysql = self::connect();
        
        }
    
    /*
     * Call as Such:
     * $core->log->getUrl()->addLog(priority, viewable, array(action, desc, uid*));
     *
     * priority: 0 - None, 1 - Low, 2 - Medium, 3 - High, 4 - Urgent
     * viewable: 0 - Admin Only, 1 - User & Admin
     * action: (example) user.login, user.start_server, admin.remove_server, admin.node_add_port, system.unknown_error
     * desc: Description of the Error
     * uid: Sent from login commands where $user is not yet defined. Optional.
     */
    public function addLog($priority, $viewable, $data = array())
        {
        
            $this->query = $this->mysql->prepare("INSERT INTO `actions_log` VALUES(NULL, :priority, :viewable, :user, :time, :ip, :url, :action, :desc)");
            
            $this->uid = (!array_key_exists(2, $data)) ? $this->uid : $data[2];
        
            $this->query->execute(array(
                ':priority' => $priority,
                ':viewable' => $viewable,
                ':user' => $this->uid,
                ':time' => time(),
                ':ip' => $_SERVER['REMOTE_ADDR'],
                ':url' => $this->url,
                ':action' => $data[0],
                ':desc' => $data[1]
            ));
        
        }
    
    public function getUrl()
        {
        
            $this->url = (isset($_SERVER['HTTPS']) == 'on' ? 'https' : 'http').'://'.$_SERVER['HTTP_HOST'].$_SERVER['REQUEST_URI'];
            return $this;
        
        }
    
}

?>

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
session_start();
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../../index.php');
}

if(!isset($_POST['node']))
	exit('No Node was Defined');

$selectData = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :name");
$selectData->execute(array(
	':name' => $_POST['node']
));
$node = $selectData->fetch();
?>

<p>
	<label for="server_ip">Server IP</label>
    <select name="server_ip" id="server_ip" class="round default-width-input">
        <?php
		
			$ips = json_decode($node['ips'], true);
			$i = 0;
			$hasFree = false;
			foreach($ips as $ip => $internal){
			
				if($i == 0)
					$firstIP = $ip;
				
	            if($internal['ports_free'] > 0){
	            	$hasFree = true;
					echo '<option value="'.$ip.'">'.$ip.' ('.$internal['ports_free'].' Avaliable Ports)</option>';
	            }else
				  echo '<option disabled="disabled">'.$ip.' ('.$internal['ports_free'].' Avaliable Ports)</option>';
                $i++;
                										
			}
		
		?>
    </select><i class="fa fa-angle-down pull-right select-arrow select-halfbox"></i>
    <?php
    	
    	if($hasFree === false)
    		echo '<div class="error-box round"><strong>Error:</strong> This node does not have any free ports avaliable. Please select another node.</div>';
    
    ?>
</p>
<p>
	<label for="server_port">Server Port</label>
  		<?php
	    
	        $ports = json_decode($node['ports'], true);
	        
	        foreach($ports as $ip => $internal){
	        
	            if($firstIP == $ip)
	                echo '<span id="node_'.$ip.'"><select name="server_port_'.$ip.'" class="round default-width-input">';
	            else
	                echo '<span style="display:none;" id="node_'.$ip.'"><select name="server_port_'.$ip.'" class="round default-width-input">';
	            
	                foreach($internal as $port => $avaliable){
	                
                        if($avaliable == 1)
                            echo '<option value="'.$port.'">'.$port.'</option>';
                        else
                            echo '<option disabled="disabled">'.$port.'</option>';
	                    
	                }
	            
	            echo '</select><i class="fa fa-angle-down pull-right select-arrow select-halfbox"></i></span>';
	        
	        }
														
  		?>
</p>
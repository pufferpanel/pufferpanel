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
require_once('../../../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../index.php');
}

if(!isset($_POST['node']))
	exit('No Node was Defined');

$selectData = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :name");
$selectData->execute(array(
	':name' => $_POST['node']
));
$node = $selectData->fetch();
?>

<div class="form-group col-6 nopad">
	<label for="server_ip" class="control-label">Server IP</label>
	<div>
		<select name="server_ip" id="server_ip" class="form-control">
        <?php

			$ips = json_decode($node['ips'], true);
			$i = 0;
			$hasFree = false;

			if(is_array($ips)){

				foreach($ips as $ip => $internal){

					if($i == 0)
						$firstIP = $ip;

		            if($internal['ports_free'] > 0){
		            	$hasFree = true;
						$append = ($internal['ports_free'] > 1) ? "s" : null;
						echo '<option value="'.$ip.'">'.$ip.' ('.$internal['ports_free'].' Avaliable Port'.$append.')</option>';
		            }else
					  echo '<option disabled="disabled">'.$ip.' (no ports avaliable)</option>';
	                $i++;

				}

			}

		?>
    	</select>
	</div>
</div>
<div class="form-group col-6 nopad-right">
	<label for="server_ip" class="control-label">Server Port</label>
		<div>
  		<?php

	        $ports = json_decode($node['ports'], true);

	        if(!empty($ports)){

                foreach($ports as $ip => $internal){

    	            if($firstIP == $ip)
    	                echo '<span id="node_'.$ip.'"><select name="server_port_'.$ip.'" class="form-control">';
    	            else
    	                echo '<span style="display:none;" id="node_'.$ip.'"><select name="server_port_'.$ip.'" class="form-control">';

    	                foreach($internal as $port => $avaliable){

                            if($avaliable == 1)
                                echo '<option value="'.$port.'">'.$port.'</option>';
                            else
                                echo '<option disabled="disabled">'.$port.' (in use)</option>';

    	                }

    	            echo '</select></span>';

    	        }

            }

  		?>
  		</div>
</div>
<?php

	if($hasFree === false)
		echo "<script type=\"text/javascript\">$('#noPorts').show();</script>";
?>

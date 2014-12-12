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
namespace PufferPanel\Core;
use \ORM as ORM;

require_once('../../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Components\Page::redirect('../../../index.php');
}

if(!isset($_POST['location']))
	exit('No location was Defined');

$nodes = ORM::forTable('nodes')->where('location', $_POST['location'])->findMany();

?>

<div class="form-group col-6 nopad-right">
	<label for="server_name" class="control-label">Location Node</label>
	<div>
		<select name="node" id="getNode" class="form-control">
        <?php

			if($nodes) {

				foreach($nodes as &$node) {

					echo '<option value="'.$node->id.'"> '.$node->node.' '.(($node->public == 0) ? '[Private]' : null).'</option>';

				}

			}

		?>
    	</select>
	</div>
</div>
<?php
if(!$nodes) {
	echo "<script type=\"text/javascript\">$('#noNodes').show();</script>";
}
?>

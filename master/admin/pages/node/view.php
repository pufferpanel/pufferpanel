<?php
session_start();
require_once('../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../index.php');
}

if(isset($_GET['do']) && $_GET['do'] == 'redirect' && isset($_GET['node'])){

	$select = $mysql->prepare("SELECT `id` FROM `nodes` WHERE `node` = :name");
	$select->execute(array(':name' => $_GET['node']));
	
	if($select->rowCount() == 1) {
		$n = $select->fetch();
		$core->framework->page->redirect('view.php?id='.$n['id']);
	}else{
		$core->framework->page->redirect('list.php');
	}

}

if(!isset($_GET['id']))
	$core->framework->page->redirect('list.php');

/*
 * Select Node Information
 */
$selectNode = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :id");
$selectNode->execute(array(
	':id' => $_GET['id']
));

	if($selectNode->rowCount() != 1)
		$core->framework->page->redirect('list.php?error=no_node');
	else
		$node = $selectNode->fetch();

?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>PufferPanel - Viewing Node `<?php echo $node['node']; ?>`</title>
	
	<!-- Stylesheets -->
	<link href='http://fonts.googleapis.com/css?family=Droid+Sans:400,700' rel='stylesheet'>
	<link rel="stylesheet" href="../../../assets/css/style.css">
	
	<!-- Optimize for mobile devices -->
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	
	<!-- jQuery & JS files -->
	<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
	<script src="../../../assets/javascript/jquery.cookie.js"></script>
</head>
<body>
	<div id="top-bar">
		<div class="page-full-width cf">
			<ul id="nav" class="fl">
				<li><a href="../../../account.php" class="round button dark"><i class="fa fa-user"></i>&nbsp;&nbsp; <strong><?php echo $core->framework->user->getData('username'); ?></strong></a></li>
			</ul>
			<ul id="nav" class="fr">
				<li><a href="../../../servers.php" class="round button dark"><i class="fa fa-sign-out"></i></a></li>
				<li><a href="../../../logout.php" class="round button dark"><i class="fa fa-power-off"></i></a></li>
			</ul>
		</div>	
	</div>
	<div id="header-with-tabs">
		<div class="page-full-width cf">
		</div>
	</div>
	<div id="content">
		<div class="page-full-width cf">
			<?php include('../../../core/templates/admin_sidebar.php'); ?>
			<div class="side-content fr">
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">Node <?php echo $node['node']; ?> Information</h3>
					</div>
				</div>
			</div>
			<div class="side-content fr">
				<div class="half-size-column fl">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Basic Node Information</h3>
						</div>
						<div class="content-module-main">
							<form action="ajax/update/basic.php" method="post">
								<fieldset>
									<p>
										<label for="name">Node Name</label>
										<input type="text" name="name" value="<?php echo $node['node']; ?>" class="round full-width-input" />
									</p>
									<p>
										<label for="link">Node Link</label>
										<input type="text" name="link" value="<?php echo $node['node_link']; ?>" class="round full-width-input" />
									</p>
									<p>
										<label for="ip">Node IP</label>
										<input type="text" name="ip" value="<?php echo $node['node_ip']; ?>" class="round full-width-input" />
									</p>
									<div class="stripe-separator"></div>
										<input type="hidden" name="nid" value="<?php echo $_GET['id']; ?>" />
										<input type="submit" value="Update Information" class="round blue ic-right-arrow" />
								</fieldset>
							</form>
						</div>
					</div>
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Edit IP &amp; Port Allocation</h3>
						</div>
						<div class="content-module-main">
							<form action="ajax/ports/add_port.php" id="addPorts" style="display: none;" method="post">
								<p>
									<label for="add_ports" id="setTitle"></label>
									<input type="text" name="add_ports" value="" placeholder="enter a comma separated list of ports to add; enter to submit" class="round full-width-input" />
									<input type="hidden" name="add_ports_ip" value=""/>
									<input type="hidden" name="add_ports_node" value=""/>
								</p>
							</form>
							<table>
								<thead>
									<tr>
										<th style="width:30%;">IP Address</th>
										<th style="width:35%">Ports</th>
										<th style="width:35%"></th>
									</tr>
								</thead>
								<tbody>
									<?php
									
										foreach(json_decode($node['ports'], true) as $ip => $ports)
											{
									
												echo "<tr><td style=\"vertical-align:top;\">{$ip}<br /><a href=\"#/add/{$ip}/{$node['id']}\" class=\"clickToAdd\" onclick=\"return false;\">Add Port(s)</a></td>";
												$counter = 1;
												$row1 = null; $row2 = null;
												foreach($ports as $port => $avaliable)
													{
											
														if($counter & 1)
															{
															
																$row1 .= ($avaliable == 1) ? "<span><a href=\"#/delete/{$ip}/{$port}/{$node['id']}\" class=\"deletePort\" onclick=\"return false;\"><i class=\"fa fa-circle-o\"></i></a>" : "<i class=\"fa fa-dot-circle-o\"></i>";
																$row1 .= "&nbsp;&nbsp;&nbsp; {$port}<br /></span>";
																
															}else{
															
																$row2 .= ($avaliable == 1) ? "<span><a href=\"#/delete/{$ip}/{$port}/{$node['id']}\" class=\"deletePort\" onclick=\"return false;\"><i class=\"fa fa-circle-o\"></i></a>" : "<i class=\"fa fa-dot-circle-o\"></i>";
																$row2 .= "&nbsp;&nbsp;&nbsp; {$port}<br /></span>";
															
															}
														
														$counter++;
										
													}
												echo "<td style=\"vertical-align:top;\">{$row1}</td><td style=\"vertical-align:top;\">{$row2}</td></tr>";
									
											}
									
									?>
								</tbody>
							</table>
							<div class="information-box round no-image"><i class="fa fa-circle-o"></i> (Port Available; Click to Delete Port)<br /><i class="fa fa-dot-circle-o"></i> (Port Used; Cannot Delete)</div>
							<div class="error-box round">Editing this information is currently unavailable.</div>
						</div>
					</div>
				</div>
				<div class="half-size-column fr">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Update SFTP IP &amp; Username</h3>
						</div>
						<div class="content-module-main">
							<p>If you have changed your SFTP IP address or the username of the main account used for provisioning servers please update it below. If the password has changed as well you can change that in the box below.</p>
							<form action="ajax/update/sftp.php?do=ip+user" method="post">
								<fieldset>
									<p>
										<label for="sftp_ip">SFTP IP Address</label>
										<input type="text" name="sftp_ip" value="<?php echo $node['sftp_ip']; ?>" class="round full-width-input" />
									</p>
									<p>
										<label for="sftp_user">SFTP Username</label>
										<input type="text" name="sftp_user" value="<?php echo $node['username']; ?>" class="round full-width-input" />
									</p>
									<div class="stripe-separator"></div>
										<div class="warning-box no-image round">
											Editing this node information will require us to update multiple records in the database for servers in order to reflect these changes. Please ensure that you have entered the above information correctly. Changing this wrongly could result in multiple clients being unable to access their server(s).<br /><br />
											<input type="checkbox" name="warning" /> I have read and understand the above statement.
										</div>
										<input type="hidden" name="nid" value="<?php echo $_GET['id']; ?>" />
										<input type="submit" value="Update SFTP" class="round blue ic-right-arrow" />
								</fieldset>
							</form>
						</div>
					</div>
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Update SFTP Password</h3>
						</div>
						<div class="content-module-main">
							<form action="ajax/update/sftp.php?do=pass" method="post">
								<fieldset>
									<p>
										<label for="pass">New SFTP Password</label>
										<input type="password" name="pass" value="" class="round full-width-input" />
									</p>
									<p>
										<label for="pass_2">New SFTP Password Again</label>
										<input type="password" name="pass_2" value="" class="round full-width-input" />
									</p>
									<div class="stripe-separator"></div>
										<div class="warning-box no-image round">
											Please ensure that you have entered the above information correctly. Changing this wrongly could result in multiple clients being unable to access their server(s).<br /><br />
											<input type="checkbox" name="warning" /> I have read and understand the above statement.
										</div>
										<input type="hidden" name="nid" value="<?php echo $_GET['id']; ?>" />
										<input type="submit" value="Update SFTP Password" class="round blue ic-right-arrow" />
								</fieldset>
							</form>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function() {
			$(".clickToAdd").click(function(){
				var rawUrl = $(this).attr("href");
				var exploded = rawUrl.split('/');
				var ip = exploded[2];
				var node_id = exploded[3];
				$("#addPorts").slideUp(function(){
					$("#setTitle").html('Add Ports for '+ip);
					$("input[name='add_ports']").val("");
					$("input[name='add_ports_ip']").val(ip);
					$("input[name='add_ports_node']").val(node_id);
					$("#addPorts").slideDown();
				});
			});
			$(".deletePort").click(function(){
				
				var rawUrl = $(this).attr("href");
				var exploded = rawUrl.split('/');
				var ip = exploded[2];
				var port = exploded[3];
				var node_id = exploded[4];
				var conf = confirm("Are you sure you want to delete "+ip+":"+port);
				
					if(conf == true)
						{
							$.ajax({
								type: "POST",
								url: "ajax/ports/delete.php",
								data: { ip: ip, port: port, node: node_id},
								success: function(data) {
									$(".deletePort[href='#/delete/"+ip+"/"+port+"/"+node_id+"']").parent().fadeOut();
								}
							});
						}else{
							return false;
						}
				
			});
			// Advanced Error Stuff
			$.urlParam = function(name){
			    var results = new RegExp('[\\?&]' + name + '=([^&#]*)').exec(decodeURIComponent(window.location.href));
			    if (results==null){
			       return null;
			    }
			    else{
			       return results[1] || 0;
			    }
			}
			if($.urlParam('error') != null){
			
				var field = $.urlParam('error');
				var exploded = field.split('|');
				
					$.each(exploded, function(key, value) {
						
						$('[name="' + value + '"]').addClass('error-input');
						
					});
					
				var obj = $.parseJSON($.cookie('__TMP_pp_admin_newnode'));
				
					$.each(obj, function(key, value) {
						
						$('[name="' + key + '"]').val(value);
						
					});
			
			}
		});
	</script>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>
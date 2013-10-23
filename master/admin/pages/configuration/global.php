<?php
session_start();
require_once('../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../index.php');
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>PufferPanel - Admin Global Settings</title>
	
	<!-- Stylesheets -->
	<link href='http://fonts.googleapis.com/css?family=Droid+Sans:400,700' rel='stylesheet'>
	<link rel="stylesheet" href="../../../assets/css/style.css">
	
	<!-- Optimize for mobile devices -->
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	
	<!-- jQuery & JS files -->
	<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
	<script src="//ajax.googleapis.com/ajax/libs/jqueryui/1.10.3/jquery-ui.min.js"></script>
	
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
				<div class="half-size-column fl">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Company Name</h3>
						</div>
						<div class="content-module-main cf">
							<form action="account.php?action=email" method="post">
								<fieldset>
									<p>
										<label for="newemail">Company Name</label>
										<input type="text" name="company_name" class="round full-width-input" value="<?php echo $core->framework->settings->get('company_name'); ?>" />
									</p>
									<div class="stripe-separator"><!--  --></div>
									<input type="submit" value="Update Company Name" class="round blue ic-right-arrow" />
								</fieldset>
							</form>	
						</div>
					</div>
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">URL Settings</h3>
						</div>
						<div class="content-module-main">
							<form action="account.php?action=password" method="post">
								<fieldset>
									<p>
										<label for="p_password">Main Website URL</label>
										<input type="text" name="main_url" class="round full-width-input" value="<?php echo $core->framework->settings->get('main_website'); ?>"/>
										<em>The URL corresponding to your main website.</em>
									</p>
									<p>
										<label for="p_password_new">PufferPanel Master URL</label>
										<input type="text" name="master_url" class="round full-width-input" value="<?php echo $core->framework->settings->get('master_url'); ?>"/>
										<em>The URL corresponding to this PufferPanel installation.</em>
									</p>
									<p>
										<label for="p_password_new_2">PufferPanel Assets URL</label>
										<input type="text" name="assets_url" class="round full-width-input" value="<?php echo $core->framework->settings->get('assets_url'); ?>"/>
										<em>The URL corresponding to the assets for PufferPanel. Update this only if you are using a CDN or Caching Service that modifies where the assets are stored.</em>
									</p>
									<div class="stripe-separator"><!--  --></div>
									<div class="warning-box round">Trailing slashes are <strong>required</strong>.</div>
									<input type="submit" value="Update Information" class="round blue ic-right-arrow" />
								</fieldset>
							</form>							
						</div>
					</div>
				</div>
				<div class="half-size-column fr">
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">Email Settings</h3>
						</div>
						<div class="content-module-main cf">
							<form action="account.php?action=email" method="post">
								<fieldset>
									<?php
										/*
										 * Select Current Email Method
										 */
										$method = $core->framework->settings->get('sendmail_method');
										$marray = array('php' => '', 'postmark' => '', 'mandrill' => '');
										
											foreach($marray as $id => $value){
											
												if($method == $id){
													$marray[$id] = 'selected="selected"';
												}
											
											}
										
									?>
									<p>
										<label for="p_password">Sendmail Method</label>
										<select name="smail_method" class="round" id="smail_method">
											<option value="php" <?php echo $marray['php']; ?>>PHP mail()</option>
											<option value="postmark" <?php echo $marray['postmark']; ?>>Postmark</option>
											<option value="mandrill" <?php echo $marray['mandrill']; ?>>Mandrill</option>
										</select>
										<i class="fa fa-angle-down pull-right select-arrow select-halfbox"></i>
									</p>
									<p>
										<label for="p_password">From Address</label>
										<input type="text" name="main_url" class="round full-width-input" value="<?php echo $core->framework->settings->get('sendmail_email'); ?>"/>
										<em>The email address all outgoing emails should use. If using Postmark or Mandrill this must match the email used in their settings.</em>
									</p>
									<p id="postmark">
										<label for="p_password">Postmark API Key</label>
										<input type="text" name="main_url" class="round full-width-input" value="<?php echo $core->framework->settings->get('postmark_api_key'); ?>"/>
										<em>The API key generated on <a href="https://postmarkapp.com/">Postmark</a>. Leave blank if not using.</em>
									</p>
									<p id="mandrill">
										<label for="p_password_new">Mandrill API Key</label>
										<input type="text" name="master_url" class="round full-width-input" value="<?php echo $core->framework->settings->get('mandrill_api_key'); ?>"/>
										<em>The API key generated on <a href="https://mandrill.com/">Mandrill</a>. Leave blank if not using.</em>
									</p>
									<div class="stripe-separator"><!--  --></div>
									<input type="submit" value="Update Email" class="round blue ic-right-arrow" />
								</fieldset>
							</form>	
						</div>
					</div>
					<div class="content-module">
						<div class="content-module-heading cf">
							<h3 class="fl">reCAPTCHA Settings</h3>
						</div>
						<div class="content-module-main cf">
							<form action="account.php?action=notifications" method="post">
								<fieldset>
									<p>
										<label for="p_password">Public Key</label>
										<input type="text" name="main_url" class="round full-width-input" value="<?php echo $core->framework->settings->get('captcha_pub'); ?>"/>
									</p>
									<p>
										<label for="p_password_new">Private Key</label>
										<input type="text" name="master_url" class="round full-width-input" value="<?php echo $core->framework->settings->get('captcha_priv'); ?>"/>
									</p>
									<div class="stripe-separator"><!--  --></div>
									<input type="submit" value="Update reCAPTCHA" class="round blue ic-right-arrow" />
									<p><em>reCAPTCHA is the system used to help prevent people from abusing the password reset page on PufferPanel. You may use these default provided keys, or you may generate your own at <a href="https://www.google.com/recaptcha">Google reCAPTCHA</a>.</em></p>
								</fieldset>
							</form>	
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			
			var method = $("#smail_method :selected").val();
				if(method == "postmark"){
					$("#mandrill").hide();
				}else if(method == "mandrill"){
					$("#postmark").hide();
				}else{
					$("#mandrill").hide();
					$("#postmark").hide();
				}
		
			$("#smail_method").change(function(){
				var method = $("#smail_method :selected").val();
					if(method == "postmark"){
						if($("#mandrill").is(':visible')){
							$("#mandrill").toggle("drop", "right", function(){
								$("#postmark").toggle("drop", "right");
							});
						}else{
							$("#postmark").toggle("drop", "right");
						}
					}else if(method == "mandrill"){
						if($("#postmark").is(':visible')){
							$("#postmark").toggle("drop", "right", function(){
								$("#mandrill").toggle("drop", "right");
							});
						}else{
							$("#mandrill").toggle("drop", "right");
						}
					}else{
						if($("#postmark").is(':visible')){
							$("#postmark").toggle("drop", "up");
						}
						if($("#mandrill").is(':visible')){
							$("#mandrill").toggle("drop", "down");
						}
					}
			});
		});
	</script>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>
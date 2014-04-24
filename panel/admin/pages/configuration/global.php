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
require_once('../../../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../index.php');
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../../../assets/include/header.php'); ?>
	<title>PufferPanel Admin Control Panel</title>
</head>
<body>
	<div class="container">
		<div class="navbar navbar-default">
			<div class="navbar-header">
				<a class="navbar-brand" href="#"><?php echo $core->settings->get('company_name'); ?></a>
			</div>
			<div class="navbar-collapse navbar-responsive-collapse collapse" style="height: 1px;">
				<ul class="nav navbar-nav navbar-right">
					<li class="dropdown">
						<a href="#" class="dropdown-toggle" data-toggle="dropdown">Account <b class="caret"></b></a>
							<ul class="dropdown-menu">
								<li><a href="<?php echo $core->settings->get('master_url'); ?>logout.php">Logout</a></li>
								<li><a href="<?php echo $core->settings->get('master_url'); ?>servers.php">View All Servers</a></li>
							</ul>
					</li>
				</ul>
			</div>
		</div>
		<div class="row">
			<div class="col-3"><?php include('../../../assets/include/admin.php'); ?></div>
			<div class="col-9">
				<ul class="nav nav-tabs" id="config_tabs">
					<li class="active"><a href="#company" data-toggle="tab">Company Name</a></li>
					<li><a href="#url" data-toggle="tab">URL Settings</a></li>
					<li><a href="#cookies" data-toggle="tab">Cookies</a></li>
					<li><a href="#email" data-toggle="tab">Email Settings</a></li>
					<li><a href="#captcha" data-toggle="tab">reCAPTCHA</a></li>
				</ul>
				<div class="tab-content">
					<div class="tab-pane active" id="company">
						<h3>Company Name</h3><hr />
						<form action="actions/cname.php" method="POST">
							<fieldset>
								<div class="input-group">
									<input type="text" name="company_name" class="form-control" value="<?php echo $core->settings->get('company_name'); ?>" />
									<span class="input-group-btn">
										<input type="submit" value="Update Company Name" class="btn btn-primary" />
									</span>
								</div>
							</fieldset>
						</form>
					</div>
					<div class="tab-pane" id="url">
						<h3>URL Settings</h3><span class="label label-warning">Trailing Slashes Required</span><hr />
						<form action="actions/url.php" method="POST">
							<fieldset>
								<div class="form-group">
									<label for="main_url" class="control-label">Main Website URL</label>
									<div>
										<input type="text" name="main_url" class="form-control" value="<?php echo $core->settings->get('main_website'); ?>"/>
										<p><small class="text-muted"><em>The URL corresponding to your main website, not necessarily this PufferPanel installation.</em></small></p>
									</div>
								</div>
								<div class="form-group">
									<label for="master_url" class="control-label">PufferPanel Master URL</label>
									<div>
										<input type="text" name="master_url" class="form-control" value="<?php echo $core->settings->get('master_url'); ?>"/>
										<p><small class="text-muted"><em>The URL corresponding to this PufferPanel installation.</em></small></p>
									</div>
								</div>
								<div class="form-group">
									<label for="assets_url" class="control-label">Assets Directory URL</label>
									<div>
										<input type="text" name="assets_url" class="form-control" value="<?php echo $core->settings->get('assets_url'); ?>"/>
										<p><small class="text-muted"><em>The URL corresponding to the assets for PufferPanel. Update this only if you are using a CDN or Caching Service that modifies where the assets are stored.</em></small></p>
									</div>
								</div>
								<div class="form-group">
									<div>
										<input type="submit" value="Update URL Settings" class="btn btn-primary" />
									</div>
								</div>
							</fieldset>
						</form>
					</div>
					<div class="tab-pane" id="cookies">
						<h3>Cookie Settings</h3><hr />
						<form action="actions/cookies.php" method="POST">
							<fieldset>
								<div class="form-group">
									<label for="cookie_website" class="control-label">Cookie Website</label>
									<div>
										<input type="text" name="cookie_website" class="form-control" value="<?php echo $core->settings->get('cookie_website'); ?>"/>
										<p><small class="text-muted"><em>This should be the website that PufferPanel is running on (but can be left blank). <strong>Setting this to the wrong value will lock you and all other users out of the panel.</strong></em></small></p>
									</div>
								</div>
								<div class="form-group">
									<div>
										<input type="submit" value="Update Cookie Settings" class="btn btn-primary" />
									</div>
								</div>
							</fieldset>
						</form>
					</div>
					<div class="tab-pane" id="email">
						<h3>Email Settings</h3><hr />
						<form action="actions/email.php" method="POST">
							<fieldset>
								<?php
									/*
									 * Select Current Email Method
									 */
									$method = $core->settings->get('sendmail_method');
									$marray = array('php' => '', 'postmark' => '', 'mandrill' => '', 'mailgun' => '');
									
										foreach($marray as $id => $value){
										
											if($method == $id){
												$marray[$id] = 'selected="selected"';
											}
										
										}
									
								?>
								<div class="form-group col-6 nopad">
									<label for="smail_method" class="control-label">Sendmail Method</label>
									<div>
										<select name="smail_method" class="form-control" id="smail_method">
											<option value="php" <?php echo $marray['php']; ?>>PHP mail()</option>
											<option value="postmark" <?php echo $marray['postmark']; ?>>Postmark</option>
											<option value="mandrill" <?php echo $marray['mandrill']; ?>>Mandrill</option>
											<option value="mailgun" <?php echo $marray['mailgun']; ?>>Mailgun</option>
										</select>
									</div>
								</div>
								<div class="form-group col-6 nopad-right">
									<label for="sendmail_email" class="control-label">Sendmail Address</label>
									<div>
										<input type="text" name="sendmail_email" class="form-control" value="<?php echo $core->settings->get('sendmail_email'); ?>"/>
										<p><small class="text-muted"><em>The email address all outgoing emails should use. If using Postmark, Mandrill, or Mailgun this must match the email used in their settings.</em></small></p>
									</div>
								</div>
								<div class="form-group" id="postmark">
									<label for="postmark_api_key" class="control-label">PostmarkApp API Key</label>
									<div>
										<input type="text" name="postmark_api_key" class="form-control" value="<?php echo $core->settings->get('postmark_api_key'); ?>"/>
										<p><small class="text-muted"><em>The API key generated on <a href="https://postmarkapp.com/">Postmark</a>. Leave blank if not using.</em></small></p>
									</div>
								</div>
								<div class="form-group" id="mandrill">
									<label for="mandrill_api_key" class="control-label">Mandrill API Key</label>
									<div>
										<input type="text" name="mandrill_api_key" class="form-control" value="<?php echo $core->settings->get('mandrill_api_key'); ?>"/>
										<p><small class="text-muted"><em>The API key generated on <a href="https://mandrill.com/">Mandrill</a>. Leave blank if not using.</em></small></p>
									</div>
								</div>
								<div class="form-group" id="mailgun">
									<label for="mailgun_api_key" class="control-label">Mailgun API Key</label>
									<div>
										<input type="text" name="mailgun_api_key" class="form-control" value="<?php echo $core->settings->get('mailgun_api_key'); ?>"/>
										<p><small class="text-muted"><em>The API key generated on <a href="https://mailgun.com/">Mailgun</a>. Leave blank if not using.</em></small></p>
									</div>
								</div>
								<div class="form-group">
									<div>
										<input type="submit" value="Update Email Settings" class="btn btn-primary" />
									</div>
								</div>
							</fieldset>
						</form>
					</div>
					<div class="tab-pane" id="captcha">
						<h3>reCAPTCHA Settings</h3><hr />
						<form action="actions/captcha.php" method="POST">
							<fieldset>
								<div class="form-group">
									<label for="pub_key" class="control-label">Public Key</label>
									<div>
										<input type="text" name="pub_key" class="form-control" value="<?php echo $core->settings->get('captcha_pub'); ?>"/>
									</div>
								</div>
								<div class="form-group">
									<label for="priv_key" class="control-label">Private Key</label>
									<div>
										<input type="text" name="priv_key" class="form-control" value="<?php echo $core->settings->get('captcha_priv'); ?>"/>
									</div>
								</div>
								<div class="form-group">
									<div>
										<input type="submit" value="Update reCAPTCHA Settings" class="btn btn-primary" />
										<p><small class="text-muted"><em>reCAPTCHA is the system used to help prevent people from abusing the password reset page on PufferPanel. You may use these default provided keys, or you may generate your own at <a href="https://www.google.com/recaptcha">Google reCAPTCHA</a>.</em></small></p>
									</div>
								</div>
							</fieldset>
						</form>
					</div>
				</div>
			</div>
		</div>
		<div class="footer">
			<?php include('../../../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			var method = $("#smail_method :selected").val();
			if(method == "postmark"){
				$("#mandrill").hide();
				$("#mailgun").hide();
			}else if(method == "mandrill"){
				$("#postmark").hide();
				$("#mailgun").hide();
			}else if(method == "mailgun"){
				$("#mandrill").hide();
				$("#postmark").hide();
			}else{
				$("#mandrill").hide();
				$("#postmark").hide();
				$("#mailgun").hide();
			}
			$("#smail_method").change(function(){
					var method = $("#smail_method :selected").val();
					if(method == "postmark"){
						if($("#postmark").not(':visible')){
							$("#mandrill").hide();
							$("#mailgun").hide();
							$("#postmark").show();
						}
					}else if(method == "mandrill"){
						if($("#mandrill").not(':visible')){
							$("#postmark").hide();
							$("#mailgun").hide();
							$("#mandrill").show();
						}
					}else if(method == "mailgun"){
						if($("#mailgun").not(':visible')){
							$("#postmark").hide();
							$("#mandrill").hide();
							$("#mailgun").show();
						}
					}else{
						$("#mandrill").hide();
						$("#mailgun").hide();
						$("#postmark").hide();
					}
			});
			if($.urlParam('tab') != null){
				$('#config_tabs a[href="#'+$.urlParam('tab')+'"]').tab('show');
			}
			if($.urlParam('error') != null){
				var field = $.urlParam('error');
				var exploded = field.split('|');
					$.each(exploded, function(key, value) {
						$('[name="' + value + '"]').parent().parent().addClass('has-error');
					});
				var obj = $.parseJSON($.cookie('__TMP_pp_admin_updateglobal'));
					$.each(obj, function(key, value) {
						$('[name="' + key + '"]').val(value);
					});			
			}
		});
	</script>
</body>
</html>

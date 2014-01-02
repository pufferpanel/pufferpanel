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

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../index.php');
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../../../assets/include/header.php'); ?>
	<title>PufferPanel - Admin Global Settings</title>
</head>
<body>
	<div class='container'>
		<?php include('../../../core/templates/admin_top.php'); ?>
		<div class='row'>
			<div class="col-3">
				<?php include('../../../core/templates/admin_sidebar.php'); ?>
			</div>
			<div class="col-9">
				<div class="row">
					<div class="col-6">
						<h3 style="margin-top:0;">URL Settings</h3><hr />
						<form action="actions/url.php" method="post">
							<div class="form-group">
								<label for="main_url" class="control-label">Main Website URL</label>
								<div>
									<input type="text" name="main_url" class="form-control" value="<?php echo $core->framework->settings->get('main_website'); ?>"/>
								</div>
								<span class='text-muted'><small>The URL corresponding to your main website.</small></span>
							</div>
							<div class="form-group">
								<label for="master_url" class="control-label">PufferPanel Master URL</label>
								<div>
									<input type="text" name="master_url" class="form-control" value="<?php echo $core->framework->settings->get('master_url'); ?>"/>
								</div>
								<span class='text-muted'><small>The URL corresponding to this PufferPanel installation.</small></span>
							</div>
							<div class="form-group">
								<label for="assets_url" class="control-label">PufferPanel Assets URL</label>
								<div>
									<input type="text" name="assets_url" class="form-control" value="<?php echo $core->framework->settings->get('assets_url'); ?>"/>
								</div>
								<span class='text-muted'><small>The URL corresponding to the assets for PufferPanel. Update this only if you are using a CDN or Caching Service that modifies where the assets are stored.</small></span>
							</div>
							<div class="form-group">
								<div>
									<input type="submit" value="Update Information" class="btn btn-primary" />
								</div>
							</div>
						</form>
					</div>
					<div class="col-6">
						<h3 style="margin-top:0;">Email Settings</h3><hr />
						<form action="actions/email.php" method="post">
							<?php
								/*
								 * Select Current Email Method
								 */
								$method = $core->framework->settings->get('sendmail_method');
								$marray = array('php' => '', 'postmark' => '', 'mandrill' => '', 'mailgun' => '');
								
									foreach($marray as $id => $value){
									
										if($method == $id){
											$marray[$id] = 'selected="selected"';
										}
									
									}
								
							?>
							<div class="form-group">
								<label for="smail_method" class="control-label">Sendmail Method</label>
								<select name="smail_method" class="form-control" id="smail_method">
									<option value="php" <?php echo $marray['php']; ?>>PHP mail()</option>
									<option value="postmark" <?php echo $marray['postmark']; ?>>Postmark</option>
									<option value="mandrill" <?php echo $marray['mandrill']; ?>>Mandrill</option>
									<option value="mailgun" <?php echo $marray['mailgun']; ?>>Mailgun</option>
								</select>
							</div>
							<div class="form-group">
								<label for="sendmail_email" class="control-label">From Address</label>
								<input type="text" name="sendmail_email" class="form-control" value="<?php echo $core->framework->settings->get('sendmail_email'); ?>"/>
								<span class='text-muted'><small>The email address all outgoing emails should use. If using Postmark, Mandrill, or Mailgun this must match the email used in their settings.</small></span>
							</div>
							<div class="form-group">
								<label for="postmark_api_key" class="control-label">Postmark API Key</label>
								<input type="text" name="postmark_api_key" class="form-control" value="<?php echo $core->framework->settings->get('postmark_api_key'); ?>"/>
								<span class='text-muted'><small><a href="https://postmarkapp.com/">Postmark</a>. Leave blank if not using.</small></span>
							</div>
							<div class="form-group">
								<label for="mandrill_api_key" class="control-label">Mandrill API Key</label>
								<input type="text" name="mandrill_api_key" class="form-control" value="<?php echo $core->framework->settings->get('mandrill_api_key'); ?>"/>
								<span class='text-muted'><small>The API key generated on <a href="https://mandrill.com/">Mandrill</a>. Leave blank if not using.</small></span>
							</div>
							<div class="form-group">
								<label for="mailgun_api_key" class="control-label">Mailgun API Key</label>
								<input type="text" name="mailgun_api_key" class="form-control" value="<?php echo $core->framework->settings->get('mailgun_api_key'); ?>"/>
								<span class='text-muted'><small>The API key generated on <a href="https://mailgun.com/">Mailgun</a>. Leave blank if not using.</small></span>
							</div>
							<div class="form-group">
								<input type="submit" value="Update Email Settings" class="btn btn-primary" />
							</div>
						</form>
					</div>
				</div>
				<div class='row'>
					<div class="col-6">
						<h3 style="margin-top:0;">Company Name</h3><hr />
						<form action="actions/cname.php" method="POST">
							<div class="form-group">
								<label for="company_name" class="control-label">Company Name</label>
								<div>
									<input type="text" name="company_name" class="form-control" value="<?php echo $core->framework->settings->get('company_name'); ?>" />
								</div>
							</div>
							<div class="form-group">
								<div>
									<input type="submit" class="btn btn-primary" value="Update Company Name" />
								</div>
							</div>
						</form>
					</div>
					<div class="col-6">
						<h3 style="margin-top:0;">reCAPTCHA Settings</h3><hr />
						<form action="actions/captcha.php" method="post">
									<div class="form-group">
										<label for="pub_key" class="control-label">Public Key</label>
										<input type="text" name="pub_key" class="form-control" value="<?php echo $core->framework->settings->get('captcha_pub'); ?>"/>
									</div>
									<div class="form-group">
										<label for="priv_key" class="control-label">Private Key</label>
										<input type="text" name="priv_key" class="form-control" value="<?php echo $core->framework->settings->get('captcha_priv'); ?>"/>
									</div>
									<div class="form-group">
										<input type="submit" value="Update reCAPTCHA" class="btn btn-primary" />
									</div>
									<span class='text-muted'><small>reCAPTCHA is the system used to help prevent people from abusing the password reset page on PufferPanel. You may use these default provided keys, or you may generate your own at <a href="https://www.google.com/recaptcha">Google reCAPTCHA</a>.</small></span>
						</form>
					</div>
					
				</div>
			</div>
		</div>
		<script type='text/javascript'>
			$( document ).ready(function() {
				$( "#admin-21" ).addClass( "active" );
			});
		</script>
		<div class='footer'>
			<?php include('../../../assets/include/footer.php'); ?>
		</div>
	</div>
</body>
</html>
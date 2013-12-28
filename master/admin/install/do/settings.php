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
if(file_exists('../install.lock'))
	exit('Installer is Locked.');
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>PufferPanel - Install</title>
	
	<!-- Stylesheets -->
	<link href='http://fonts.googleapis.com/css?family=Droid+Sans:400,700' rel='stylesheet'>
	<link rel="stylesheet" href="../../../assets/css/style.css">
	
	<!-- Optimize for mobile devices -->
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	
	<!-- jQuery & JS files -->
	<script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>
</head>
<body>
	<div id="top-bar">
		<div class="page-full-width cf">
            &nbsp;
		</div>	
	</div>
	<div id="header-with-tabs">
		<div class="page-full-width cf">
		</div>
	</div>
	<div id="content">
		<div class="page-full-width cf">
            <div class="content-module">
				<div class="content-module-main">
				    <h1>General Settings</h1>
                    <p>This information can be changed later on. Please provide accurate information for URLs, using the wrong link can break the system.</p>
                    <div class="stripe-separator"><!--  --></div>
                    <?php
                    
                        if(isset($_POST['do_settings'])){
                        
                            include('../../../core/framework/framework.database.connect.php');
                            $mysql = dbConn::getConnection();
                            
                            $prepare = $mysql->prepare("INSERT INTO `acp_settings` (`setting_ref`, `setting_val`) VALUES
                                ('company_name', :cname),
                                ('master_url', :murl),
                                ('cookie_website', :cwebsite),
                                ('postmark_api_key', NULL),
                                ('mandrill_api_key', NULL),
                                ('mailgun_api_key', NULL),
                                ('sendmail_email', :smail),
                                ('main_website', :mwebsite),
                                ('sendmail_method','php'),
                                ('captcha_pub','6LdSzuYSAAAAAHkmq8LlvmhM-ybTfV8PaTgyBDII'),
                                ('captcha_priv','6LdSzuYSAAAAAISSAYIJrFGGGJHi5a_V3hGRvIAz'),
                                ('assets_url', :aurl),
                                ('use_api','0'),
                                ('api_key', NULL),
                                ('api_allowed_ips','*'),
                                ('api_module_controls_all','0')");
                            
                            $cookieSite = (strtolower($_POST['cookie_website']) == 'null' || empty($_POST['cookie_website'])) ? null : $_POST['cookie_website'];
                            $prepare->execute(array(
                                ':cname' => $_POST['company_name'],
                                ':murl' => $_POST['master_url'],
                                ':cwebsite' => $cookieSite,
                                ':smail' => $_POST['sendmail_email'],
                                ':mwebsite' => $_POST['main_website'],
                                ':aurl' => $_POST['assets_url']
                            ));
                            
                            exit('<meta http-equiv="refresh" content="0;url=hash.php"/>');
                            
                        }
                    
                    ?>
                    <form action="settings.php" method="post">
                        <p>
                            <label for="company_name">Company Name</label>
                            <input type="text" name="company_name" class="round default-width-input" />
                        </p>
                        <p>
                            <label for="main_website">Main Website URL</label>
                            <input type="text" name="main_website" placeholder="http://example.com/" class="round default-width-input" />
                        </p>
                        <p>
                            <label for="master_url">PufferPanel Master URL</label>
                            <input type="text" name="master_url" placeholder="http://example.com/pufferpanel/" class="round default-width-input" />
                            <em>Trailing slashes are required.</em>
                        </p>
                        <p>
                            <label for="assets_url">PufferPanel Assets URL</label>
                            <input type="text" name="assets_url" placeholder="http://example.com/pufferpanel/assets/" class="round default-width-input" />
                            <em>Trailing slashes are required.</em>
                        </p>
                        <p>
                            <label for="cookie_website">Cookie Website</label>
                            <input type="text" name="cookie_website" placeholder="example.com" class="round default-width-input" />
                        </p>
                        <p>
                            <label for="sendmail_email">Sendmail Email</label>
                            <input type="text" name="sendmail_email" class="round default-width-input" />
                        </p>
                        <input type="submit" name="do_settings" value="Setup Database" class="round blue ic-right-arrow" />
                    </form>
				</div>
            </div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.4 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>
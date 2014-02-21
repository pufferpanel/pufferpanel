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
				    <h1>Create Your Account</h1>
                    <p>You've reached the final step of the process. Please fill out the information below to create your admin account. After finishing this step you will be redirected to the login page where you will be able to access the Admin Control Panel to add nodes, users, and servers. Thank you for installing PufferPanel on your server. Please contact us on IRC (irc.esper.net) in #pufferpanel if you encounter any problems or have questions, comments, or concerns.</p>
                    <div class="stripe-separator"><!--  --></div>
                    <?php
                    
                        if(isset($_POST['do_account'])){
                        
                            include('../../../core/framework/framework.database.connect.php');
                            $mysql = dbConn::getConnection();
                            
                            $prepare = $mysql->prepare("INSERT INTO `users` VALUES(NULL, NULL, :username, :email, :password, :time, 'owner', NULL, NULL, NULL, 1, 0, 1)");
                            
                            include('../../../core/framework/configuration.php');
                            $salt = crypt($_POST['password'], '$6$rounds=5000$'.$_INFO['salt'].'$');
                            $password = hash('ripemd320', $salt);
                            
                            $prepare->execute(array(
                                ':username' => $_POST['username'],
                                ':email' => $_POST['email'],
                                ':password' => $password,
                                ':time' => time()
                            ));
                            
                            rename('../install.lock.dist', '../install.lock');
                            
                            exit('<meta http-equiv="refresh" content="0;url=../../../index.php"/>');
                            
                        }
                    
                    ?>
                    <form action="account.php" method="post">
                        <p>
                            <label for="email">Email</label>
                            <input type="text" name="email" class="round default-width-input" />
                        </p>
                        <p>
                            <label for="username">Username</label>
                            <input type="text" name="username" class="round default-width-input" />
                        </p>
                        <p>
                            <label for="password">Password</label>
                            <input type="password" name="password" class="round default-width-input" />
                        </p>
                        <input type="submit" name="do_account" value="Setup Account" class="round blue ic-right-arrow" />
                    </form>
				</div>
            </div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.4.2 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>
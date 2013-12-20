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
				    <h1>Database Information</h1>
                    <p>Please fill out the database connection information that you will be using. <strong>For security reasons ensure it has a strong password and <span style="color:red;">do not</span> run this database under root credentials.</strong> This database user only needs the following permissions: CREATE, INSERT, SELECT, UPDATE, DELETE. This database user must be accessable from other servers.</p>
                    <div class="stripe-separator"><!--  --></div>
                    <?php
                    
                        /* Check Connection Information */
                        if(isset($_POST['do_connect'])){
                            
                            try {
                                
                                $database = new PDO('mysql:host='.$_POST['sql_h'].';dbname='.$_POST['sql_db'], $_POST['sql_u'], $_POST['sql_p'], array(
                                    PDO::ATTR_PERSISTENT => true,
                                    PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
                                ));
                        
                                $database->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
                                
                                /* Make File */
                                $keyset  = "abcdefghijklmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!^*%#@";
                                $randkey = "";

                                for ($i=0; $i<30; $i++){
                                    $randkey .= substr($keyset, rand(0, strlen($keyset)-1), 1);
                                }
                                
                                    $fp = fopen('../../../core/framework/master_configuration.php.dist', 'w+');
                                    fwrite($fp, "<?php
\$_INFO['sql_u'] = '".$_POST['sql_u']."';
\$_INFO['sql_p'] = '".$_POST['sql_p']."';
\$_INFO['sql_h'] = '".$_POST['sql_h']."';
\$_INFO['sql_db'] = '".$_POST['sql_db']."';
\$_INFO['salt'] = '".$randkey."';");
                                    fclose($fp);
                                
                                    rename('../../../core/framework/master_configuration.php.dist', '../../../core/framework/master_configuration.php');
                                    exit('<meta http-equiv="refresh" content="0;url=tables.php"/>');
                        
                            }catch (PDOException $e) {
                        
                                echo '<div class="error-box round">MySQL Connection Error: ' . $e->getMessage() . '</div>';
                        
                            }
                            
                        }
                    
                    ?>
                    <form action="start.php" method="post">
                        <p>
                            <label for="sql_db">Database Name</label>
                            <input type="text" name="sql_db" value="pufferpanel" class="round default-width-input" />
                        </p>
                        <p>
                            <label for="sql_h">Database Host</label>
                            <input type="text" name="sql_h" value="<?php echo $_SERVER['SERVER_ADDR']; ?>" class="round default-width-input" />
                        </p>
                        <p>
                            <label for="sql_u">Database User</label>
                            <input type="text" name="sql_u" class="round default-width-input" />
                        </p>
                        <p>
                            <label for="sql_p">Database User Password</label>
                            <input type="password" name="sql_p" class="round default-width-input" />
                        </p>
                        <input type="submit" name="do_connect" value="Setup Database" class="round blue ic-right-arrow" />
                    </form>
				</div>
            </div>
		</div>
	</div>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>
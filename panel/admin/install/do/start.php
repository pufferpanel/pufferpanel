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
	<link rel="stylesheet" href="../../../assets/css/bootstrap.css">
	<title>PufferPanel Installer</title>
</head>
<body>
	<div class="container">
		<div class="alert alert-danger">
			<strong>WARNING:</strong> Do not run this version on a live environment! There are known security holes that we are working on getting patched. This is extremely beta software and this version is to get the features in place while we work on security enhancements.
		</div>
		<div class="navbar navbar-default">
			<div class="navbar-header">
				<a class="navbar-brand" href="#">Install PufferPanel - Database Information</a>
			</div>
		</div>
		<div class="col-12">
			<div class="row">
				<div class="col-2"></div>
				<div class="col-8">
					<div class="alert alert-info">For security reasons ensure it has a strong password and <strong>do not</strong> run this database under root credentials.</div>
					<p>Please fill out the database connection information that you will be using. This database user only needs the following permissions: <code>CREATE</code>, <code>INSERT</code>, <code>SELECT</code>, <code>UPDATE</code>, <code>DELETE</code>. This database user must be accessible from other servers.</p>
					<?php
                    
                        /* Check Connection Information */
                        if(isset($_POST['do_connect'])){
                            
                            try {
                                
                                $database = new PDO('mysql:host='.$_POST['sql_h'].';dbname='.$_POST['sql_db'], $_POST['sql_u'], $_POST['sql_p'], array(
                                    PDO::ATTR_PERSISTENT => true,
                                    PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
                                ));
                        
                                $database->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
                                                                
                                    $fp = fopen('../../../core/framework/configuration.php.dist', 'w+');
                                    fwrite($fp, "<?php
\$_INFO['sql_u'] = '".$_POST['sql_u']."';
\$_INFO['sql_p'] = '".$_POST['sql_p']."';
\$_INFO['sql_h'] = '".$_POST['sql_h']."';
\$_INFO['sql_db'] = '".$_POST['sql_db']."';
\$_INFO['sql_ssl'] = false;
\$_INFO['sql_ssl_client-key'] = '/path/to/client-key.pem';
\$_INFO['sql_ssl_client-cert'] = '/path/to/client-cert.pem';
\$_INFO['sql_ssl_ca-cert'] = '/path/to/ca-cert.pem';");
                                    fclose($fp);
                                
                                    if(!rename('../../../core/framework/configuration.php.dist', '../../../core/framework/configuration.php')){
                                    
                                    	echo '<div class="alert alert-danger">Permission error encountered when trying to rename your configuration file. Please ensure its directory is 0777.</div>'; 
                                    
                                    }
                                    exit('<meta http-equiv="refresh" content="0;url=tables.php"/>');
                        
                            }catch (PDOException $e) {
                        
                                echo '<div class="alert alert-danger">MySQL Connection Error: ' . $e->getMessage() . '</div>';
                        
                            }
                            
                        }
                    
                    ?>
                    <form action="start.php" method="post">
                    	<fieldset>
                    		<div class="form-group">
                    			<label for="sql_db" class="control-label">Database Name</label>
                    			<div>
                    				<input type="text" class="form-control" name="sql_db" autocomplete="off" />
                    			</div>
                    		</div>
                    		<div class="form-group">
                    			<label for="sql_h" class="control-label">Database Host</label>
                    			<div>
                    				<input type="text" class="form-control" name="sql_h" autocomplete="off" />
                    			</div>
                    		</div>
                    		<div class="form-group">
                    			<label for="sql_u" class="control-label">Database User</label>
                    			<div>
                    				<input type="text" class="form-control" name="sql_u" autocomplete="off" />
                    			</div>
                    		</div>
                    		<div class="form-group">
                    			<label for="sql_p" class="control-label">Database User Password</label>
                    			<div>
                    				<input type="password" class="form-control" name="sql_p" autocomplete="off" />
                    			</div>
                    		</div>
                    		<div class="form-group">
                    			<div>
                    				<input type="submit" name="do_connect" value="Continue &rarr;" class="btn btn-primary" />
                    			</div>
                    		</div>
                        </fieldset>
                    </form>
				</div>
				<div class="col-2"></div>
			</div>
		</div>
		<div class="footer">
			<div class="col-8 nopad"><p>PufferPanel is licensed under a <a href="https://github.com/DaneEveritt/PufferPanel/blob/master/LICENSE">GPL-v3 License</a>.<br />Running Version 0.6.0.2 Beta distributed by <a href="http://kelp.in">Kelpin' Systems</a>.</p></div>
		</div>
	</div>
</body>
</html>
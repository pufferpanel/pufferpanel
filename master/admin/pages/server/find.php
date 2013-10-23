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
	<title>PufferPanel - Find Server</title>
	
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
			<ul id="nav" class="fl">
				<li><a href="../../../account.php" class="round button dark"><i class="icon-user"></i>&nbsp;&nbsp; <strong><?php echo $core->framework->user->getData('username'); ?></strong></a></li>
			</ul>
			<ul id="nav" class="fr">
				<li><a href="../../../servers.php" class="round button dark"><i class="icon-signout"></i></a></li>
				<li><a href="../../../logout.php" class="round button dark"><i class="icon-off"></i></a></li>
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
						<h3 class="fl"><i class="icon-cog" id="toggle_search_simple"></i>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<i class="icon-cogs" id="toggle_search_advanced"></i>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Search Parameters</h3>
					</div>
					<div class="content-module-main cf" id="search_simple">
						<form id="ss_form" onsubmit="return false">
							<fieldset>
								<div style="width:20%;display:inline-block;vertical-align:top;margin-right: 0.6em;">
									<p>
										<select class="round" name="field">
											<option value="node">Node</option>
											<option value="name" selected="selected">Server Name</option>
											<option value="server_ip">Server IP</option>
											<option value="owner_email">Owner Email</option>
											<option value="active">Active</option>
										</select><i class="icon-angle-down pull-right select-arrow"></i>
									</p>
								</div>
								<div style="width:20%;display:inline-block;vertical-align:top;margin-right: 0.6em;">
									<p>
										<select class="round" name="operator">
											<option value="equal">Equals</option>
											<option value="not_equal">Not Equal</option>
											<option value="like" selected="selected">Like</option>
											<option value="starts_w">Starts With</option>
											<option value="ends_w">Ends With</option>
										</select><i class="icon-angle-down pull-right select-arrow"></i>
									</p>
								</div>
								<div style="width:50%;display:inline-block;vertical-align:top;margin-right: 0.6em;">
									<p>
										<input class="round full-width-input" style="width: 90%;" name="term" type="text" />
									</p>
								</div>
								<div class="stripe-separator" style="margin: 0 0 1em 0;"><!--  --></div>
								<div class="confirmation-box round pull-left" id="search_active" style="display: none;margin-bottom: -1em; background: #e7fae6; padding-left: 0.833em;"><i class="icon-refresh icon-spin"></i> &nbsp;Searching!</div>
								<input type="submit" value="Perform Simple Search" class="round blue ic-right-arrow pull-right" style="margin-bottom: -1em;" />
							</fieldset>
						</form>
					</div>
					<div class="content-module-main cf" style="display: none;" id="search_advanced">
						<form id="as_form" onsubmit="return false">
							<fieldset>
								<div style="width:10%;display:inline-block;vertical-align:top;margin-right: 0.6em;">
									<p>
										<select class="round" name="field_1">
											<option value="node">Node</option>
											<option value="name">Server Name</option>
											<option value="server_ip">Server IP</option>
											<option value="owner_email">Owner Email</option>
											<option value="active">Active</option>
										</select><i class="icon-angle-down pull-right select-arrow"></i>
									</p>
								</div>
								<div style="width:10%;display:inline-block;vertical-align:top;margin-right: 0.6em;">
									<p>
										<select class="round" name="operator_1">
											<option value="equal">Equals</option>
											<option value="not_equal">Not Equal</option>
											<option value="like">Like</option>
											<option value="starts_w">Starts With</option>
											<option value="ends_w">Ends With</option>
										</select><i class="icon-angle-down pull-right select-arrow"></i>
									</p>
								</div>
								<div style="width:20%;display:inline-block;vertical-align:top;margin-right: 0.6em;">
									<p>
										<input class="round full-width-input" name="term_1" style="width: 90%;" type="text" />
									</p>
								</div>
								<div style="width:10%;display:inline-block;vertical-align:top;margin-right: 0.6em;">
									<p>
										<select class="round" name="middle_operator">
											<option value="and">And</option>
											<option value="or">Or</option>
										</select><i class="icon-angle-down pull-right select-arrow"></i>
									</p>
								</div>
								<div style="width:10%;display:inline-block;vertical-align:top;margin-right: 0.6em;">
									<p>
										<select class="round" name="field_2">
											<option value="node">Node</option>
											<option value="name">Server Name</option>
											<option value="server_ip">Server IP</option>
											<option value="owner_email">Owner Email</option>
											<option value="active">Active</option>
										</select><i class="icon-angle-down pull-right select-arrow"></i>
									</p>
								</div>
								<div style="width:10%;display:inline-block;vertical-align:top;margin-right: 0.6em;">
									<p>
										<select class="round" name="operator_2">
											<option value="equal">Equals</option>
											<option value="not_equal">Not Equal</option>
											<option value="like">Like</option>
											<option value="starts_w">Starts With</option>
											<option value="ends_w">Ends With</option>
										</select><i class="icon-angle-down pull-right select-arrow"></i>
									</p>
								</div>
								<div style="width:20%;display:inline-block;vertical-align:top;margin-right: 0.6em;">
									<p>
										<input class="round full-width-input" name="term_2" style="width:90%;" type="text" />
									</p>
								</div>
								<div class="stripe-separator" style="margin: 0 0 1em 0;"><!--  --></div>
								<div class="confirmation-box round pull-left" id="search_active_2" style="display: none;margin-bottom: -1em; background: #e7fae6; padding-left: 0.833em;"><i class="icon-refresh icon-spin"></i> &nbsp;Searching!</div>
								<input type="submit" value="Perform Advanced Search" class="round blue ic-right-arrow pull-right" style="margin-bottom: -1em;" />
							</fieldset>
						</form>
					</div>
				</div>
				<div class="content-module">
					<div class="content-module-heading cf">
						<h3 class="fl">Search Results</h3>
					</div>
					<div class="content-module-main cf" id="search_results">
						<div class="information-box round">Please enter a search command.</div>	
					</div>
				</div>
			</div>
		</div>
	</div>
	<script type="text/javascript">
		$("#toggle_search_simple").click(function(){
			$("#search_advanced").slideUp(function(){
				$("#search_simple").slideDown();
			});
		});
		$("#toggle_search_advanced").click(function(){
			$("#search_simple").slideUp(function(){
				$("#search_advanced").slideDown();
			});
		});
		$("#ss_form").submit(function(){
			$("#search_active").show();
			var search_field = $('select[name="field"] :selected').val();
			var search_operator = $('select[name="operator"] :selected').val();
			var search_term = $('input[name="term"]').val();
			
			$.ajax({
				type: "POST",
				url: "ajax/search/simple.php",
				data: { method: "simple", field: search_field, operator: search_operator, term: search_term },
			  	success: function(data) {
			  		$("#search_results").slideUp(function(){
			  			$("#search_results").html(data);
		  				$("#search_results").fadeIn(function(){
		  					$("#search_active").fadeOut();
		  					return false;
		  				});
			  		});
			  	}
			});
		});
		$("#as_form").submit(function(){
			$("#search_active_2").show();
			var search_field_1 = $('select[name="field_1"] :selected').val();
			var search_operator_1 = $('select[name="operator_1"] :selected').val();
			var search_term_1 = $('input[name="term_1"]').val();
			var middle_op = $('select[name="middle_operator"] :selected').val();
			var search_field_2 = $('select[name="field_2"] :selected').val();
			var search_operator_2 = $('select[name="operator_2"] :selected').val();
			var search_term_2 = $('input[name="term_2"]').val();
			
			$.ajax({
				type: "POST",
				url: "ajax/search/advanced.php",
				data: { method: "advanced", field_1: search_field_1, operator_1: search_operator_1, term_1: search_term_1, mid_op: middle_op, field_2: search_field_2, operator_2: search_operator_2, term_2: search_term_2},
			  	success: function(data) {
			  		$("#search_results").slideUp(function(){
			  			$("#search_results").html(data);
		  				$("#search_results").fadeIn(function(){
		  					$("#search_active_2").fadeOut();
		  					return false;
		  				});
			  		});
			  	}
			});
		});
	</script>
	<div id="footer">
		<p>Copyright &copy; 2012 - 2013. All Rights Reserved.<br />Running PufferPanel Version 0.3 Beta distributed by <a href="http://pufferfi.sh">Puffer Enterprises</a>.</p>
	</div>
</body>
</html>
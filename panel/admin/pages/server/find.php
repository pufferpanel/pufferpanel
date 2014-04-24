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
					<li class="active"><a href="#simple" data-toggle="tab">Simple Search</a></li>
					<li><a href="#advanced" data-toggle="tab">Advanced Search</a></li>
					<li><a href="#list_all" id="list_all_tab" data-toggle="tab">All Servers</a></li>
				</ul>
				<div class="tab-content">
					<div class="tab-pane active" id="simple">
						<h3>Simple Search</h3><hr />
						<form id="ss_form" onsubmit="return false">
							<fieldset>
								<div class="form-group col-3 nopad">
									<div>
										<select class="form-control" name="field">
											<option value="name" selected="selected">Server Name</option>
											<option value="server_ip">Server IP</option>
											<option value="server_port">Server Port</option>
											<option value="owner_email">Owner Email</option>
											<option value="active">Active</option>
										</select>
									</div>
								</div>
								<div class="form-group col-2">
									<div>
										<select class="form-control" name="operator">
											<option value="equal">Equals</option>
											<option value="not_equal">Not Equal</option>
											<option value="like" selected="selected">Like</option>
											<option value="starts_w">Starts With</option>
											<option value="ends_w">Ends With</option>
										</select>
									</div>
								</div>
								<div class="form-group col-7 nopad-right">
									<div class="input-group">
										<input type="text" class="form-control" name="term" />
										<span class="input-group-btn">
											<button class="btn btn-primary" id="ss_active_spin" type="submit">&rarr;</button>
										</span>
									</div>
								</div>
							</fieldset>
						</form>
					</div>
					<div class="tab-pane" id="advanced">
						<h3>Advanced Search</h3><hr />
						<form id="as_form" onsubmit="return false">
							<fieldset>
								<div class="form-group col-3 nopad">
									<div>
										<select class="form-control" name="field_1">
											<option value="name" selected="selected">Server Name</option>
											<option value="server_ip">Server IP</option>
											<option value="server_port">Server Port</option>
											<option value="owner_email">Owner Email</option>
											<option value="active">Active</option>
										</select>
									</div>
								</div>
								<div class="form-group col-2">
									<div>
										<select class="form-control" name="operator_1">
											<option value="equal">Equals</option>
											<option value="not_equal">Not Equal</option>
											<option value="like" selected="selected">Like</option>
											<option value="starts_w">Starts With</option>
											<option value="ends_w">Ends With</option>
										</select>
									</div>
								</div>
								<div class="form-group col-7 nopad-right">
									<div class="input-group">
										<input type="text" class="form-control" name="term_1" />
										<span class="input-group-btn">
											<div class="btn-group">
							                  <button type="button" class="btn btn-success dropdown-toggle" data-toggle="dropdown"><span class="caret"></span></button>
							                  <ul class="dropdown-menu" id="special_case_dropdown">
							                    <li><a href="and" id="radio_and"><input type="radio" name="middle_operator" value="and" checked="checked"/> AND</a></li>
							                    <li><a href="or" id="radio_or"><input type="radio" name="middle_operator" value="or"/> OR</a></li>
							                  </ul>
							                </div>
										</span>
									</div>
								</div>
								<div class="form-group col-3 nopad">
									<div>
										<select class="form-control" name="field_2">
											<option value="name">Server Name</option>
											<option value="server_ip">Server IP</option>
											<option value="server_port">Server Port</option>
											<option value="owner_email" selected="selected">Owner Email</option>
											<option value="active">Active</option>
										</select>
									</div>
								</div>
								<div class="form-group col-2">
									<div>
										<select class="form-control" name="operator_2">
											<option value="equal">Equals</option>
											<option value="not_equal">Not Equal</option>
											<option value="like" selected="selected">Like</option>
											<option value="starts_w">Starts With</option>
											<option value="ends_w">Ends With</option>
										</select>
									</div>
								</div>
								<div class="form-group col-7 nopad-right">
									<div class="input-group">
										<input type="text" class="form-control" name="term_2" />
										<span class="input-group-btn">
											<button class="btn btn-primary" id="as_active_spin" type="submit">&rarr;</button>
										</span>
									</div>
								</div>
							</fieldset>
						</form>
					</div>
					<div class="tab-pane" id="list_all">
						<h3>All Servers</h3><hr />
					</div>
				</div>
				<div id="search_results"></div>
			</div>
		</div>
		<div class="footer">
			<?php include('../../../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(document).ready(function(){
			$("#ss_form").submit(function(){
				$("#ss_active_spin").html('<i class="fa fa-refresh fa-spin"></i>');
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
			  					$("#ss_active_spin").html('&rarr;');
			  					return false;
			  				});
				  		});
				  	}
				});
			});
			$("#as_form").submit(function(){
				$("#as_active_spin").html('<i class="fa fa-refresh fa-spin"></i>');
				var search_field_1 = $('select[name="field_1"] :selected').val();
				var search_operator_1 = $('select[name="operator_1"] :selected').val();
				var search_term_1 = $('input[name="term_1"]').val();
				var middle_op = $('input[name="middle_operator"]:checked').val();
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
			  					$("#as_active_spin").html('&rarr;');
			  					return false;
			  				});
				  		});
				  	}
				});
			});
			$("#list_all_tab").click(function(e){
				$("#list_all_tab").append(' <i class="fa fa-refresh fa-spin"></i>');			
				$.ajax({
					type: "POST",
					url: "ajax/search/getall.php",
				  	success: function(data) {
				  		$("#search_results").slideUp(function(){
				  			$("#search_results").html(data);
			  				$("#search_results").fadeIn(function(){
			  					$("#list_all_tab").html("All Servers");
			  					return false;
			  				});
				  		});
				  	}
				});
			});
			$('#special_case_dropdown a').click(function(e) {
				var radio = $(this).attr("href");
				$("input[value='"+radio+"']").prop("checked", true)
				e.preventDefault();
			    e.stopPropagation();
			});
		});
	</script>
</body>
</html>
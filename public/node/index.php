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
require_once('../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Page\components::redirect($core->settings->get('master_url').'index.php?login');
	exit();
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<?php include('../assets/include/header.php'); ?>
	<title>PufferPanel - Manage Your Server</title>
	<script src="http://code.highcharts.com/highcharts.js"></script>
	<script src="http://code.highcharts.com/highcharts-more.js"></script>
	<script src="http://code.highcharts.com/modules/solid-gauge.src.js"></script>
</head>
<body>
	<div class="container">
		<?php include('../assets/include/navbar.php'); ?>
		<div class="row">
			<div class="col-3">
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.acc_actions'); ?></strong></a>
					<a href="../account.php" class="list-group-item"><?php echo $_l->tpl('sidebar.settings'); ?></a>
					<a href="../servers.php" class="list-group-item"><?php echo $_l->tpl('sidebar.servers'); ?></a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.server_acc'); ?></strong></a>
					<a href="index.php" class="list-group-item active"><?php echo $_l->tpl('sidebar.overview'); ?></a>
					<a href="console.php" class="list-group-item"><?php echo $_l->tpl('sidebar.console'); ?></a>
					<a href="files/index.php" class="list-group-item"><?php echo $_l->tpl('sidebar.files'); ?></a>
				</div>
				<div class="list-group">
					<a href="#" class="list-group-item list-group-item-heading"><strong><?php echo $_l->tpl('sidebar.server_sett'); ?></strong></a>
					
					<a href="settings.php" class="list-group-item"><?php echo $_l->tpl('sidebar.manage'); ?></a>
					<a href="plugins/index.php" class="list-group-item"><?php echo $_l->tpl('sidebar.plugins'); ?></a>
				</div>
			</div>
			<div class="col-9">
				<div class="col-12">
					<h3 class="nopad"><?php echo $_l->tpl('node.overview.stats_h1'); ?></h3><hr />
					<!--<div id="online_notice" class="alert alert-info"><i class="fa fa-spinner fa-spin"></i> <?php echo $_l->tpl('node.overview.collect_usage'); ?></div>-->
					<div class="row" style="margin: -30px auto -30px auto;width: 600px;">
						<div id="container-cpu-loading" class="stats_loading_box"><i class="fa fa-refresh fa-spin"></i></div>
						<div id="container-memory-loading" class="stats_loading_box" style="left: 255px;"><i class="fa fa-refresh fa-spin"></i></div>
						<div id="container-disk-loading" class="stats_loading_box" style="left: 455px;"><i class="fa fa-refresh fa-spin"></i></div>
						<div id="container-cpu" style="width: 200px; height: 200px; float: left"></div>
						<div id="container-memory" style="width: 200px; height: 200px; float: left"></div>
						<div id="container-disk" style="width: 200px; height: 200px; float: left"></div>
					</div>
				</div>
				<div class="col-12">
					<h3><?php echo $_l->tpl('node.overview.players_h5'); ?></h3><hr />
					<div id="players_notice" class="alert alert-info"><i class="fa fa-spinner fa-spin"></i> <?php echo $_l->tpl('node.overview.collect_usage'); ?></div>
					<span id="toggle_players" style="display:none;">
						<p class="text-muted"><?php echo $_l->tpl('node.overview.no_players'); ?></p>
				</div>
				<div class="col-12">
					<h3><?php echo $_l->tpl('node.overview.information_h1'); ?></h3><hr />
					<table class="table table-striped table-bordered table-hover">
						<tbody>
							<tr>
								<td><strong><?php echo $_l->tpl('string.connection'); ?></strong></td>
								<td><?php echo $core->server->getData('server_ip').':'.$core->server->getData('server_port'); ?></td>
							</tr>
							<tr>
								<td><strong><?php echo $_l->tpl('string.node'); ?></strong></td>
								<td><?php echo $core->settings->nodeName($core->server->getData('node')); ?></td>
							</tr>
							<tr>
								<td><strong><?php echo $_l->tpl('string.mem_alloc'); ?></strong></td>
								<td><?php echo $core->server->getData('max_ram').' MB'; ?></td>
							</tr>
							<tr>
								<td><strong><?php echo $_l->tpl('string.disk_alloc'); ?></strong></td>
								<td><?php echo $core->server->getData('disk_space').' MB'; ?></td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
		<div class="footer">
			<?php include('../assets/include/footer.php'); ?>
		</div>
	</div>
	<script type="text/javascript">
		$(window).load(function(){
			var socket = io.connect('http://<?php echo $core->server->nodeData('sftp_ip'); ?>:8031/<?php echo $core->server->getData('gsd_id'); ?>');
			socket.on('process', function (data) {
				if($("#container-cpu-loading").is(":visible")){
					$("#container-cpu-loading").fadeOut();
				}
				if($("#container-memory-loading").is(":visible")){
					$("#container-memory-loading").fadeOut();
				}
				var chart = $('#container-cpu').highcharts();
				if(chart){
					var point = chart.series[0].points[0];
					point.update(data.process.cpu);
				}
				var chart = $('#container-memory').highcharts();
				if(chart){
					var point = chart.series[0].points[0];
					if(parseInt(data.process.memory / (1024 * 1024)) > <?php echo (int) $core->server->getData('max_ram'); ?>){
						point.update(<?php echo (int) $core->server->getData('max_ram'); ?>);
					}else{
						point.update(parseInt(data.process.memory / (1024 * 1024)));
					}
				}
			});
			socket.on('query', function (data) {
				if($("#players_notice").is(":visible")){
					$("#players_notice").hide();
					$("#toggle_players").show();
				}
				if(data.query.players.length !== 0){
					$("#toggle_players").html("");
					$.each(data.query.players, function(id, name) {
						$("#toggle_players").append('<img data-toggle="tooltip" src="http://i.fishbans.com/helm/'+name+'/32" title="'+name+'" style="padding: 0 2px 6px 0;"/>');
					});
				}else{
					$("#toggle_players").html('<p class="text-muted">No players are currently online.</p>');
				}
				$("img[data-toggle='tooltip']").tooltip();
			});
			$(function () {
			    var gaugeOptions = {
				    chart: {
				        type: 'solidgauge'
				    },
				    title: null,
				    pane: {
				    	center: ['50%', '85%'],
				    	size: '100%',
				        startAngle: -90,
				        endAngle: 90,
			            background: {
			                backgroundColor: (Highcharts.theme && Highcharts.theme.background2) || '#EEE',
			                innerRadius: '60%',
			                outerRadius: '100%',
			                shape: 'arc'
			            }
				    },
				    tooltip: {
				    	enabled: false
				    },
				    yAxis: {
						stops: [
							[0.3, '#06BA00'], // green
				        	[0.6, '#DBC200'], // yellow
				        	[0.9, '#BA0000'] // red
						],
						lineWidth: 0,
			            minorTickInterval: null,
			            tickPixelInterval: 400,
			            tickWidth: 0,
				        title: {
			                y: -70
				        },
			            labels: {
			                y: 16
			            }        
				    },
			        plotOptions: {
			            solidgauge: {
			                dataLabels: {
			                    y: -30,
			                    borderWidth: 0,
			                    useHTML: true
			                }
			            }
			        }
			    };
			    $('#container-cpu').highcharts(Highcharts.merge(gaugeOptions, {
			        yAxis: {
			        	min: 0,
			        	max: 100,
			            title: {
			                text: 'CPU Usage'
			            },
			            labels: {
			            	enabled: false
			            }    
			        },
			        credits: {
			        	enabled: false
			        },
			        series: [{
			            name: 'CPU',
			            data: [0],
			            dataLabels: {
			            	format: '<div style="text-align:center"><span style="font-size:16px;color:' + 
			                    ((Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black') + '">{y} %</span><br/>' + 
			                   	'</div>'
			            }
			        }]
			    }));
			    $('#container-memory').highcharts(Highcharts.merge(gaugeOptions, {
			        yAxis: {
				        min: 0,
				        max: <?php echo (int) $core->server->getData('max_ram'); ?>,
				        title: {
				            text: 'Memory Usage'
				        },
				        labels: {
				        	enabled: false
				        }
				    },
				    credits: {
				    	enabled: false
				    },
				    series: [{
				        name: 'Memory',
				        data: [0],
				        dataLabels: {
				        	format: '<div style="text-align:center"><span style="font-size:16px;color:' + 
			                    ((Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black') + '">{y} MB</span><br/>' + 
			                   	'<span style="font-size:10px;color:silver">/ 128 MB</span></div>'
				        }
				    }]
				
				}));
				$('#container-disk').highcharts(Highcharts.merge(gaugeOptions, {
				    yAxis: {
				        min: 0,
				        max: <?php echo (int) $core->server->getData('disk_space'); ?>,
				        title: {
				            text: 'Disk Usage'
				        },
				        labels: {
				        	enabled: false
				        }
				    },
				    credits: {
				    	enabled: false
				    },
				    series: [{
				        name: 'Disk',
				        data: [0],
				        dataLabels: {
				        	format: '<div style="text-align:center"><span style="font-size:16px;color:' + 
				                ((Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black') + '">{y} MB</span><br/>' + 
				               	'<span style="font-size:10px;color:silver">/ <?php echo $core->server->getData('disk_space'); ?> MB</span></div>'
				        }
				    }]
				
				}));
			});
			$.ajax({
				type: "POST",
				url: "ajax/overview/data.php",
				data: { command: 'stats' },
			  		success: function(data) {
						var chart = $('#container-disk').highcharts();
						if(chart){
							$("#container-disk-loading").fadeOut();
							var point = chart.series[0].points[0];
							point.update(parseInt(data));
						}
			 		}
			});
		});
	</script>
</body>
</html>
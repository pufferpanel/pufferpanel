@extends('layouts.master')

@section('title')
    Viewing Server: {{ $server->uuidShort }}
@endsection

@section('content')
<div class="col-md-9">
	<ul class="nav nav-tabs" id="config_tabs">
		<li class="active"><a href="#stats" data-toggle="tab">Statistics</a></li>
		<li><a href="#console" data-toggle="tab">Control Server</a></li>
		<li><a href="#remote" data-toggle="tab">Remote Requests</a></li>
	</ul><br />
	<div class="tab-content">
		<div class="tab-pane active" id="stats">
			<div class="row">
				<div class="col-md-6">
					<h3 class="nopad">Memory Usage</h3><hr />
					<div class="row centered">
						<canvas id="memoryChart" width="280" height="150" style="margin-left:20px;"></canvas>
						<p style="text-align:center;margin-top:-15px;" class="text-muted"><small>Time (2s Increments)</small></p>
						<p class="graph-yaxis hidden-xs hidden-sm text-muted" style="margin-top:-50px !important;"><small>Memory Usage (Mb)</small></p>
						<p class="graph-yaxis hidden-lg hidden-md text-muted" style="margin-top:-65px !important;margin-left: 100px !important;"><small>Memory Usage (%)</small></p>
					</div>
				</div>
				<div class="col-md-6">
					<h3 class="nopad">CPU Usage</h3><hr />
					<div class="row centered">
						<canvas id="cpuChart" width="280" height="150" style="margin-left:20px;"></canvas>
						<p style="text-align:center;margin-top:-15px;" class="text-muted"><small>Time (2s Increments)</small></p>
						<p class="graph-yaxis hidden-sm hidden-xs text-muted" style="margin-top:-65px !important;"><small>CPU Usage (%)</small></p>
						<p class="graph-yaxis hidden-lg hidden-md text-muted" style="margin-top:-65px !important;margin-left: 100px !important;"><small>CPU Usage (%)</small></p>
					</div>
				</div>
			</div>
			<div class="row">
				<div class="col-md-12" id="stats_players" style="display: none;">
					<h3 class="nopad">Active Players</h3><hr />
					<div id="players_notice" class="alert alert-info"><i class="fa fa-spinner fa-spin"></i>Currently Collecting Usage Information</div>
					<span id="toggle_players" style="display:none;">
						<p class="text-muted">No players are online.</p>
				</div>
				<div class="col-md-12">
					<h3>Server Information</h3><hr />
					<table class="table table-striped table-bordered table-hover">
						<tbody>
							<tr>
								<td><strong>Default Connection</strong></td>
								<td><code>{{ $server->ip }}:{{ $server->port }}</code></td>
							</tr>
							<tr>
								<td><strong>Node</strong></td>
								<td>{{ $node->name }}</td>
							</tr>
							<tr>
								<td><strong>Memory Limit</strong></td>
								<td>{{ $server->memory }} MB</td>
							</tr>
							<tr>
								<td><strong>Disk Space</strong></td>
								<td>{{ $server->disk }} MB</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
		<div class="tab-pane" id="console">
			<div class="row">
				<div class="col-md-12">
					<textarea id="live_console" class="form-control console" readonly="readonly"></textarea>
				</div>
				<div class="col-md-6">
					<hr />
					<form action="#" method="post" id="console_command">
						<fieldset>
							<div class="input-group">
								<input type="text" class="form-control" name="command" id="ccmd" placeholder="command here" />
								<span class="input-group-btn">
									<button id="sending_command" class="btn btn-primary btn-sm">&rarr;</button>
								</span>
							</div>
						</fieldset>
					</form>
					<div class="alert alert-danger" id="sc_resp" style="display:none;margin-top: 15px;"></div>
				</div>
				<div class="col-md-6" style="text-align:center;">
					<hr />
					<button class="btn btn-success btn-sm start disabled" id="server_start">Start</button>
					<button class="btn btn-primary btn-sm restart disabled" id="server_restart">Restart</button>
					<button class="btn btn-danger btn-sm stop disabled" id="server_stop">Stop</button>
					<button class="btn btn-primary btn-sm" data-toggle="modal" data-target="#pauseConsole" id="pause_console"><small><i class="fa fa-pause"></i></small></button>
					<div style="margin-top:5px;" id="kill_process_text" style="display:none;"><small><p class="text-muted">My server isn't responding! Please <code id="kill_proc" style="cursor: pointer;">kill it</code>.</p></small></div>
					<div id="pw_resp" style="display:none;margin-top: 15px;"></div>
				</div>
			</div>
		</div>
		<div class="tab-pane" id="remote">
			<h3>Remote Connection Information</h3>
			<p>Remote connections allows you to optionally control your server from outside of the panel. Using this feature you can create custom scripts to start or stop your server at specified times, or even track server information on another site.</p>
			<div class="panel panel-default">
				<div class="panel-heading">
					<h3 class="panel-title">Connection Information</h3>
				</div>
				<div class="panel-body">
					<div class="alert alert-info">This is still a work in progress, and as such it is not neccessarily fully documented.</div>
					<table class="table table-hover" style="margin-bottom:0;">
						<tr>
							<td style="border-top:0;">Connection Base</td>
							<td style="border-top:0;">https://{{ $node->fqdn }}:{{ $node->daemonListen }}/</td>
						</tr>
						<tr>
							<td>X-Access-Server</td>
							<td><code>{{ $server->uuid }}</code></td>
						</tr>
						<tr>
							<td>X-Access-Token</td>
							<td><code>{{ $server->daemonSecret }}</code></td>
						</tr>
					</table>
				</div>
			</div>
			<div class="panel panel-default">
				<div class="panel-heading">
					<h3 class="panel-title">Connection Endpoints</h3>
				</div>
				<div class="panel-body">
					<table class="table table-hover" style="margin-bottom:0;">
						<thead>
							<tr>
								<th>Action</th>
								<th>Method</th>
								<th>Endpoint</th>
								<th></th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td>Basic Server Information</td>
								<td><code>GET</code></td>
								<td><code>/server</code></td>
								<td><a href="http://scales.pufferpanel.com/docs/server" target="_blank"><i class="fa fa-link"></i></a></td>
							</tr>
							<tr>
								<td>Get Last Lines of Console</td>
								<td><code>GET</code></td>
								<td><code>/server/log/:lines</code></td>
								<td><a href="http://scales.pufferpanel.com/docs/serverloglines" target="_blank"><i class="fa fa-link"></i></a></td>
							</tr>
							<tr>
								<td>Control Server Power</td>
								<td><code>GET</code></td>
								<td><code>/server/power/(start|stop|restart|kill)</code></td>
								<td><a href="http://scales.pufferpanel.com/docs/serverpowerstatus" target="_blank"><i class="fa fa-link"></i></a></td>
							</tr>
							<tr>
								<td>Send Console Commands</td>
								<td><code>POST</code></td>
								<td><code>/server/console</code></td>
								<td><a href="#" target="_blank"><i class="fa fa-link"></i></a></td>
							</tr>
							<tr>
								<td>View Server Files</td>
								<td><code>GET</code></td>
								<td><code>/server/directory/:directory</code></td>
								<td><a href="http://scales.pufferpanel.com/docs/serverdirectorydirectory" target="_blank"><i class="fa fa-link"></i></a></td>
							</tr>
							<tr>
								<td>Edit Server File</td>
								<td><code>GET</code></td>
								<td><code>/server/file/:file</code></td>
								<td><a href="http://scales.pufferpanel.com/docs/serverfilefile" target="_blank"><i class="fa fa-link"></i></a></td>
							</tr>
							<tr>
								<td>Save Server File</td>
								<td><code>PUT</code></td>
								<td><code>/server/file/:file</code></td>
								<td><a href="#" target="_blank"><i class="fa fa-link"></i></a></td>
							</tr>
							<tr>
								<td>Delete Server File</td>
								<td><code>DELETE</code></td>
								<td><code>/server/file/:file</code></td>
								<td><a href="#" target="_blank"><i class="fa fa-link"></i></a></td>
							</tr>
							<tr>
								<td>Change SFTP Password</td>
								<td><code>POST</code></td>
								<td><code>/server/reset-password</code></td>
								<td><a href="#" target="_blank"><i class="fa fa-link"></i></a></td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
	</div>
</div>
<div class="modal fade" id="pauseConsole" tabindex="-1" role="dialog" aria-labelledby="PauseConsole" aria-hidden="true">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
				<h4 class="modal-title" id="PauseConsole">ScrollStop</h4>
			</div>
			<div class="modal-body">
				<div class="row">
					<div class="col-md-12">
						<textarea id="paused_console" class="form-control console" readonly="readonly"></textarea>
					</div>
				</div>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
			</div>
		</div>
	</div>
</div>
<script>
$(document).ready(function () {
    $('#sidebar_links').find('a[href=\'/node/index\']').addClass('active');
});
</script>
@endsection

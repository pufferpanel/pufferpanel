<!DOCTYPE html>
<html lang="en">
<head>
	@section('scripts')
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta name="robots" content="noindex">
		<link rel="stylesheet" href="http://pffr.me/assets/css/bootstrap.css">
		<link rel="stylesheet" href="http://pffr.me/assets/css/pufferpanel.css">
		<link rel="stylesheet" href="http://pffr.me/assets/css/animate.css">
		<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/font-awesome/4.4.0/css/font-awesome.min.css">
		<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.2/js/bootstrap.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/socket.io/1.3.7/socket.io.min.js"></script>
		<script src="http://pffr.me/assets/javascript/admin.min.js"></script>
		<script src="http://pffr.me/assets/javascript/bootstrap-notify.min.js"></script>
	@show
	<title>PufferPanel - @yield('title')</title>
</head>
<body>
	<div class="container">
		<div class="navbar navbar-default">
			<div class="navbar-header">
				<button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-responsive-collapse">
					<span class="icon-bar"></span>
					<span class="icon-bar"></span>
					<span class="icon-bar"></span>
				</button>
				<a class="navbar-brand" href="/index">PufferPanel - Laravel</a>
			</div>
			<div class="navbar-collapse collapse navbar-responsive-collapse">
				@section('server-name')
					@if (isset($server->name))
						<ul class="nav navbar-nav">
							<li class="active" id="{{ $server->name }}"><a href="/server/{{ $server->id }}/index"><i id="applyUpdate" class="fa fa-circle-o-notch fa-spinner fa-spin spin-light"></i> {{ $server->name }}</a></li>
						</ul>
					@endif
				@show
				@section('right-nav')
					<ul class="nav navbar-nav navbar-right">
						<li class="dropdown">
							<a href="#" class="dropdown-toggle" data-toggle="dropdown">Langauge<b class="caret"></b></a>
							<ul class="dropdown-menu">
								<li><a href="/language/de">Deutsch</a></li>
								<li><a href="/language/en">English</a></li>
								<li><a href="/language/es">Espa&ntilde;ol</a></li>
								<li><a href="/language/fr">Fran&ccedil;ais</a></li>
								<li><a href="/language/it">Italiano</a></li>
								<li><a href="/language/pl">Polski</a></li>
								<li><a href="/language/pt">Portugu&ecirc;s</a></li>
								<li><a href="/language/ru">&#1088;&#1091;&#1089;&#1089;&#1082;&#1080;&#1081;</a></li>
								<li><a href="/language/se">Svenska</a></li>
								<li><a href="/language/zh">&#20013;&#22269;&#30340;的</a></li>
							</ul>
						</li>
						<li class="hidden-xs"><a href="/admin/index"><i class="fa fa-cogs"></i></a></li>
						<li class="hidden-xs"><a href="/auth/logout"><i class="fa fa-power-off"></i></a></li>
					</ul>
				@show
			</div>
		</div>
		<!-- Add Back Mobile Support -->
		<div class="row">
			<div class="col-md-3 hidden-xs hidden-sm" id="sidebar_links">
				@section('sidebar')
					<div class="list-group">
						<a href="#" class="list-group-item list-group-item-heading"><strong>Account Controls</strong></a>
						<a href="/account" class="list-group-item">Account Settings</a>
						<a href="/totp" class="list-group-item">Account Security</a>
						<a href="/index" class="list-group-item">Your Servers</a>
					</div>
					@if (isset($server->name))
						<div class="list-group">
							<a href="#" class="list-group-item list-group-item-heading"><strong>Server Controls</strong></a>
							<a href="/node/index" class="list-group-item">Server Overview</a>
							<a href="/node/files" class="list-group-item">File Manager</a>
							<a href="/node/users" class="list-group-item">Manage Sub-Users</a>
							<a href="/node/settings" class="list-group-item">Manage Server</a>
						</div>
					@endif
				@show
			</div>
			@yield('content')
		</div>
		<div class="footer">
			<div class="row">
			</div>
		</div>
	</div>
</body>
</html>

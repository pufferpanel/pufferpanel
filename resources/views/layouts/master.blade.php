<!DOCTYPE html>
<html lang="en">
<head>
    @section('scripts')
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta name="robots" content="noindex">
        <link rel="stylesheet" href="{{ asset('css/bootstrap.css') }}">
        <link rel="stylesheet" href="{{ asset('css/pufferpanel.css') }}">
        <link rel="stylesheet" href="{{ asset('css/animate.css') }}">
        <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/font-awesome/4.4.0/css/font-awesome.min.css">
        <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/font-awesome/4.4.0/css/font-awesome.min.css">
        <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/font-awesome/4.4.0/css/font-awesome.min.css">
        <script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.2/js/bootstrap.min.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/socket.io/1.3.7/socket.io.min.js"></script>
        <script src="{{ asset('js/admin.min.js') }}"></script>
        <script src="{{ asset('js/bootstrap-notify.min.js') }}"></script>
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
                <a class="navbar-brand" href="/">PufferPanel - Laravel</a>
            </div>
            <div class="navbar-collapse collapse navbar-responsive-collapse">
                @section('server-name')
                    @if (isset($server->name) && isset($node->name))
                        <ul class="nav navbar-nav">
                            <li class="active" id="{{ $server->name }}"><a href="/server/{{ $server->id }}/index"><i id="applyUpdate" class="fa fa-circle-o-notch fa-spinner fa-spin spin-light"></i> {{ $server->name }}</a></li>
                        </ul>
                    @endif
                @show
                @section('right-nav')
                    <ul class="nav navbar-nav navbar-right">
                        <li class="dropdown">
                            <a href="#" class="dropdown-toggle" data-toggle="dropdown">{{ trans('strings.language') }}<b class="caret"></b></a>
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
                        <a href="#" class="list-group-item list-group-item-heading"><strong>{{ trans('pagination.sidebar.account_controls') }}</strong></a>
                        <a href="/account" class="list-group-item">{{ trans('pagination.sidebar.account_settings') }}</a>
                        <a href="/account/totp" class="list-group-item">{{ trans('pagination.sidebar.account_security') }}</a>
                        <a href="/" class="list-group-item">{{ trans('pagination.sidebar.servers') }}</a>
                    </div>
                    @section('sidebar-server')
                        @if (isset($server->name) && isset($node->name))
                            <div class="list-group">
                                <a href="#" class="list-group-item list-group-item-heading"><strong>{{ trans('pagination.sidebar.server_controls') }}</strong></a>
                                <a href="/server/{{ $server->uuidShort }}/" class="list-group-item server-index">{{ trans('pagination.sidebar.overview') }}</a>
                                <a href="/server/{{ $server->uuidShort }}/files" class="list-group-item server-files">{{ trans('pagination.sidebar.files') }}</a>
                                <a href="/server/{{ $server->uuidShort }}/users" class="list-group-item server-users">{{ trans('pagination.sidebar.subusers') }}</a>
                                <a href="/server/{{ $server->uuidShort }}/settings" class="list-group-item server-settings">{{ trans('pagination.sidebar.manage') }}</a>
                            </div>
                        @endif
                    @show
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

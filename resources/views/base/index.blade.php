@extends('layouts.master')

@section('title', 'Your Servers')

@section('sidebar-server')
@endsection

@section('content')
<div class="col-md-9">
    @if (Auth::user()->root_admin == 1)
        <div class="alert alert-info">
            You are viewing this server listing as an admin. As such, all servers installed on the system are displayed. Any servers that you are set as the owner of are marked with a blue dot to the left of their name.
        </div>
    @endif
    @if (!$servers->isEmpty())
        <table class="table table-striped table-bordered table-hover">
            <thead>
                <tr>
                    @if (Auth::user()->root_admin == 1)
                        <th></th>
                    @endif
                    <th>Server Name</th>
                    <th>Location</th>
                    <th>Node</th>
                    <th>Connection</th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                @foreach ($servers as $server)
                    <tr class="dynUpdate" id="{{ $server->uuidShort }}">
                        @if (Auth::user()->root_admin == 1)
                            <td style="width:26px;">
                                @if ($server->owner === Auth::user()->id)
                                    <i class="fa fa-circle" style="color:#008cba;"></i>
                                @else
                                    <i class="fa fa-circle" style="color:#ddd;"></i>
                                @endif
                            </td>
                        @endif
                        <td><a href="/server/{{ $server->uuidShort }}">{{ $server->name }}</a></td>
                        <td>{{ $server->location }}</td>
                        <td>{{ $server->nodeName }}</td>
                        <td><code>{{ $server->ip }}:{{ $server->port }}</code></td>
                        <td style="width:26px;"><i class="fa fa-circle-o-notch fa-spinner fa-spin applyUpdate"></i></td>
                    </tr>
                @endforeach
            </tbody>
        </table>
    @else
        <div class="alert alert-info">
            You do not currently have any servers listed on your account.
        </div>
    @endif
</div>
<script>
$(window).load(function () {
    $('#sidebar_links').find('a[href=\'/\']').addClass('active');
    function updateServerStatus () {
        $('.dynUpdate').each(function (index, data) {

            var element = $(this);
            var serverShortUUID = $(this).attr('id');
            var updateElement = $(this).find('.applyUpdate');

            updateElement.removeClass('fa-check-circle fa-times-circle').css({ color: '#000' });
            updateElement.addClass('fa-circle-o-notch fa-spinner fa-spin');

            $.ajax({
                type: 'GET',
                url: '/server/' + serverShortUUID + '/ajax/status',
                timeout: 10000
            }).done(function (data) {

                var selector = (data == 'true') ? 'fa-check-circle' : 'fa-times-circle';
                var selectorColor = (data == 'true') ? 'rgb(83, 179, 12)' : 'rgb(227, 50, 0)';

                updateElement.removeClass('fa-circle-o-notch fa-spinner fa-spin');
                updateElement.addClass(selector).css({ color: selectorColor });

            }).fail(function (jqXHR) {

                updateElement.removeClass('fa-circle-o-notch fa-spinner fa-spin');
                updateElement.addClass('fa-question-circle').css({ color: 'rgb(227, 50, 0)' });

            });

        });
    }
    updateServerStatus();
    setInterval(updateServerStatus, 30000);
});
</script>
@endsection

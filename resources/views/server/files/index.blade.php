@extends('layouts.master')

@section('title')
    Managing Files for: {{ $server->name }}
@endsection

@section('content')
<div class="col-md-9" id="internal_alert">
    <div class="alert alert-info">
        <i class="fa fa-spinner fa-spin"></i> {{ trans('server.files.loading') }}
    </div>
</div>
<div class="col-md-9">
    <div class="files_loading_box"><i class="fa fa-refresh fa-spin" id="position_me"></i></div>
</div>
<div class="col-md-9" id="load_files"></div>
<script>
    $(document).ready(function () {
        $('.server-files').addClass('active');
    });
    $(window).load(function(){
        var doneLoad = false;

        // Show Loading Animation
        function handleLoader (show) {

            // Hide animation if no files displayed.
            if ($('#load_files').height() < 5) { return; }

            // Show Animation
            if (show === true){
                var height = $('#load_files').height();
                var width = $('.files_loading_box').width();
                var center_height = (height / 2) - 30;
                var center_width = (width / 2) - 30;
                $('#position_me').css({
                    'top': center_height,
                    'left': center_width,
                    'font-size': '60px'
                });
                $(".files_loading_box").css('height', (height + 5)).fadeIn();
            } else {
                $('.files_loading_box').fadeOut(100);
            }

        }

        // Handle folder clicking to load new contents
        function reloadActionClick() {
            $('a.load_new').click(function (e) {
                e.preventDefault();
                window.history.pushState(null, null, $(this).attr('href'));
                loadDirectoryContents($.urlParam('dir', $(this).attr('href')));
            });
        }

        //
        function loadDirectoryContents (dir) {

            handleLoader(true);
            var outputContent;
            var urlDirectory = (dir === null) ? '/' : dir;

            $.ajax({
                type: 'POST',
                url: '/server/{{ $server->uuidShort }}/ajax/files/directory-list',
                headers: { 'X-CSRF-Token': '{{ csrf_token() }}' },
                data: { directory: urlDirectory }
            }).done(function (data) {
                handleLoader(false);
                $("#load_files").slideUp(function () {
                    $("#load_files").html(data).slideDown();
                    $('#internal_alert').slideUp();
                    reloadActionClick();
                });
            }).fail (function (jqXHR) {
                $("#internal_alert").html('<div class="alert alert-danger">An error occured while attempting to process this request. Please try again.</div>').show();
                console.log(jqXHR);
            });

        }
        loadDirectoryContents($.urlParam('dir'));
    });
</script>
@endsection

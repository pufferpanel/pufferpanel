@extends('layouts.master')

@section('title', 'Your Servers')

@section('content')
<div class="col-md-9">
	@if (Auth::user()->root_admin == 1)
		<div class="alert alert-info">
			You are viewing this server listing as an admin. As such, all servers installed on the system are displayed. Any servers that you are set as the owner of are marked with a blue dot to the left of their name.
		</div>
	@endif
	<div class="alert alert-info">
		You do not currently have any servers listed on your account.
	</div>
</div>
<script type="text/javascript">
$(document).ready(function() {
	$("#sidebar_links").find("a[href='/']").addClass('active');
});
</script>
@endsection

@extends('layouts.master')

@section('title', 'Your Account')

@section('sidebar-server')
@endsection

@section('content')
<div class="col-md-9">
    @if (session('flash-error'))
        <div class="alert alert-danger alert-dismissible" role="alert">
            <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
            {{ session('flash-error') }}
        </div>
    @endif
    @if (session('flash-success'))
        <div class="alert alert-success alert-dismissible" role="alert">
            <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
            {{ session('flash-success') }}
        </div>
    @endif
	<div class="row">
		<div class="col-md-6">
			<h3 class="nopad">{{ trans('base.account.update_pass') }}</h3><hr />
				<form action="/account/password" method="post">
					<div class="form-group">
						<label for="current_password" class="control-label">{{ trans('strings.current_password') }}</label>
						<div>
							<input type="password" class="form-control" name="current_password" />
						</div>
					</div>
					<div class="form-group">
						<label for="new_password" class="control-label">{{ trans('base.account.new_password') }}</label>
						<div>
							<input type="password" class="form-control" name="new_password" />
						</div>
					</div>
					<div class="form-group">
						<label for="new_password_again" class="control-label">{{ trans('base.account.new_password') }} {{ trans('strings.again') }}</label>
						<div>
							<input type="password" class="form-control" name="new_password_again" />
						</div>
					</div>
					<div class="form-group">
						<div>
							{!! csrf_field() !!}
							<input type="submit" class="btn btn-primary btn-sm" value="{{ trans('base.account.update_pass') }}" />
						</div>
					</div>
				</form>
		</div>
		<div class="col-md-6">
			<h3 class="nopad">{{ trans('base.account.update_email') }}</h3><hr />
				<form action="/account/email" method="post">
					<div class="form-group">
						<label for="new_email" class="control-label">{{ trans('base.account.new_email') }}</label>
						<div>
							<input type="text" class="form-control" name="new_email" />
						</div>
					</div>
					<div class="form-group">
						<label for="password" class="control-label">{{ trans('strings.current_password') }}</label>
						<div>
							<input type="password" class="form-control" name="password" />
						</div>
					</div>
					<div class="form-group">
						<div>
							{!! csrf_field() !!}
							<input type="submit" class="btn btn-primary btn-sm" value="{{ trans('base.account.update_email') }}" />
						</div>
					</div>
				</form>
		</div>
	</div>
</div>
<script>
$(document).ready(function () {
    $('#sidebar_links').find('a[href=\'/account\']').addClass('active');
});
</script>
@endsection

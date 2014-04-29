<div class="row">
	<div class="col-8 ">
		<p><?php echo $_l->tpl('footer.license'); ?><br /><?php echo sprintf($_l->tpl('footer.version'), '0.6.1 Beta'); ?></p>
	</div>
	<div class="col-4">
		<p class="pull-right"><?php echo sprintf($_l->tpl('footer.generated'), number_format((microtime(true) - $pageStartTime), 4)); ?><br /><?php echo sprintf($_l->tpl('footer.queries'), Database\databaseInit::getCount()); ?></p>
	</div>
</div>
<script type="text/javascript">
if($.urlParam('tab') != null){$('#config_tabs a[href="#'+$.urlParam('tab')+'"]').tab('show');}
</script>
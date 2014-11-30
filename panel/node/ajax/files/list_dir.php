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
 
namespace PufferPanel\Core;

$klein->respond('*', function($request, $response) use ($core, $twig) {
	if(!$core->user->hasPermission('files.view')) {
		exit('<div class="alert alert-danger">You do not have permission to save files.</div>');
	}

	/*
	* Set Defaults
	*/
	$displayFolders = array();
	$displayFiles = array();
	$entries = array();
	$previousDir = array();

	if(isset($_POST['dir']) && !empty($_POST['dir'])) {
		$_POST['dir'] = str_replace('..', '', urldecode(rtrim($_POST['dir'], '/')));
	} else{
		$_POST['dir'] = null;
	}

	/*
	* Gather Files and Folders
	*/
	$getDirectory = (is_null($_POST['dir'])) ? "/" : $_POST['dir'];
	$url = "http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id')."/folder/".$getDirectory;

	$context = stream_context_create(array(
		"http" => array(
			"method" => "GET",
			"header" => 'X-Access-Token: '.$core->server->getData('gsd_secret'),
			"ignore_errors" => true,
			"timeout" => 3
		)
	));

	$rawcontent = @file_get_contents($url, 0, $context);
	$content = json_decode($rawcontent, true);

	if(json_last_error() != JSON_ERROR_NONE) {
		exit('<div class="alert alert-danger">GSD ERROR: '.$rawcontent.'</div>');
	}

	if(!is_array($content)) {
		exit('<div class="alert alert-danger">'.$_l->tpl('node_files_ajax_no_dl').' This usually occurs because of a networking error.</div>');
	}

	/*
	* Parse Through Files
	*/
	if(isset($content['code']) && isset($content['message'])) {
		exit('<div class="alert alert-danger">GSD ERROR: '.$content['message'].'</div>');
	}
	
	/*
	* Handle Directory
	*/
	if(isset($_POST['dir']) && !empty($_POST['dir'])) {

		/*
		* In dir, show first arrow in display
		*/
		$previousDir['first'] = true;

		/*
		* Check First Character
		*/
		if(substr($_POST['dir'], 0, 1) == '/') {
			$_POST['dir'] = substr($_POST['dir'], 1);
		} else {
			$_POST['dir'] = $_POST['dir'].'/';
		}

	}

	/*
	* Inside a Directory
	*/
	$goBack = explode('/', $_POST['dir']);
	if(array_key_exists(1, $goBack) && !empty($goBack[1])){

		/*
		* Do we show previous-dir arrow?
		*/
		if(!empty($goBack[1])) {
			$previousDir['show'] = true;
		}

		unset($goBack[count($goBack) - 2]);
		$previousDir['link'] = rtrim(implode('/', $goBack), '/');

	}

	/*
	* Setting More Variables
	*/
	if(array_key_exists('link', $previousDir) && strpos(rtrim($previousDir['link'], '/'), '/')) {
		$previousDir['link_show'] = end(explode('/', $previousDir['link']));
	} elseif(array_key_exists('link', $previousDir)) {
		$previousDir['link_show'] = $previousDir['link'];
	}

	/*
	* Loop Through
	*/
	foreach($content as $value) {

		/*
		* Iterate into Arrays
		*/
		if($value['filetype'] == 'folder') {

			$displayFolders = array_merge($displayFolders, array(array(
				"entry" => $value['name'],
				"directory" => $_POST['dir'],
				"size" => null,
				"date" => strtotime($value['mtime'])
			)));

		} else {

			$displayFiles = array_merge($displayFiles, array(array(
				"entry" => $value['name'],
				"directory" => $_POST['dir'],
				"extension" => pathinfo($value['name'], PATHINFO_EXTENSION),
				"size" => $core->files->formatSize($value['size']),
				"date" => strtotime($value['mtime'])
			)));

		}

	}

	/*
	* Render Page
	*/
	echo $twig->render('node/files/ajax/list_dir.html', array(
		'files' => $displayFiles,
		'folders' => $displayFolders,
		'extensions' => array('txt', 'yml', 'log', 'conf', 'html', 'json', 'properties', 'props', 'cfg', 'lang'),
		'directory' => $previousDir
	));
});

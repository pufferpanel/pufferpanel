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
use \ORM as ORM;

require_once('../../../src/core/core.php');

$filesIncluded = true;

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false)
	Components\Page::redirect($core->settings->get('master_url').'index.php?login');

$results = array();
$pluginAuthors = null;
$i = 0;

if(!isset($_GET['slug']) || empty($_GET['slug']))
	Components\Page::redirect('search.php');
else {

	/*
	 * Viewing Plugin
	 */
	$_GET['slug'] = str_replace(array(' ', '+', '%20'), '', $_GET['slug']);
	$context = stream_context_create(array(
		"http" => array(
			"method" => "GET",
			"header" => 'User-Agent: PufferPanel',
			"timeout" => 5
		)
	));
	$data = json_decode(file_get_contents('http://api.bukget.org/3/plugins/bukkit/'.$_GET['slug'], false, $context), true);
	
	/*
	 * Handle Authors
	 */
	if(empty($data['authors']))
		$pluginAuthors = 'none specified';
	else
		foreach($data['authors'] as $id => $name)
			$pluginAuthors .= '<a href="http://dev.bukkit.org/profiles/'.$name.'/" target="_blank">'.$name.'</a>, ';

	/*
	 * Parse Data
	 */
	$data['versions'] = (is_array($data['versions'])) ? $data['versions'] : array($data['versions']);
	
	foreach($data['versions'] as $id => $value){
		
		$results = array_merge($results, array(array(
			'id' => $i,
			'filename' => $value['filename'],
			'version' => $value['version'],
			'date' => $value['date'],
			'versions' => $value['game_versions'],
			'md5' => $value['md5']
		)));
		
		$i++;
	
	}

}

/*
 * Display Page
 */
echo $twig->render(
		'node/plugins/view.html', array(
			'plugin' => array(
				'slug' => $data['slug'],
				'name' => $data['plugin_name'],
				'authors' => rtrim($pluginAuthors, ', '),
				'description' => (strlen($data['description']) == 0) ? 'No description is avaliable for this plugin.' : $data['description']
			),
			'results' => $results,
			'footer' => array(
				
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
	));
?>

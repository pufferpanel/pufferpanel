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

require_once('../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false)
	Components\Page::redirect($core->settings->get('master_url').'index.php?login');

$results = array();
$error = null;

if(!isset($_GET['term']) || empty($_GET['term']))
	$error = true;
else

	/*
	 * Searching for Plugin
	 */
	((isset($_GET['start']) && $_GET['start'] >= 1) ? $_GET['start'] = $_GET['start'] : $_GET['start'] = '0');
	$_GET['term'] = str_replace(array(' ', '+', '%20'), '', $_GET['term']);
	$data = file_get_contents('http://api.bukget.org/3/search/plugin_name/like/'.$_GET['term'].'?start='.$_GET['start'].'&size=25');
	$data = json_decode($data, true);
	
		foreach($data as $item => $value){
			
			$results = array_merge($results, array(array(
				'slug' => $value['slug'],
				'name' => $value['plugin_name'],
				'description' => (strlen($value['description']) == 0) ? 'No description is avaliable for this plugin.' : ((strlen($value['description']) > 150) ? substr($value['description'], 0, 147).'...' : $value['description'])
			)));
		
		}
	
	/*
	 * Build Pagination
	 */	
	if(count($data) == 25)
		if(isset($_GET['start']) && $_GET['start'] > 24){
			$previous = $_GET['start'];
			$next = $_GET['start'];
		}else{
			$previous = false;
			$next = $_GET['start'];
		}
	else
		if(isset($_GET['start']) && $_GET['start'] != 0){
			$previous = $_GET['start'];
			$next = false;
		}else{
			$previous = false;
			$next = false;
		}

/*
 * Display Page
 */
echo $twig->render(
		'node/plugins/search.html', array(
			'pagination' => array(
				'term' => $_GET['term'],
				'previous' => $previous,
				'next' => $next
			),
			'results' => $results,
			'error' => $error,
			'footer' => array(
				'queries' => Database_Initiator::getCount(),
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
	));
?>

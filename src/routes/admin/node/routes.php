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
use \ORM, \PufferPanel\Core\Components\Functions;

$klein->respond('GET', '/admin/node', function($request, $response, $service) use($core) {

	$response->body($core->twig->render(
		'admin/node/list.html',
		array(
			'flash' => $service->flashes(),
			'nodes' => ORM::forTable('nodes')
				->select('nodes.*')->select('locations.long', 'l_location')
				->join('locations', array('nodes.location', '=', 'locations.short'))
				->findMany()
		)
	))->send();

});

$klein->respond('GET', '/admin/node/new', function($request, $response, $service) use($core) {

	$response->body($core->twig->render(
		'admin/node/new.html',
		array(
			'flash' => $service->flashes(),
			'locations' => ORM::forTable('locations')->findMany()
		)
	))->send();

});

$klein->respond('POST', '/admin/node/new', function($request, $response, $service) use($core) {

	if(
		!is_numeric($request->param('gsd_listen')) ||
		!is_numeric($request->param('gsd_console')) ||
		!is_numeric($request->param('allocate_memory')) ||
		!is_numeric($request->param('allocate_disk'))
	) {

		$service->flash('<div class="alert alert-danger">You seem to have passed some non-integers through. Try double checking the GSD listening ports as well as the disk and memory allocation.</div>');
		$response->redirect('/admin/node/new')->send();
		return;

	}

	if(!preg_match('/^([\/][\d\w.\-\/]+[\/])$/', $request->param('gsd_server_dir'))) {

		$service->flash('<div class="alert alert-danger">That seems to be an invalid directory that you passed.</div>');
		$response->redirect('/admin/node/new')->send();
		return;

	}

	if(!preg_match('/^[\w.-]{1,15}$/', $request->param('node_name'))) {

		$service->flash('<div class="alert alert-danger">The node name you passed does not meet server requirements. Names must be between 1 and 15 characters, and may not contain any special characters.</div>');
		$response->redirect('/admin/node/new')->send();
		return;

	}

	if(
		!filter_var(gethostbyname($request->param('fqdn')), FILTER_VALIDATE_IP) ||
		!filter_var($request->param('ip') , FILTER_VALIDATE_IP)
	) {

		$service->flash('<div class="alert alert-danger">The Fully Qualified Domain Name or Server IP provided were invalid.</div>');
		$response->redirect('/admin/node/new')->send();
		return;

	}

	$IPP = array();
	$IPA = array();

	$lines = explode("\r\n", str_replace(" ", "", $request->param('ip_port')));
	foreach($lines as $id => $values) {

		list($ip, $ports) = explode('|', $values);

		$IPA = array_merge($IPA, array($ip => array()));
		$IPP = array_merge($IPP, array($ip => array()));

		$portList = Functions::processPorts($ports);

		$portCount = count($portList);
		for($l=0; $l<$portCount; $l++) {
			$IPP[$ip][$portList[$l]] = 1;
		}

		if(count($IPP[$ip]) > 0) {
			$IPA[$ip] = array_merge($IPA[$ip], array("ports_free" => count($IPP[$ip])));
		} else {

			$service->flash('<div class="alert alert-danger">You must enter ports to be used with the IP.</div>');
			$response->redirect('/admin/node/new')->send();
			return;

		}

	}

	$node = ORM::forTable('nodes')->create();
	$node->set(array(
		'node' => $request->param('node_name'),
		'location' => $request->param('location'),
		'allocate_memory' => ($request->param('mem_selector') == 1) ? ($request->param('allocate_memory') * 1024) : $request->param('allocate_memory'),
		'allocate_disk' => ($request->param('disk_selector') == 1) ? ($request->param('allocate_disk') * 1024) : $request->param('allocate_disk'),
		'fqdn' => $request->param('fqdn'),
		'ip' => $request->param('ip'),
		'gsd_secret' => $core->auth->generateUniqueUUID('nodes', 'gsd_secret'),
		'gsd_listen' => $request->param('gsd_listen'),
		'gsd_console' => $request->param('gsd_console'),
		'gsd_server_dir' => $request->param('gsd_server_dir'),
		'ips' => json_encode($IPA),
		'ports' => json_encode($IPP),
		'public' => (!$request->param('is_public')) ? 0 : 1
	));
	$node->save();

	$service->flash('<div class="alert alert-success">Node successfully created. Please make sure to setup GSD properly on the node and then begin adding servers.</div>');
	$response->redirect('/admin/node/view/'.$node->id())->send();
	return;

});

$klein->respond('GET', '/admin/node/locations', function($request, $response, $service) use($core) {

	$response->body($core->twig->render(
		'admin/node/locations.html',
		array(
			'flash' => $service->flashes(),
			'locations' => ORM::forTable('locations')
				->select_many('locations.*')
				->select_expr('COUNT(nodes.id)', 'totalnodes')
				->left_outer_join('nodes', array('locations.short', '=', 'nodes.location'))
				->group_by('locations.id')
				->find_many()
		)
	))->send();

});

$klein->respond('POST', '/admin/node/locations', function($request, $response, $service) use($core) {

	if(!preg_match('/^[\w-]{1,10}$/', $request->param('shortcode'))) {

		$service->flash('<div class="alert alert-danger">Location shotcode must be between 1 and 10 characters, and not contain any special characters.</div>');
		$response->redirect('/admin/node/locations')->send();
		return;

	}

	$location = ORM::forTable('locations')->create();
	$location->short = $request->param('shortcode');
	$location->long = $request->param('location');
	$location->save();

	$service->flash('<div class="alert alert-success">Successfully added a new location.</div>');
	$response->redirect('/admin/node/locations')->send();
	return;

});

$klein->respond('GET', '/admin/node/locations/[:shortcode]/edit', function($request, $response) use($core) {

	if(!preg_match('/^[\w-]{1,10}$/', $request->param('shortcode'))) {

		$response->code(404)->body('<div class="alert alert-danger">No location provided or the location was invalid.</div>')->send();
		return;

	}

	$location = ORM::forTable('locations')->where('short', $request->param('shortcode'))->findOne();
	if(!$location) {

		$response->code(404)->body('<div class="alert alert-danger">That location could not be found in the system.</div>')->send();
		return;

	}

	$response->body($core->twig->render(
		'admin/node/location-popup.html',
		array(
			'location' => $location
		)
	))->send();

});

$klein->respond('GET', '/admin/node/locations/[:shortcode]/delete', function($request, $response, $service) use($core) {

	$location = ORM::forTable('locations')->where('short', $request->param('shortcode'))->findOne();
	if(!$location) {

		$service->flash('<div class="alert alert-danger">The requested location could not be found in the system.</div>');
		$response->redirect('/admin/node/locations')->send();
		return;

	}

	if(ORM::forTable('nodes')->where('location', $location->short)->findMany()) {

		$service->flash('<div class="alert alert-danger">You may not delete locations with currently active nodes.</div>');
		$response->redirect('/admin/node/locations')->send();
		return;

	}

	$location->delete();
	$service->flash('<div class="alert alert-success">The requested location has been deleted from the system.</div>');
	$response->redirect('/admin/node/locations')->send();

});


$klein->respond('POST', '/admin/node/locations/[:shortcode]/edit', function($request, $response, $service) use($core) {

	if(!preg_match('/^[\w-]{1,10}$/', $request->param('shortcode'))) {

		$service->flash('<div class="alert alert-danger">Location shotcode must be between 1 and 10 characters, and not contain any special characters.</div>');
		$response->redirect('/admin/node/locations')->send();
		return;

	}

	$location = ORM::forTable('locations')->where('short', $request->param('shortcode'))->findOne();
	if(!$location) {

		$service->flash('<div class="alert alert-danger">The requested location could not be found in the system.</div>');
		$response->redirect('/admin/node/locations')->send();
		return;

	}

	$nodes = ORM::forTable('nodes')->where('location', $location->short)->findMany();
	foreach($nodes as $node) {

		$node->location = $request->param('location-short');
		$node->save();

	}

	$location->short = $request->param('location-short');
	$location->long = $request->param('location-long');
	$location->save();

	$service->flash('<div class="alert alert-success">Location information has been successfully updated.</div>');
	$response->redirect('/admin/node/locations')->send();
	return;

});

$klein->respond('GET', '/admin/node/view/[i:id]', function($request, $response, $service) use ($core) {

	$node = ORM::forTable('nodes')->findOne($request->param('id'));

	if(!$node) {

		$service->flash('<div class="alert alert-danger">A node by that ID does not exist in the system.</div>');
		$response->redirect('/admin/node')->send();
		return;

	}

	$response->body($core->twig->render(
		'admin/node/view.html',
		array(
			'flash' => $service->flashes(),
			'node' => $node,
			'portlisting' => json_decode($node->ports, true),
		)
	))->send();

});

$klein->respond('POST', '/admin/node/view/[i:id]/settings', function($request, $response, $service) use($core) {

	if(!preg_match('/^[\w.-]{1,15}$/', $request->param('name'))) {

		$service->flash('<div class="alert alert-danger">The node name did not meet server requirements. Node names must be between 1 and 15 charatcers and contain no special characters.</div>');
		$response->redirect('/admin/node/view/'.$request->param('id'))->send();
		return;

	}

	if(!filter_var(gethostbyname($request->param('fqdn')), FILTER_VALIDATE_IP, FILTER_FLAG_NO_RES_RANGE)) {

		$service->flash('<div class="alert alert-danger">The node Fully Qualified Domain Name is not valid. Domains must resolve to a non-reserved IP.</div>');
		$response->redirect('/admin/node/view/'.$request->param('id'))->send();
		return;

	}


	$node = ORM::forTable('nodes')->findOne($request->param('id'));
	if($node) {

		$node->node = $request->param('name');
		$node->fqdn = $request->param('fqdn');
		$node->save();

		$service->flash('<div class="alert alert-success">Your node settings have been updated.</div>');
		$response->redirect('/admin/node/view/'.$request->param('id'))->send();

	} else {

		$service->flash('<div class="alert alert-danger">The requested node does not exist in the system.</div>');
		$response->redirect('/admin/node')->send();

	}

});

$klein->respond('POST', '/admin/node/view/[i:id]/add-port', function($request, $response, $service) use($core) {

	if(!preg_match('/^[\d-,]+$/', $request->param('add_ports'))) {

		$service->flash('<div class="alert alert-danger">Hold up there cowboy. Those ports don\'t seem to be in a valid format.</div>');
		$response->redirect('/admin/node/view/'.$request->param('id').'?tab=allocation')->send();
		return;

	}

	$portList = Functions::processPorts($request->param('add_ports'));

	$node = ORM::forTable('nodes')->findOne($request->param('add_ports_node'));

	if(!$node) {

		$service->flash('<div class="alert alert-danger">The requested node does not exist in the system.</div>');
		$response->redirect('/admin/node')->send();
		return;

	}

	$saveips = json_decode($node->ips, true);
	$saveports = json_decode($node->ports, true);

	foreach($portList as $id => $port) {

		if(
			(strlen($port) > 0 && strlen($port) < 6) &&
			array_key_exists($request->param('add_ports_ip'), $saveports) &&
			!array_key_exists($port, $saveports[$request->param('add_ports_ip')])
		) {

			$saveports[$request->param('add_ports_ip')][$port] = 1;
			$saveips[$request->param('add_ports_ip')]['ports_free']++;

		}

	}

	$node->ips = json_encode($saveips);
	$node->ports = json_encode($saveports);
	$node->save();

	$service->flash('<div class="alert alert-success">Ports were successfully added to <strong>'.$request->param('add_ports_ip').'</strong>.</div>');
	$response->redirect('/admin/node/view/'.$request->param('id').'?tab=allocation')->send();

});

$klein->respond('POST', '/admin/node/view/[i:id]/add-ip', function($request, $response, $service) use($core) {

	$lines = explode("\r\n", str_replace(" ", "", $request->param('ip_port')));

	$node = ORM::forTable('nodes')->findOne($request->param('id'));

	if(!$node) {

		$service->flash('<div class="alert alert-danger">The requested node does not exist in the system.</div>');
		$response->redirect('/admin/node')->send();
		return;

	}

	$IPA = json_decode($node->ips, true);
	$IPP = json_decode($node->ports, true);

	foreach($lines as $id => $values) {

		list($ip, $ports) = explode('|', $values);

		$IPA = array_merge($IPA, array($ip => array()));
		$IPP = array_merge($IPP, array($ip => array()));

		$portList = Functions::processPorts($ports);

		$portCount = count($portList);
		for($l=0; $l<$portCount; $l++) {

				$IPP[$ip][$portList[$l]] = 1;

		}

		if(count($IPP[$ip]) > 0) {

			$IPA[$ip] = array_merge($IPA[$ip], array("ports_free" => count($IPP[$ip])));

		} else {

			$service->flash('<div class="alert alert-danger">You must enter ports to be used with the IP.</div>');
			$response->redirect('/admin/node/view/'.$request->param('id').'?tab-allocation')->send();
			return;

		}

	}

	$node->ips = json_encode($IPA);
	$node->ports = json_encode($IPP);
	$node->save();

	$service->flash('<div class="alert alert-success">New IP address has been successfully allocated to this node.</div>');
	$response->redirect('/admin/node/view/'.$request->param('id').'?tab=allocation')->send();

});

$klein->respond('POST', '/admin/node/view/[i:id]/ftp', function($request, $response, $service) use($core) {

	if(!filter_var($request->param('ip') , FILTER_VALIDATE_IP, FILTER_FLAG_NO_RES_RANGE)) {

		$service->flash('<div class="alert alert-danger">The IP provided for FTP was not valid. FTP IPs must resolve to a non-reserved address.</div>');
		$response->redirect('/admin/node/view/'.$request->param('id').'?tab=ftp')->send();
		return;

	}

	$node = ORM::forTable('nodes')->findOne($request->param('id'));

	if($node) {

		$node->ip = $request->param('ip');
		$node->save();

		$service->flash('<div class="alert alert-success">Your node settings have been updated.</div>');
		$response->redirect('/admin/node/view/'.$request->param('id').'?tab=ftp')->send();

	} else {

		$service->flash('<div class="alert alert-danger">The requested node does not exist in the system.</div>');
		$response->redirect('/admin/node')->send();

	}


});

$klein->respond('POST', '/admin/node/view/[i:id]/reset-token', function($request, $response) use($core) {

	$node = ORM::forTable('nodes')->findOne($request->param('id'));
	if(!$core->gsd->avaliable($node->ip, $node->gsd_listen)) {

		$node->gsd_secret = $core->auth->generateUniqueUUID('nodes', 'gsd_secret');
		$node->save();

		$response->body($node->gsd_secret)->send();

	} else {

		$response->body("GSD Must be stopped before running this command.")->send();

	}

});

$klein->respond('POST', '/admin/node/view/[i:id]/delete-port', function($request, $response) use($core) {

	$node = ORM::forTable('nodes')->findOne($request->param('node'));

	if(!$node) {
		$response->body("The requested node does not exist in the system.")->send();
		return;
	}

	$ips = json_decode($node->ips, true);
	$ports = json_decode($node->ports, true);

	if(
		array_key_exists($request->param('ip'), $ports) &&
		array_key_exists($request->param('port'), $ports[$request->param('ip')]) &&
		$ports[$request->param('ip')][$request->param('port')] == 1
	) {

		unset($ports[$request->param('ip')][$request->param('port')]);
		$ips[$request->param('ip')]['ports_free']--;

	} else {
		$response->body("That port either doesn't exist or is currently in use.")->send();
		return;
	}

	$node->ips = json_encode($ips);
	$node->ports = json_encode($ports);
	$node->save();

	$response->body('Done')->send();

});
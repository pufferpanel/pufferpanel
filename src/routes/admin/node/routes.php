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

use \ORM;

$klein->respond('GET', BASE_URL.'/admin/node', function($request, $response, $service) use($core) {

    $response->body($core->twig->render('admin/node/list.html', array(
            'flash' => $service->flashes(),
            'nodes' => ORM::forTable('nodes')
                        ->select('nodes.*')->select('locations.long', 'l_location')
                        ->join('locations', array('nodes.location', '=', 'locations.id'))
                        ->findMany()
        )
    ));

});

$klein->respond('GET', BASE_URL.'/admin/node/new', function($request, $response, $service) use($core) {

    if (ORM::forTable('locations')->count() == 0) {
        $service->flash('<div class="alert alert-danger">You must have at least one location defined before creating a node.</div>');
        $response->redirect('/admin/node/locations');
    } else {
        $response->body($core->twig->render('admin/node/new.html', array(
            'flash' => $service->flashes(),
            'locations' => ORM::forTable('locations')->findMany()
        )));
    }

});

$klein->respond('POST', BASE_URL.'/admin/node/new', function($request, $response, $service) use($core) {

    if (!is_numeric($request->param('daemon_listen'))) {

        $service->flash('<div class="alert alert-danger">You seem to have passed some non-integers through. Try double checking the daemon listening ports as well as the disk and memory allocation.</div>');
        $response->redirect('/admin/node/new');
        return;
    }

    if (!preg_match('/^[\w.-]{1,15}$/', $request->param('node_name'))) {

        $service->flash('<div class="alert alert-danger">The node name you passed does not meet server requirements. Names must be between 1 and 15 characters, and may not contain any special characters.</div>');
        $response->redirect('/admin/node/new');
        return;
    }

    if (!filter_var(gethostbyname($request->param('fqdn')), FILTER_VALIDATE_IP)) {

        $service->flash('<div class="alert alert-danger">The Fully Qualified Domain Name or Server IP provided were invalid.</div>');
        $response->redirect('/admin/node/new');
        return;
    }

    if (ORM::forTable('nodes')->where_any_is(array(
                array('name' => $request->param('node_name')),
                array('fqdn' => $request->param('fqdn'))
            ))->findOne()) {

        $service->flash('<div class="alert alert-danger">A node with that name or IP address is already in use on this system.</div>');
        $response->redirect('/admin/node/new');
        return;
    }
    
    $internalip = $request->param('internalip');
    if (trim($internalip) === '') {
        $internalip = $request->param('fqdn');
    }

    $node = ORM::forTable('nodes')->create();
    $node->set(array(
        'name' => $request->param('node_name'),
        'location' => $request->param('location'),
        'allocate_memory' => 0,
        'allocate_disk' => 0,
        'fqdn' => $request->param('fqdn'),
        'ip' => $internalip,
        'daemon_secret' => $core->auth->generateUniqueUUID('nodes', 'daemon_secret'),
        'daemon_listen' => $request->param('daemon_listen'),
        'daemon_sftp' => '5657',
        'ips' => '{}',
        'ports' => '{}',
        'public' => 1,
        'docker' => 0
    ));
    $node->save();

    $service->flash('<div class="alert alert-success">Node successfully created. Please make sure to setup the daemon properly on the node and then begin adding servers.</div>');
    $response->redirect('/admin/node/view/' . $node->id());

});

$klein->respond('GET', BASE_URL.'/admin/node/locations', function($request, $response, $service) use($core) {

    $response->body($core->twig->render('admin/node/locations.html', array(
        'flash' => $service->flashes(),
        'locations' => ORM::forTable('locations')
                ->select_many('locations.*')
                ->select_expr('COUNT(DISTINCT nodes.id)', 'totalnodes')
                ->select_expr('COUNT(servers.id)', 'totalservers')
                ->left_outer_join('nodes', array('locations.id', '=', 'nodes.location'))
                ->left_outer_join('servers', array('servers.node', '=', 'nodes.id'))
                ->group_by('locations.id')
                ->find_many()
    )));
});

$klein->respond('POST', BASE_URL.'/admin/node/locations', function($request, $response, $service) use($core) {

    if (!preg_match('/^[\w-]{1,10}$/', $request->param('shortcode'))) {

        $service->flash('<div class="alert alert-danger">Location shotcode must be between 1 and 10 characters, and not contain any special characters.</div>');
        $response->redirect('/admin/node/locations')->send();
        return;
    }

    $location = ORM::forTable('locations')->create();
    $location->short = $request->param('shortcode');
    $location->long = $request->param('location');
    $location->save();

    $service->flash('<div class="alert alert-success">Successfully added a new location.</div>');
    $response->redirect('/admin/node/locations');

});

$klein->respond('GET', BASE_URL.'/admin/node/locations/[i:id]/edit', function($request, $response) use($core) {

    $location = ORM::forTable('locations')->findOne((int) $request->param('id'));
    if (!$location) {

        $response->code(404)->body('<div class="alert alert-danger">That location could not be found in the system.</div>')->send();
        return;
    }

    $response->body($core->twig->render(
        'admin/node/location-popup.html',
        [ 'location' => $location ]
    ))->send();
});

$klein->respond('GET', BASE_URL.'/admin/node/locations/[i:id]/delete', function($request, $response, $service) use($core) {

    $location = ORM::forTable('locations')->findOne((int) $request->param('id'));
    if (!$location) {

        $service->flash('<div class="alert alert-danger">The requested location could not be found in the system.</div>');
        $response->redirect('/admin/node/locations')->send();
        return;
    }

    if (ORM::forTable('nodes')->where('location', $location->short)->findMany()) {

        $service->flash('<div class="alert alert-danger">You may not delete locations with currently active nodes.</div>');
        $response->redirect('/admin/node/locations')->send();
        return;
    }

    $location->delete();
    $service->flash('<div class="alert alert-success">The requested location has been deleted from the system.</div>');
    $response->redirect('/admin/node/locations')->send();
});


$klein->respond('POST', BASE_URL.'/admin/node/locations/[i:id]/edit', function($request, $response, $service) use($core) {

    $location = ORM::forTable('locations')->findOne((int) $request->param('id'));
    if (!$location) {

        $service->flash('<div class="alert alert-danger">The requested location could not be found in the system.</div>');
        $response->redirect('/admin/node/locations')->send();
        return;
    }

    if (!preg_match('/^[\w-]{1,10}$/', $request->param('location-short'))) {

        $service->flash('<div class="alert alert-danger">Location shortcode must be between 1 and 10 characters, and not contain any special characters.</div>');
        $response->redirect('/admin/node/locations')->send();
        return;
    }

    $location->short = $request->param('location-short');
    $location->long = $request->param('location-long');
    $location->save();

    $service->flash('<div class="alert alert-success">Location information has been successfully updated.</div>');
    $response->redirect('/admin/node/locations')->send();
    return;
});

$klein->respond('GET', BASE_URL.'/admin/node/view/[i:id]', function($request, $response, $service) use ($core) {

    $node = ORM::forTable('nodes')->findOne($request->param('id'));

    if (!$node) {

        $service->flash('<div class="alert alert-danger">A node by that ID does not exist in the system.</div>');
        $response->redirect('/admin/node')->send();
        return;
    }

    $response->body($core->twig->render(
                    'admin/node/view.html', array(
                'flash' => $service->flashes(),
                'node' => $node,
                'locations' => ORM::forTable('locations')->findMany(),
                'portlisting' => json_decode($node->ports, true),
                'autodeploy' => ORM::forTable('autodeploy')->where('node', $request->param('id'))->where_gt('expires', time())->findOne()
                    )
    ))->send();
});

$klein->respond('POST', BASE_URL.'/admin/node/view/[i:id]/settings', function($request, $response, $service) use($core) {

    if (!preg_match('/^[\w.-]{1,15}$/', $request->param('name'))) {

        $service->flash('<div class="alert alert-danger">The node name did not meet server requirements. Node names must be between 1 and 15 charatcers and contain no special characters.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id'))->send();
        return;
    }

    if (!filter_var(gethostbyname($request->param('fqdn')), FILTER_VALIDATE_IP, FILTER_FLAG_NO_RES_RANGE) && !filter_var($request->param('fqdn'), FILTER_VALIDATE_IP, FILTER_FLAG_NO_RES_RANGE)) {

        $service->flash('<div class="alert alert-danger">The node Fully Qualified Domain Name is not valid. Domains must resolve to a non-reserved IP.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id'))->send();
        return;
    }
    
    $internalip = $request->param('internalip');
    if (trim($internalip) === '') {
        $internalip = $request->param('fqdn');
    }

    if (!filter_var(gethostbyname($internalip), FILTER_VALIDATE_IP) && !filter_var($internalip, FILTER_VALIDATE_IP)) {

        $service->flash('<div class="alert alert-danger">The node\'s internal IP is not valid. Domains must resolve to an IP.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id'))->send();
        return;
    }


    $location = ORM::forTable('locations')->select('id')->findOne($request->param('location'));
    $node = ORM::forTable('nodes')->findOne($request->param('id'));
    if ($node && $location) {

        $node->name = $request->param('name');
        $node->location = $location->id;
        $node->fqdn = $request->param('fqdn');
        $node->ip = $internalip;
        $node->save();

        $service->flash('<div class="alert alert-success">Your node settings have been updated.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id'))->send();
    } else {

        $service->flash('<div class="alert alert-danger">The requested node or location does not exist in the system.</div>');
        $response->redirect('/admin/node')->send();
    }
});

$klein->respond('POST', BASE_URL.'/admin/node/view/[i:id]/generate-autodeploy', function($request, $response, $service) use($core) {

    $node = ORM::forTable('nodes')->findOne($request->param('id'));

    if (ORM::forTable('autodeploy')->where('node', $node->id)->where_gt('expires', time())->count() > 0) {
        $service->flash('<div class="alert alert-danger">There is already an exisiting auto-deploy script for this server.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id') . '?tab=deploy')->send();
        return;
    }

    $deploy = ORM::forTable('autodeploy')->create();
    $deploy->set(array(
        'node' => $node->id,
        'code' => $core->auth->generateUniqueUUID('autodeploy', 'code'),
        'expires' => time() + 900
    ));
    $deploy->save();

    $service->flash('<div class="alert alert-success">A new deployment script endpoint was successfully configured for this node.</div>');
    $response->redirect('/admin/node/view/' . $request->param('id') . '?tab=deploy')->send();
    return;
});

$klein->respond('POST', BASE_URL.'/admin/node/view/[i:id]/delete', function($request, $response, $service) use($core) {

    $node = ORM::forTable('nodes')->findOne($request->param('id'));

    if (ORM::forTable('servers')->where('node', $node->id)->count() > 0) {
        $service->flash('<div class="alert alert-danger">In order to delete this node there cannot be servers associated with it.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id') . '?tab=delete')->send();
        return;
    }

    $node->delete();

    $service->flash('<div class="alert alert-success">The node (' . $node->name . ') was successfully deleted.</div>');
    $response->redirect('/admin/node')->send();
    return;
});

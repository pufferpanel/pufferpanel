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

use \ORM,
    \PufferPanel\Core\Components\Functions;

$klein->respond('GET', '/admin/node', function($request, $response, $service) use($core) {

    $response->body($core->twig->render(
                'admin/node/list.html', array(
              'flash' => $service->flashes(),
              'nodes' => ORM::forTable('nodes')
                    ->select('nodes.*')->select('locations.long', 'l_location')
                    ->join('locations', array('nodes.location', '=', 'locations.id'))
                    ->findMany()
                )
    ))->send();
});

$klein->respond('GET', '/admin/node/new', function($request, $response, $service) use($core) {

    if (ORM::forTable('locations')->count() == 0) {

        $service->flash('<div class="alert alert-danger">You must have at least one location defined before creating a node.</div>');
        $response->redirect('/admin/node/locations')->send();
        return;
    }

    $response->body($core->twig->render(
                'admin/node/new.html', array(
              'flash' => $service->flashes(),
              'locations' => ORM::forTable('locations')->findMany()
                )
    ))->send();
});

$klein->respond('POST', '/admin/node/new', function($request, $response, $service) use($core) {

    if (
          !is_numeric($request->param('daemon_listen')) ||
          !is_numeric($request->param('daemon_sftp'))
    ) {

        $service->flash('<div class="alert alert-danger">You seem to have passed some non-integers through. Try double checking the daemon listening ports as well as the disk and memory allocation.</div>');
        $response->redirect('/admin/node/new')->send();
        return;
    }

    if (!preg_match('/^([\/][\d\w.\-\/]+[\/])$/', $request->param('daemon_base_dir'))) {

        $service->flash('<div class="alert alert-danger">That seems to be an invalid directory that you passed.</div>');
        $response->redirect('/admin/node/new')->send();
        return;
    }

    if (!preg_match('/^[\w.-]{1,15}$/', $request->param('node_name'))) {

        $service->flash('<div class="alert alert-danger">The node name you passed does not meet server requirements. Names must be between 1 and 15 characters, and may not contain any special characters.</div>');
        $response->redirect('/admin/node/new')->send();
        return;
    }

    if (!filter_var(gethostbyname($request->param('fqdn')), FILTER_VALIDATE_IP)) {

        $service->flash('<div class="alert alert-danger">The Fully Qualified Domain Name or Server IP provided were invalid.</div>');
        $response->redirect('/admin/node/new')->send();
        return;
    }

    if (ORM::forTable('nodes')->where_any_is(array(
              array('name' => $request->param('node_name')),
              array('fqdn' => $request->param('fqdn'))
          ))->findOne()) {

        $service->flash('<div class="alert alert-danger">A node with that name or IP address is already in use on this system.</div>');
        $response->redirect('/admin/node/new')->send();
        return;
    }

    $IPP = array();
    $IPA = array();

    if (!$request->param('ip_port') || empty($request->param('ip_port'))) {

        $service->flash('<div class="alert alert-danger">You must define at least one IP and Port for this node to add servers to.</div>');
        $response->redirect('/admin/node/new')->send();
        return;
    }

    $lines = explode("\r\n", str_replace(" ", "", $request->param('ip_port')));
    foreach ($lines as $id => $values) {

        list($ip, $ports) = explode('|', $values);

        if (!trim($ip)) {

            $service->flash('<div class="alert alert-danger">An IP must be specified with a port list.</div>');
            $response->redirect('/admin/node/new')->send();
            return;
        }

        if (!trim($ports)) {

            $service->flash('<div class="alert alert-danger">You must provide at least one port with an IP</div>');
            $response->redirect('/admin/node/new')->send();
            return;
        }

        $IPA = array_merge($IPA, array($ip => array()));
        $IPP = array_merge($IPP, array($ip => array()));

        try {
            $portList = Functions::processPorts($ports);
        } catch (Exception $ex) {
            $service->flash('<div class="alert alert-danger">' . $ex . getMessage() . '</div>');
            $response->redirect('/admin/node/new')->send();
            return;
        }

        $portCount = count($portList);
        for ($l = 0; $l < $portCount; $l++) {
            $IPP[$ip][$portList[$l]] = 1;
        }

        if (count($IPP[$ip]) > 0) {
            $IPA[$ip] = array_merge($IPA[$ip], array("ports_free" => count($IPP[$ip])));
        } else {

            $service->flash('<div class="alert alert-danger">You must enter ports to be used with the IP.</div>');
            $response->redirect('/admin/node/new')->send();
            return;
        }
    }

    $node = ORM::forTable('nodes')->create();
    $node->set(array(
        'name' => $request->param('node_name'),
        'location' => $request->param('location'),
        'allocate_memory' => 0,
        'allocate_disk' => 0,
        'fqdn' => $request->param('fqdn'),
        'ip' => $request->param('fqdn'),
        'daemon_secret' => $core->auth->generateUniqueUUID('nodes', 'daemon_secret'),
        'daemon_listen' => $request->param('daemon_listen'),
        'daemon_sftp' => $request->param('daemon_sftp'),
        'daemon_base_dir' => $request->param('daemon_base_dir'),
        'ips' => json_encode($IPA),
        'ports' => json_encode($IPP),
        'public' => 1,
        'docker' => ($request->param('is_docker')) ? 1 : 0
    ));
    $node->save();

    $service->flash('<div class="alert alert-success">Node successfully created. Please make sure to setup the daemon properly on the node and then begin adding servers.</div>');
    $response->redirect('/admin/node/view/' . $node->id())->send();
    return;
});

$klein->respond('GET', '/admin/node/locations', function($request, $response, $service) use($core) {

    $response->body($core->twig->render(
                'admin/node/locations.html', array(
              'flash' => $service->flashes(),
              'locations' => ORM::forTable('locations')
                    ->select_many('locations.*')
                    ->select_expr('COUNT(DISTINCT nodes.id)', 'totalnodes')
                    ->select_expr('COUNT(servers.id)', 'totalservers')
                    ->left_outer_join('nodes', array('locations.id', '=', 'nodes.location'))
                    ->left_outer_join('servers', array('servers.node', '=', 'nodes.id'))
                    ->group_by('locations.id')
                    ->find_many()
                )
    ))->send();
});

$klein->respond('POST', '/admin/node/locations', function($request, $response, $service) use($core) {

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
    $response->redirect('/admin/node/locations')->send();
    return;
});

$klein->respond('GET', '/admin/node/locations/[:shortcode]/edit', function($request, $response) use($core) {

    if (!preg_match('/^[\w-]{1,10}$/', $request->param('shortcode'))) {

        $response->code(404)->body('<div class="alert alert-danger">No location provided or the location was invalid.</div>')->send();
        return;
    }

    $location = ORM::forTable('locations')->where('short', $request->param('shortcode'))->findOne();
    if (!$location) {

        $response->code(404)->body('<div class="alert alert-danger">That location could not be found in the system.</div>')->send();
        return;
    }

    $response->body($core->twig->render(
                'admin/node/location-popup.html', array(
              'location' => $location
                )
    ))->send();
});

$klein->respond('GET', '/admin/node/locations/[:shortcode]/delete', function($request, $response, $service) use($core) {

    $location = ORM::forTable('locations')->where('short', $request->param('shortcode'))->findOne();
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


$klein->respond('POST', '/admin/node/locations/[:shortcode]/edit', function($request, $response, $service) use($core) {

    if (!preg_match('/^[\w-]{1,10}$/', $request->param('location-short'))) {

        $service->flash('<div class="alert alert-danger">Location shortcode must be between 1 and 10 characters, and not contain any special characters.</div>');
        $response->redirect('/admin/node/locations')->send();
        return;
    }

    $location = ORM::forTable('locations')->where('short', $request->param('location-short'))->findOne();
    if (!$location) {

        $service->flash('<div class="alert alert-danger">The requested location could not be found in the system.</div>');
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

$klein->respond('GET', '/admin/node/view/[i:id]', function($request, $response, $service) use ($core) {

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

$klein->respond('POST', '/admin/node/view/[i:id]/settings', function($request, $response, $service) use($core) {

    if (!preg_match('/^[\w.-]{1,15}$/', $request->param('name'))) {

        $service->flash('<div class="alert alert-danger">The node name did not meet server requirements. Node names must be between 1 and 15 charatcers and contain no special characters.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id'))->send();
        return;
    }

    if (!filter_var(gethostbyname($request->param('fqdn')), FILTER_VALIDATE_IP, FILTER_FLAG_NO_RES_RANGE)) {

        $service->flash('<div class="alert alert-danger">The node Fully Qualified Domain Name is not valid. Domains must resolve to a non-reserved IP.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id'))->send();
        return;
    }


    $location = ORM::forTable('locations')->select('id')->findOne($request->param('location'));
    $node = ORM::forTable('nodes')->findOne($request->param('id'));
    if ($node && $location) {

        $node->name = $request->param('name');
        $node->location = $location->id;
        $node->fqdn = $request->param('fqdn');
        $node->save();

        $service->flash('<div class="alert alert-success">Your node settings have been updated.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id'))->send();
    } else {

        $service->flash('<div class="alert alert-danger">The requested node or location does not exist in the system.</div>');
        $response->redirect('/admin/node')->send();
    }
});

$klein->respond('POST', '/admin/node/view/[i:id]/add-port', function($request, $response, $service) use($core) {

    if (!preg_match('/^[\d-,]+$/', $request->param('add_ports'))) {

        $service->flash('<div class="alert alert-danger">Hold up there cowboy. Those ports don\'t seem to be in a valid format.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id') . '?tab=allocation')->send();
        return;
    }

    $portList = Functions::processPorts($request->param('add_ports'));

    $node = ORM::forTable('nodes')->findOne($request->param('add_ports_node'));

    if (!$node) {

        $service->flash('<div class="alert alert-danger">The requested node does not exist in the system.</div>');
        $response->redirect('/admin/node')->send();
        return;
    }

    $saveips = json_decode($node->ips, true);
    $saveports = json_decode($node->ports, true);

    foreach ($portList as $id => $port) {

        if (
              (strlen($port) > 0 && strlen($port) < 6) &&
              array_key_exists($request->param('add_ports_ip'), $saveports) &&
              !array_key_exists($port, $saveports[$request->param('add_ports_ip')])
        ) {

            $saveports[$request->param('add_ports_ip')][$port] = 1;
            $saveips[$request->param('add_ports_ip')]['ports_free'] ++;
        }
    }

    $node->ips = json_encode($saveips);
    $node->ports = json_encode($saveports);
    $node->save();

    $service->flash('<div class="alert alert-success">Ports were successfully added to <strong>' . $request->param('add_ports_ip') . '</strong>.</div>');
    $response->redirect('/admin/node/view/' . $request->param('id') . '?tab=allocation')->send();
});

$klein->respond('POST', '/admin/node/view/[i:id]/add-ip', function($request, $response, $service) use($core) {

    $lines = explode("\r\n", str_replace(" ", "", $request->param('ip_port')));

    $node = ORM::forTable('nodes')->findOne($request->param('id'));

    if (!$node) {

        $service->flash('<div class="alert alert-danger">The requested node does not exist in the system.</div>');
        $response->redirect('/admin/node')->send();
        return;
    }

    $IPA = json_decode($node->ips, true);
    $IPP = json_decode($node->ports, true);

    foreach ($lines as $id => $values) {

        list($ip, $ports) = explode('|', $values);

        $IPA = array_merge($IPA, array($ip => array()));
        $IPP = array_merge($IPP, array($ip => array()));

        $portList = Functions::processPorts($ports);

        $portCount = count($portList);
        for ($l = 0; $l < $portCount; $l++) {

            $IPP[$ip][$portList[$l]] = 1;
        }

        if (count($IPP[$ip]) > 0) {

            $IPA[$ip] = array_merge($IPA[$ip], array("ports_free" => count($IPP[$ip])));
        } else {

            $service->flash('<div class="alert alert-danger">You must enter ports to be used with the IP.</div>');
            $response->redirect('/admin/node/view/' . $request->param('id') . '?tab-allocation')->send();
            return;
        }
    }

    $node->ips = json_encode($IPA);
    $node->ports = json_encode($IPP);
    $node->save();

    $service->flash('<div class="alert alert-success">New IP address has been successfully allocated to this node.</div>');
    $response->redirect('/admin/node/view/' . $request->param('id') . '?tab=allocation')->send();
});

$klein->respond('POST', '/admin/node/view/[i:id]/sftp', function($request, $response, $service) use($core) {

    if (!filter_var($request->param('ip'), FILTER_VALIDATE_IP, FILTER_FLAG_NO_RES_RANGE)) {

        $service->flash('<div class="alert alert-danger">The IP provided for SFTP was not valid. SFTP IPs must resolve to a non-reserved address.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id') . '?tab=sftp')->send();
        return;
    }

    $node = ORM::forTable('nodes')->findOne($request->param('id'));

    if ($node) {

        $node->ip = $request->param('ip');
        $node->save();

        $service->flash('<div class="alert alert-success">Your node settings have been updated.</div>');
        $response->redirect('/admin/node/view/' . $request->param('id') . '?tab=sftp')->send();
    } else {

        $service->flash('<div class="alert alert-danger">The requested node does not exist in the system.</div>');
        $response->redirect('/admin/node')->send();
    }
});

$klein->respond('POST', '/admin/node/view/[i:id]/reset-token', function($request, $response) use($core) {

    $node = ORM::forTable('nodes')->findOne($request->param('id'));
    if (!$core->daemon->avaliable($node->ip, $node->daemon_listen)) {

        $node->daemon_secret = $core->auth->generateUniqueUUID('nodes', 'daemon_secret');
        $node->save();

        $response->body($node->daemon_secret)->send();
    } else {

        $response->body("GSD Must be stopped before running this command.")->send();
    }
});

$klein->respond('POST', '/admin/node/view/[i:id]/delete-port', function($request, $response) use($core) {

    $node = ORM::forTable('nodes')->findOne($request->param('id'));

    if (!$node) {
        $response->body("The requested node does not exist in the system.")->send();
        return;
    }

    $ips = json_decode($node->ips, true);
    $ports = json_decode($node->ports, true);

    if (
          array_key_exists($request->param('ip'), $ports) &&
          array_key_exists($request->param('port'), $ports[$request->param('ip')]) &&
          $ports[$request->param('ip')][$request->param('port')] == 1
    ) {

        unset($ports[$request->param('ip')][$request->param('port')]);
        $ips[$request->param('ip')]['ports_free'] --;
    } else {
        $response->body("That port either doesn't exist or is currently in use.")->send();
        return;
    }

    $node->ips = json_encode($ips);
    $node->ports = json_encode($ports);
    $node->save();

    $response->body('Done')->send();
});

$klein->respond('DELETE', '/admin/node/view/[i:id]/delete-ip/[:ip_address]', function($request, $response) use($core) {

    $node = ORM::forTable('nodes')->findOne($request->param('id'));

    if (ORM::forTable('servers')->where(array(
              'node' => $request->param('id'),
              'server_ip' => $request->param('ip_address')
          ))->count() > 0) {
        $response->code(500)->body('Unable to delete this IP because there are still servers associated with it.')->send();
        return;
    }

    $ips = json_decode($node->ips, true);
    $ports = json_decode($node->ports, true);

    if (!array_key_exists($request->param('ip_address'), $ips) || !array_key_exists($request->param('ip_address'), $ports)) {
        $response->code(500)->body('The requested IP does not exist on this node.')->send();
        return;
    }

    unset($ips[$request->param('ip_address')]);
    unset($ports[$request->param('ip_address')]);

    $node->ips = json_encode($ips);
    $node->ports = json_encode($ports);
    $node->save();

    $response->code(204)->send();
    return;
});

$klein->respond('POST', '/admin/node/view/[i:id]/generate-autodeploy', function($request, $response, $service) use($core) {

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

$klein->respond('POST', '/admin/node/view/[i:id]/delete', function($request, $response, $service) use($core) {

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

$klein->respond('GET', '/admin/node/plugins', function($request, $response, $service) use($core) {

    $response->body($core->twig->render(
                'admin/node/plugins/index.html', array(
              'flash' => $service->flashes(),
              'plugins' => ORM::forTable('plugins')->findMany(),
                )
    ))->send();
});

$klein->respond('GET', '/admin/node/plugins/view/[:hash]', function($request, $response, $service) use($core) {

    $orm = ORM::forTable('plugins')->where('hash', $request->param('hash'))->findOne();

    if (!$orm) {

        $service->flash('<div class="alert alert-danger">The requested plugin does not seem to exist in the system.</div>');
        $response->redirect('/admin/node/plugins')->send();
        return;
    }

    $response->body($core->twig->render(
                'admin/node/plugins/view.html', array(
              'flash' => $service->flashes(),
              'plugin' => $orm,
              'vars' => json_decode($orm->variables, true),
              'servers' => ORM::forTable('servers')->select('servers.*')->select('nodes.name', 'node_name')
                    ->join('nodes', array('servers.node', '=', 'nodes.id'))
                    ->where('servers.plugin', $orm->id)
                    ->findArray()
                )
    ))->send();
});

$klein->respond('POST', '/admin/node/plugins/view/[:hash]/delete', function($request, $response, $service) use($core) {

    $plugin = ORM::forTable('plugins')->where('hash', $request->param('hash'))->findOne();

    if (ORM::forTable('servers')->where('plugin', $plugin->id)->count() > 0) {

        $service->flash('<div class="alert alert-danger">The selected plugin cannot be deleted because there are servers running on it.</div>');
        $response->redirect('/admin/node/plugins/view/' . $request->param('hash'))->send();
        return;
    }

    $plugin->delete();
    $service->flash('<div class="alert alert-success">The plugin (' . $plugin->name . ') has been successfully deleted.</div>');
    $response->redirect('/admin/node/plugins')->send();
    return;
});

$klein->respond('GET', '/admin/node/plugins/new', function($request, $response, $service) use($core) {

    $response->body($core->twig->render(
                'admin/node/plugins/new.html', array(
              'flash' => $service->flashes()
                )
    ))->send();
});

$klein->respond('POST', '/admin/node/plugins/new', function($request, $response, $service) use($core) {

    if (!preg_match('/^(.{1,100})$/', $request->param('name')) || !preg_match('/^[\w.-]{1,100}$/', $request->param('slug'))) {

        $service->flash('<div class="alert alert-danger">The name or slug provided for this plugin must not exceede 100 characters.</div>');
        $response->redirect('/admin/node/plugins/new')->send();
        return;
    }

    if (ORM::forTable('plugins')->where('slug', $request->param('slug'))->findOne()) {

        $service->flash('<div class="alert alert-danger">A plugin with that slug already exists in the system.</div>');
        $response->redirect('/admin/node/plugins/new')->send();
        return;
    }

    if (!$request->param('variables_name') || empty($request->param('variables_name')[0])) {
        $built_variables = array();
    } else {

        $count = count($request->param('variables_name'));
        $built_variables = [];

        for ($i = 0; $i < $count; $i++) {

            $built_variables[$request->param('variables_identifier')[$i]] = array(
                "name" => $request->param('variables_name')[$i],
                "description" => $request->param('variables_description')[$i],
                "required" => ($request->param('variables_required')[$i] == "true") ? true : false,
                "editable" => ($request->param('variables_editable')[$i] == "true") ? true : false,
                "default" => $request->param('variables_default')[$i]
            );
        }
    }

    $new = ORM::forTable('plugins')->create();
    $hash = $core->auth->generateUniqueUUID('plugins', 'hash');
    $new->set(array(
        'hash' => $hash,
        'name' => $request->param('name'),
        'description' => $request->param('description'),
        'slug' => $request->param('slug'),
        'default_startup' => $request->param('default_startup'),
        'variables' => json_encode($built_variables)
    ));
    $new->save();

    $response->redirect('/admin/node/plugins/view/' . $hash)->send();
});

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

$klein->respond('GET', '/admin/account', function($request, $response, $service) use ($core) {

    $users = ORM::forTable('users')->findArray();

    $response->body($core->twig->render(
        'admin/account/find.html',
        array(
            'users' => $users
        )
    ))->send();

});

$klein->respond('GET', '/admin/account/new', function($request, $response, $service) use ($core) {

    $response->body($core->twig->render('admin/account/new.html'))->send();

});

$klein->respond('GET', '/admin/account/view/[i:id]', function($request, $response, $service) use ($core) {

    $core->user->rebuildData($request->param('id'));

    $date1 = new \DateTime(date('Y-m-d', $core->user->getData('register_time')));
    $date2 = new \DateTime(date('Y-m-d', time()));

    $user = $core->user->getData();
    $user['register_time'] = date('F j, Y g:ia', $core->user->getData('register_time')).' ('.$date2->diff($date1)->format("%a Days Ago").')';

    /*
     * Select Servers Owned by the User
     */
    $servers = ORM::forTable('servers')->select('servers.*')->select('nodes.node', 'node_name')
        ->join('nodes', array('servers.node', '=', 'nodes.id'))
        ->where(array('servers.owner_id' => $core->user->getData('id'), 'servers.active' => 1))
        ->findArray();

    $response->body($core->twig->render(
        'admin/account/view.html',
        array(
            'user' => $user,
            'servers' => $servers
        )
    ))->send();

});
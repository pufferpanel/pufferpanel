<?php

/*
  PufferPanel - A Minecraft Server Management Panel
  Copyright (c) 2014 PufferPanel

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

session_start();

require_once('../src/core/core.php');

$klein = new \Klein();

$klein->respond(function($request, $sresponse, $service, $app) {
    $app->register('core', function() {
        return $core;
    });
    $app->register('twig', function() {
        return $twig;
    });
    $app->register('startTime', function() {
        return $startTime;
    });
});

$klein->with('/account', 'account.php');

$klein->dispatch();

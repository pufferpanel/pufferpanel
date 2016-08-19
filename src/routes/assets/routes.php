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

$klein->respond('GET', '/assets/[:type]/[*:file]', function($request, $response) {

    $file = APP_DIR . 'assets/' . $request->param('type') . '/' . $request->param('file');

    if (!file_exists($file)) {

        $response->code(404)->body("The requested asset does not exist on this server.")->send();
    }

    $etag = md5_file($file);

    if ($request->server()["HTTP_IF_MODIFIED_SINCE"] && $request->server()["HTTP_IF_NONE_MATCH"]) {

        if (filemtime($file) <= strtotime($request->server()["HTTP_IF_MODIFIED_SINCE"]) || trim($_SERVER['HTTP_IF_NONE_MATCH']) == $etag) {
            $response->code(304);
        }
    }

    if ($request->param('type') == 'css') {
        header('Content-Type: text/' . $request->param('type'));
    } else if ($request->param('type') == 'javascript') {
        header('Content-Type: application/' . $request->param('type'));
    }

    header('X-Content-Type-Options: nosniff');
    header('Cache-control: public');
    header('Pragma: cache');
    header('Etag: "' . $etag . '"');
    header('Expires: ' . gmdate('D, d M Y H:i:s', time() + 60 * 60) . ' GMT');
    header('Last-Modified: ' . gmdate('D, d M Y H:i:s', filemtime($file)) . ' GMT');

    $response->body(file_get_contents($file))->send();
});

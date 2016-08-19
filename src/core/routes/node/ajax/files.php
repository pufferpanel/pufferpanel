<?php

/*
  PufferPanel - A Game Server Management Panel
  Copyright (c) 2015 Dane Everitt

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

namespace PufferPanel\Core\Router\Node\Ajax;

use \Unirest,
    \Tracy\Debugger,
    PufferPanel\Core\OAuthService as OAuth2Service;

class Files extends \PufferPanel\Core\Files {

    use \PufferPanel\Core\Components\Error_Handler;

    /**
     * @param object
     */
    protected $server;

    /**
     * @param object
     */
    protected $params;

    /**
     * @param array
     */
    protected $display_folders = array();

    /**
     * @param array
     */
    protected $display_files = array();

    /**
     * Constructor class.
     */
    public function __construct(\PufferPanel\Core\Server $server) {

        $this->server = $server;
    }

    /**
     * Builds the contents of a given directory and makes them accessable through getFiles() and getFolders() functions.
     *
     * @param object $params The array of $klein->params()
     * @return bool
     */
    public function buildContents(array $params) {

        $this->params = $params;
        if (isset($this->params['dir']) && !empty($this->params['dir'])) {
            $this->params['dir'] = rtrim($this->params['dir'], '/');
        } else {
            $this->params['dir'] = null;
        }

        $contents = self::_retrieveFolderListing();

        if (!$contents) {

            self::_setError("Unable to connect to daemon to process this request.");
            return false;
        }

        if ($contents->code !== 200 || isset($contents->body->message)) {

            self::_setError("The daemon returned an error. (" . $contents->body->message . ")");
            return false;
        }

        foreach (json_decode($contents->raw_body, true) as $value) {

            if ($value['file'] !== true) {

                // @TODO handle symlinks

                $this->display_folders = array_merge($this->display_folders, array(array(
                        "entry" => $value['name'],
                        "directory" => trim($this->params['dir'], "/"),
                        "size" => null,
                        "date" => strtotime($value['modified'])
                )));
            } else {

                $this->display_files = array_merge($this->display_files, array(array(
                        "entry" => $value['name'],
                        "directory" => trim($this->params['dir'], "/"),
                        "extension" => pathinfo($value['name'], PATHINFO_EXTENSION),
                        "size" => $this->formatSize($value['size']),
                        "date" => strtotime($value['modified'])
                )));
            }
        }

        return true;
    }

    /**
     * Returns a listing of all files in a given directory after calling \PufferPanel\Core\Router\Node\Ajax\Files\buildContents();
     *
     * @return array
     */
    public final function getFiles() {

        return $this->display_files;
    }

    /**
     * Returns a listing of all folders in a given directory after calling \PufferPanel\Core\Router\Node\Ajax\Files\buildContents();
     *
     * @return array
     */
    public final function getFolders() {

        return $this->display_folders;
    }

    /**
     * Returns the contents of a given folder in GSD.
     *
     * @return bool|array
     */
    protected final function _retrieveFolderListing() {

        try {

            $attached_folder = (!is_null($this->params['dir'])) ? $this->params['dir'] : "/";

            $request = Unirest\Request::get(sprintf("http://%s:%s/server/%s/file/%s", $this->server->nodeData('fqdn'), $this->server->nodeData('daemon_listen'), $this->server->getData('hash'), $attached_folder), array(
                        'Authorization' => 'Bearer ' + OAuth2Service::get()->getPanelAccessToken()));

            return $request;
        } catch (\Exception $e) {

            Debugger::log($e);
            return false;
        }
    }

}

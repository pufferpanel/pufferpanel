<?php

namespace PufferPanel\Http\Controllers\Server;

use Log;
use Debugbar;
use PufferPanel\Models\Server;
use PufferPanel\Models\Node;
use PufferPanel\Http\Helpers;

use PufferPanel\Http\Controllers\Controller;
use Illuminate\Http\Request;

use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;

class AjaxController extends Controller
{

    /**
     * @var array
     */
    protected $files = [];

    /**
     * @var array
     */
    protected $folders = [];

    /**
     * @var string
     */
    protected $directory;

    /**
     * Listing of editable files in the control panel.
     * @var array
     */
    protected $editable = [
        'txt',
        'yml',
        'yaml',
        'log',
        'conf',
        'config',
        'html',
        'json',
        'properties',
        'props',
        'cfg',
        'lang',
        'ini',
        'cmd',
        'sh',
        'lua',
        '0' // Supports BungeeCord Files
    ];

    /**
     * Controller Constructor
     */
    public function __construct()
    {

        // All routes in this controller are protected by the authentication middleware.
        $this->middleware('auth');

        // Routes in this file are also checked aganist the server middleware. If the user
        // does not have permission to view the server it will not load.
        $this->middleware('server');

    }

    /**
     * Returns true or false depending on the power status of the requested server.
     *
     * @param  \Illuminate\Http\Request $request
     * @param  string $uuid
     * @return \Illuminate\Contracts\View\View
     */
    public function getStatus(Request $request, $uuid)
    {

        $server = Server::getByUUID($uuid);
        $client = Node::guzzleRequest($server->node);

        try {

            $res = $client->request('GET', '/server', [
                'headers' => Server::getGuzzleHeaders($uuid)
            ]);

            if($res->getStatusCode() === 200) {

                $json = json_decode($res->getBody());

                if (isset($json->status) && $json->status === 1) {
                    return 'true';
                }

            }

        } catch (RequestException $e) {
            Debugbar::error($e->getMessage());
            Log::notice('An exception was raised while attempting to contact a Scales instance to get server status information.', [
                'exception' => $e->getMessage(),
                'path' => $request->path()
            ]);
        }

        return 'false';
    }

    /**
     * Returns a listing of files in a given directory for a server.
     *
     * @param  \Illuminate\Http\Request $request
     * @param  string $uuid`
     * @return \Illuminate\Contracts\View\View
     */
    public function postDirectoryList(Request $request, $uuid)
    {

        $server = Server::getByUUID($uuid);
        $this->directory = '/' . trim($request->input('directory', '/'), '/');
        $this->authorize('list-files', $server);

        $prevDir = [
            'header' => ($this->directory !== '/') ? $this->directory : ''
        ];
        if ($this->directory !== '/') {
            $prevDir['first'] = true;
        }

        // Determine if we should show back links in the file browser.
        // This code is strange, and could probably be rewritten much better.
        $goBack = explode('/', rtrim($this->directory, '/'));
        if (isset($goBack[2]) && !empty($goBack[2])) {
            $prevDir['show'] = true;
            $prevDir['link'] = '/' . trim(str_replace(end($goBack), '', $this->directory), '/');
            $prevDir['link_show'] = trim($prevDir['link'], '/');
        }

        try {

            $client = Node::guzzleRequest($server->node);
            $directory = $request->input('directory', '/');

            $res = $client->request('GET', '/server/directory/' . $this->directory, [
                'headers' => Server::getGuzzleHeaders($uuid)
            ]);

            $json = json_decode($res->getBody());
            if($res->getStatusCode() !== 200 || isset($json->error)) {
                throw new \Exception('The response code from Scales was invalid: HTTP\\' . $res->getStatusCode());
            }

            $this->buildListing($json);

            return view('server.files.list', [
                'server' => $server,
                'files' => $this->files,
                'folders' => $this->folders,
                'extensions' => $this->editable,
                'directory' => $prevDir
            ]);

        } catch (\Exception $e) {
            Debugbar::addException($e);
            Log::notice('An exception was raised while attempting to contact a Scales instance to gather a directory listing.', [
                'exception' => $e,
                'path' => $request->path()
            ]);
        }

    }

    protected function buildListing($json)
    {

        foreach($json as &$value) {

            if ($value->file !== true) {

                // @TODO Handle Symlinks
                $this->folders = array_merge($this->folders, [[
                    'entry' => $value->name,
                    'directory' => trim($this->directory, '/'),
                    'size' => null,
                    'date' => strtotime($value->modified)
                ]]);

            } else {

                $this->files = array_merge($this->files, [[
                    'entry' => $value->name,
                    'directory' => trim($this->directory, '/'),
                    'extension' => pathinfo($value->name, PATHINFO_EXTENSION),
                    'size' => Helpers::bytesToHuman($value->size),
                    'date' => strtotime($value->modified)
                ]]);

            }

        }

    }

}

<?php

namespace PufferPanel\Http\Controllers\Server;

use PufferPanel\Models\Server;
use PufferPanel\Models\Node;
use Debugbar;

use PufferPanel\Http\Controllers\Controller;
use Illuminate\Http\Request;

use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;

class AjaxController extends Controller
{

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
     * @return \Illuminate\Contracts\View\View
     */
    public function getStatus(Request $request)
    {

        $server = Server::getByUUID($request->route()->server);
        $client = Node::guzzleRequest($server->node);

        try {

            $res = $client->request('GET', '/server', [
                'headers' => Server::getGuzzleHeaders($server->uuidShort)
            ]);

            if($res->getStatusCode() === 200) {

                $json = json_decode($res->getBody());

                if (isset($json->status) && $json->status === 1) {
                    return 'true';
                }

            }

        } catch (RequestException $e) {
            Debugbar::error($e->getMessage());
        }

        return 'false';
    }

}

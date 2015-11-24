<?php

namespace PufferPanel\Http\Controllers\Server;

use Gate;
use Auth;
use PufferPanel\Models\Server;
use PufferPanel\Models\Node;
use Debugbar;

use PufferPanel\Exceptions\DisplayException;
use PufferPanel\Http\Controllers\Scales\FileController;
use PufferPanel\Http\Controllers\Controller;
use Illuminate\Http\Request;

class ServerController extends Controller
{

    /**
     * Controller Constructor
     *
     * @return void
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
     * Renders server index page for specified server.
     *
     * @param  \Illuminate\Http\Request $request
     * @return \Illuminate\Contracts\View\View
     */
    public function getIndex(Request $request)
    {
        $server = Server::getByUUID($request->route()->server);
        return view('server.index', [
            'server' => $server,
            'node' => Node::find($server->node)
        ]);
    }

    /**
     * Renders file overview page.
     *
     * @param  Request $request
     * @return \Illuminate\Contracts\View\View
     */
    public function getFiles(Request $request)
    {

        $server = Server::getByUUID($request->route()->server);
        $this->authorize('list-files', $server);

        return view('server.files.index', [
            'server' => $server,
            'node' => Node::find($server->node)
        ]);
    }

    /**
     * Renders edit file page for a given file.
     *
     * @param  Request $request
     * @param  string  $uuid
     * @param  string  $file
     * @return \Illuminate\Contracts\View\View
     */
    public function getEditFile(Request $request, $uuid, $file)
    {

        $server = Server::getByUUID($uuid);
        $this->authorize('edit-files', $server);

        $fileInfo = (object) pathinfo($file);
        $controller = new FileController($uuid);

        try {
            $fileContent = $controller->returnFileContents($file);
        } catch (\Exception $e) {

            Debugbar::addException($e);
            $exception = 'An error occured while attempting to load the requested file for editing, please try again.';

            if ($e instanceof DisplayException) {
                $exception = $e->getMessage();
            }

            return redirect()->route('files.index', $uuid)->with('flash-error', $exception);

        }

        return view('server.files.edit', [
            'server' => $server,
            'node' => Node::find($server->node),
            'file' => $file,
            'contents' => $fileContent->contents,
            'directory' => (in_array($fileInfo->dirname, ['.', './', '/'])) ? '/' : trim($fileInfo->dirname, '/') . '/',
            'extension' => $fileInfo->extension
        ]);

    }

}

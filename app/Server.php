<?php

namespace PufferPanel;

use Auth;
use Illuminate\Database\Eloquent\Model;

class Server extends Model
{

    /**
     * The table associated with the model.
     *
     * @var string
     */
    protected $table = 'servers';

    /**
     * The attributes excluded from the model's JSON form.
     *
     * @var array
     */
    protected $hidden = ['daemonSecret'];

    /**
     * @var array
     */
    protected static $serverUUIDInstance = [];

    /**
     * Returns array of all servers owned by the logged in user.
     * Returns all active servers if user is a root admin.
     *
     * @return \Illuminate\Database\Eloquent\Collection
     */
    public static function getUserServers()
    {

        if (!Auth::user()) {
            return false;
        }

        $query = self::select('servers.*', 'nodes.name as nodeName', 'locations.long as location')
                    ->join('nodes', 'servers.node', '=', 'nodes.id')
                    ->join('locations', 'nodes.location', '=', 'locations.id')
                    ->where('active', 1);

        if (Auth::user()->root_admin !== 1) {
            // ->whereIn('servers.id', Permissions::serversAsArray());
            $query->where('owner', Auth::user()->id);
        }

        return $query->get();

    }

    /**
     * Returns a single server specified by UUID
     *
     * @param  string $uuid The Short-UUID of the server to return an object about.
     * @return \Illuminate\Database\Eloquent\Collection
     */
    public static function getByUUID($uuid)
    {

        if (!Auth::user()) {
            return false;
        }

        if (array_key_exists($uuid, self::$serverUUIDInstance)) {
            return self::$serverUUIDInstance[$uuid];
        }

        $query = self::where('uuidShort', $uuid)->where('active', 1);

        if (Auth::user()->root_admin !== 1) {
            // ->whereIn('id', Permissions::serversAsArray());
            $query->where('owner', Auth::user()->id);
        }

        self::$serverUUIDInstance[$uuid] = $query->first();
        return self::$serverUUIDInstance[$uuid];

    }

    /**
     * Returns non-administrative headers for accessing a server on Scales
     *
     * @param  string $uuid
     * @return array
     */
    public static function getGuzzleHeaders($uuid)
    {

        if (array_key_exists($uuid, self::$serverUUIDInstance)) {
            return [
                'X-Access-Server' => self::$serverUUIDInstance[$uuid]->uuid,
                'X-Access-Token' => self::$serverUUIDInstance[$uuid]->daemonSecret
            ];
        }

        return [];

    }

}

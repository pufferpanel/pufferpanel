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

        $query = self::where('active', 1);

        if (Auth::user()->root_admin !== 1) {
            $query->where('owner', Auth::user()->id);
        }

        return $query->get();

    }

    /**
     * Returns a single server specified by UUID
     *
     * @param  string $uuid The UUID of the server to return an object about.
     * @return \Illuminate\Database\Eloquent\Collection
     */
    public static function getByUUID($uuid)
    {

        if (!Auth::user()) {
            return false;
        }

        $query = self::where('uuid', $uuid)->where('active', 1);

        if (Auth::user()->root_admin !== 1) {
            $query->where('owner', Auth::user()->id);
        }

        return $query->take(1)->get();

    }

}

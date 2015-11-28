<?php

namespace PufferPanel\Models;

use Log;
use PufferPanel\Exceptions\DisplayException;
use PufferPanel\Models\APIPermission;
use Illuminate\Database\Eloquent\Model;

class API extends Model
{

    /**
     * The table associated with the model.
     *
     * @var string
     */
    protected $table = 'api';

    /**
     * The attributes excluded from the model's JSON form.
     *
     * @var array
     */
    protected $hidden = ['daemonSecret'];

    public function permissions()
    {
        return $this->hasMany(APIPermission::class);
    }

    public static function findKey($key)
    {
        return self::where('key', $key)->first();
    }

    /**
     * Determine if an API key has permission to perform an action.
     *
     * @param  string $key
     * @param  string $permission
     * @return boolean
     */
    public static function checkPermission($key, $permission)
    {
        $api = self::findKey($key);

        if (!$api) {
            throw new DisplayException('The requested API key (' . $key . ') was not found in the system.');
        }

        return APIPermission::check($api->id, $permission);

    }

    public static function noPermissionError($error = 'You do not have permission to perform this action with this API key.')
    {
        return response()->json([
            'error' => 'You do not have permission to perform this action with this API key.'
        ], 403);
    }

}

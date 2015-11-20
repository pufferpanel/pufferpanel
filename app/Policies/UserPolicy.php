<?php

namespace PufferPanel\Policies;

use Debugbar;
use PufferPanel\Models\User;
use PufferPanel\Models\Server;

class UserPolicy
{

    /**
     * Create a new policy instance.
     *
     * @return void
     */
    public function __construct()
    {
        //
    }

    // public function before($user, $ability)
    // {
    //     if ($user->root_admin === 1) {
    //         return true;
    //     }
    // }

    public function power(User $user, Server $server)
    {
        return true;
    }

}

<?php

namespace PufferPanel\Policies;

use Debugbar;
use PufferPanel\Models\User;
use PufferPanel\Models\Server;

class ServerPolicy
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

    protected function isOwner(User $user, Server $server)
    {
        return $server->owner === $user->id;
    }

    public function before($user, $ability)
    {
        if ($user->root_admin === 1) {
            return true;
        }
    }

    public function power(User $user, Server $server)
    {
        if ($this->isOwner($user, $server)) {
            return true;
        }

        return $user->permissions()->server($server)->permission('power')->exists();
    }

    public function command(User $user, Server $server)
    {
        if ($this->isOwner($user, $server)) {
            return true;
        }

        return $user->permissions()->server($server)->permission('command')->exists();
    }

    public function listFiles(User $user, Server $server)
    {
        if ($this->isOwner($user, $server)) {
            return true;
        }

        return $user->permissions()->server($server)->permission('list-files')->exists();
    }

}

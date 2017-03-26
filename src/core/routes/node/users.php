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
namespace PufferPanel\Core\Router\Node;

use \ORM, \PufferPanel\Core\Settings;
use PufferPanel\Core\OAuthService;

class Users extends \PufferPanel\Core\Email
{

    use \PufferPanel\Core\Components\Authentication, \PufferPanel\Core\Components\Daemon, \PufferPanel\Core\Components\Error_Handler;

    /**
     * @param object
     */
    protected $server;

    /**
     * Constructor class for \PufferPanel\Core\Router\Node\Users
     */
    public function __construct(\PufferPanel\Core\Server $server)
    {
        $this->server = $server;
    }

    /**
     * Handles adding a new subuser to a server.
     *
     * @return bool
     */
    public function addSubuser(\Klein\DataCollection\DataCollection $data)
    {

        if (is_null($data->permissions) || !$data->permissions) {
            $data->permissions = array('console.view');
        }

        $registerToken = $this->keygen(32);
        $data->permissions = self::_rebuildUserPermissions($data->permissions);

        ORM::get_db()->beginTransaction();

        try {

            $user = ORM::forTable('users')->where('email', $data->email)->findOne();

            $checkUserExists = $user ? true : false;
            if (!$checkUserExists) {

                $user = ORM::forTable('users')->create()->set(array(
                    'uuid' => $this->generateUniqueUUID('users', 'uuid'),
                    'username' => $data->email,
                    'email' => $data->email,
                    'password' => password_hash(OAuthService::generateSecret(), PASSWORD_BCRYPT),
                    'language' => Settings::config('default_language'),
                    'register_time' => time()
                ));
                $user->save();

                ORM::forTable('account_change')->create()->set(array(
                    'user_id' => $user->id,
                    'type' => 'user_register',
                    'content' => $data->email,
                    'key' => $registerToken,
                    'time' => time()
                ))->save();

            }

            ORM::forTable('subusers')->create()->set(array(
                'user' => $user->id,
                'server' => $this->server->getData('id')
            ))->save();

            $email = null;

            foreach ($data->permissions as $id => $permission) {

                ORM::forTable('permissions')->create()->set(array(
                    'user' => $user->id,
                    'server' => $this->server->getData('id'),
                    'permission' => $permission
                ))->save();

            }
            $daemonPerms = self::_getDaemonPermissions($data->permissions);

            OAuthService::Get()->create(ORM::get_db(),
                $user->id(),
                $this->server->getData('id'),
                '.internal_' . $user->id . '_' . $this->server->getData('id'),
                implode(' ', $daemonPerms),
                'internal_use',
                'internal_use'
            );

            if ($checkUserExists) {

                $email = $this->buildEmail('new_subuser', array(
                    'SERVER' => $this->server->getData('name'),
                    'EMAIL' => $data->email
                ));

            } else {

                $email = $this->buildEmail('new_subuser_createaccount', array(
                    'REGISTER_TOKEN' => $registerToken,
                    'URLENCODE_TOKEN' => urlencode($registerToken),
                    'SERVER' => $this->server->getData('name'),
                    'EMAIL' => $data->email
                ));

            }

            $email->dispatch($data->email, Settings::config()->company_name . ' - You\'ve Been Invited to Manage a Server');

            ORM::get_db()->commit();
            return true;

        } catch (\Exception $e) {

            \Tracy\Debugger::log($e);
            self::_setError("An error occurred when trying to update the server information.");
            ORM::get_db()->rollBack();
            return false;

        }

    }

    /**
     * Updates information for a subuser.
     *
     * @param array $data
     * @return bool
     */
    public function modifySubuser(\Klein\DataCollection\DataCollection $data)
    {

        $select = ORM::forTable('subusers')->where(array(
            'user' => $data->user_id,
            'server' => $this->server->getData('id')
        ))->findOne();

        if (!$select) {
            self::_setError("Invalid user was provided.");
            return false;
        }

        ORM::forTable('permissions')->where(array(
            'user' => $data->user_id,
            'server' => $this->server->getData('id')
        ))->deleteMany();

        if (is_null($data->permissions) || !$data->permissions) {
            $data->permissions = array('console.view');
        }

        $permissions = self::_rebuildUserPermissions($data->permissions);

        foreach ($permissions as $id => $permission) {

            ORM::forTable('permissions')->create()->set(array(
                'user' => $select->user,
                'server' => $this->server->getData('id'),
                'permission' => $permission
            ))->save();

        }

        $clientIds = ORM::forTable('oauth_clients')
            ->where(array(
                'user_id' => $select->user,
                'server_id' => $this->server->getData('id')
            ))->find_many();

        $daemonPerms = self::_getDaemonPermissions($data->permissions);
        $daemonPermsString = implode(' ', $daemonPerms);

        foreach($clientIds as $client) {
            ORM::forTable('oauth_access_tokens')
                ->where('oauthClientId', $client->id)
                ->deleteMany();

            $client->scopes = $daemonPermsString;
            $client->save();
        }

        return true;

    }

    /**
     * Rebuilds user permissions to include the parent .view permission.
     *
     * @param array $data
     * @return array
     * @static
     */
    protected final static function _rebuildUserPermissions(array $data)
    {

        foreach ($data as $permission) {

            if (in_array($permission, array('files.edit', 'files.save', 'files.download', 'files.delete', 'files.create', 'files.upload', 'files.zip')) && !in_array('files.view', $data)) {
                $data = array_merge($data, array("files.view"));
            }

            if (in_array($permission, array('manage.rename.jar')) && !in_array('manage.rename.view', $data)) {
                $data = array_merge($data, array("manage.rename.view"));
            }

            if (in_array($permission, array('manage.ftp.details', 'manage.ftp.password')) && !in_array('manage.ftp.view', $data)) {
                $data = array_merge($data, array("manage.ftp.view"));
            }

            if (in_array($permission, array('manage.ftp.details', 'manage.ftp.password', 'manage.rename.jar', 'manage.ftp.view', 'manage.rename.view')) && !in_array('manage.view', $data)) {
                $data = array_merge($data, array("manage.view"));
            }

        }

        return $data;

    }

    /**
     * Revokes subuser permissions for a given user that has an active account on the panel.
     *
     * @param object $orm Database query object.
     * @return bool
     */
    public function revokeActiveUserPermissions(ORM $orm)
    {

        if ($orm->pending == 0) {

            try {
                ORM::forTable('permissions')
                    ->where(array(
                        'user' => $orm->user,
                        'server' => $this->server->getData('id')
                    ))->delete_many();

                $orm->delete();

                return true;

            } catch (\Exception $e) {

                \Tracy\Debugger::log($e, Debugger::Error);
                self::_setError("An error occurred when trying to update the server.");
                return false;

            }

        } else {

            ORM::forTable('account_change')->where(array(
                'content' => $orm->email,
                'type' => 'user_register'
            ))->findOne()->delete();
            ORM::forTable('subusers')->where(array(
                'id' => $orm->user,
                'server' => $this->server->getData('id')
            ))->findOne()->delete();
            $orm->delete();
            return true;

        }

    }

    protected static function _getDaemonPermissions(array $data)
    {
        if(!in_array('console.view', $data)) {
            $data[] = 'console.view';
        }

        $daemonPerms = array();
        foreach($data as $key => $value) {
            switch($value) {
                case "console.view":
                    $daemonPerms[] = "server.console";
                    $daemonPerms[] = "server.stats";
                    $daemonPerms[] = "server.network";
                    break;
                case "console.commands":
                    $daemonPerms[] = "server.console.send";
                    break;
                case "console.power":
                    $daemonPerms[] = "server.install";
                    $daemonPerms[] = "server.start";
                    $daemonPerms[] = "server.stop";
                    break;
                case "files.view":
                    $daemonPerms[] = "server.file.get";
                    break;
                case "files.edit":
                    $daemonPerms[] = "server.file.get";
                    break;
                case "files.save":
                    $daemonPerms[] = "server.file.put";
                    break;
                case "files.download":
                    $daemonPerms[] = "server.file.get";
                    break;
                case "files.zip":
                    break;
                case "manage.oauth2":
                    break;
                case "sftp":
                    $daemonPerms[] = "sftp";
                    break;
            }
        }
        return $daemonPerms;
    }

}

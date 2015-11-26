<?php

namespace PufferPanel\Models;

use Google2FA;
use PufferPanel\Exceptions\AccountNotFound;
use PufferPanel\Models\Permission;

use Illuminate\Auth\Authenticatable;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Auth\Passwords\CanResetPassword;
use Illuminate\Foundation\Auth\Access\Authorizable;
use Illuminate\Contracts\Auth\Authenticatable as AuthenticatableContract;
use Illuminate\Contracts\Auth\Access\Authorizable as AuthorizableContract;
use Illuminate\Contracts\Auth\CanResetPassword as CanResetPasswordContract;

class User extends Model implements AuthenticatableContract,
                                    AuthorizableContract,
                                    CanResetPasswordContract
{
    use Authenticatable, Authorizable, CanResetPassword;

    /**
     * The table associated with the model.
     *
     * @var string
     */
    protected $table = 'users';

    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = ['name', 'email', 'password'];

    /**
     * The attributes excluded from the model's JSON form.
     *
     * @var array
     */
    protected $hidden = ['password', 'remember_token', 'totp_secret'];

    public function permissions()
    {
        return $this->hasMany(Permission::class);
    }

    /**
     * Sets the TOTP secret for an account.
     *
     * @param int $id Account ID for which we want to generate a TOTP secret
     * @return string
     */
    public function setTotpSecret($id)
    {

        $totpSecretKey = Google2FA::generateSecretKey();

        $user = User::find($id);

        if (!$user) {
            throw new AccountNotFound('An account with that ID (' . $id . ') does not exist in the system.');
        }

        $user->totp_secret = $totpSecretKey;
        $user->save();

        return $totpSecretKey;

    }

    /**
     * Enables or disables TOTP on an account if the token is valid.
     *
     * @param int $id Account ID for which we want to generate a TOTP secret
     * @return boolean
     */
    public function toggleTotp($id, $token)
    {

        $user = User::find($id);

        if (!$user) {
            throw new AccountNotFound('An account with that ID (' . $id . ') does not exist in the system.');
        }

        if (!Google2FA::verifyKey($user->totp_secret, $token)) {
            return false;
        }

        $user->use_totp = ($user->use_totp === 1) ? 0 : 1;
        $user->save();

        return true;

    }

}

<?php

namespace PufferPanel\Console\Commands;

use Illuminate\Console\Command;
use PufferPanel\Models\User;

class InstallPanel extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'panel:install';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Runs the initial panel setup: creates admin user';

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        $uuid = sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x', mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0x0fff) | 0x4000, mt_rand(0, 0x3fff) | 0x8000, mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff));
        $username = $this->ask('What is your username?');
        $email = $this->ask('What is your email address?');
        $password = $this->secret('What is your password?');

        $user = User::create([
            'uuid' => $uuid,
            'username' => $username,
            'email' => $email,
            'password' => bcrypt($password),
            'language' => 'en',
            'time' => ''
        ]);

        $user->root_admin = 1;
        $user->save();

        $this->line('Your admin account has been created.');
    }
}
<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class UpdateForeign extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        // Schema::table('nodes', function (Blueprint $table) {
        //     $table->foreign('location')->references('id')->on('locations');
        // });
        // 
        // Schema::table('servers', function (Blueprint $table) {
        //     $table->foreign('owner')->references('id')->on('users');
        //     $table->foreign('node')->references('id')->on('nodes');
        // });
        //
        // Schema::table('daemon', function (Blueprint $table) {
        //     $table->foreign('server')->references('id')->on('servers');
        // });
        //
        // Schema::table('allocations', function (Blueprint $table) {
        //     $table->foreign('node')->references('id')->on('nodes');
        // });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
    }
}

<?php

/*
  PufferPanel - A Minecraft Server Management Panel
  Copyright (c) 2013 Dane Everitt

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

namespace PufferPanel\Core;

class DatabaseManager {

    private static $queries = [];
    private static $queryCount = 0;

    public final static function logQuery($query, $time) {

        self::$queries[] = array(
            'query' => $query,
            'execution_time' => $time,
            'time' => time()
        );
        self::$queryCount++;
    }

    public final static function returnQueries() {

        return self::$queries;
    }

    public final static function returnQueryCount() {

        return self::$queryCount;
    }

    public final static function logQueriesToDebugger() {

        \Tracy\Debugger::log(self::$queries);
    }

}

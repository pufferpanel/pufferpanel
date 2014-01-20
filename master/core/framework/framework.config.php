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
 
/*
 * Config Framework Core
 */

class Config
{
    /**
     * Array of cached configuration settings
     * @var array
     */
    protected $config = array();

    /**
     * Gets a config setting by its path, in dot format. For example, db.user
     * gets "username" from a file config/db.php that looks like:
     *     <?php return array('user' => 'foo') ?>
     * @param string $path
     * @return mixed
     */
    public function get($path)
    {
        if (!$this->existsInCache($path)) {
            list($file, $reference) = explode('.', $path, 1);
            $this->config[$file] = (require 'config/' . $file);
        }

        return $this->getFromCache($path);
    }

    /**
     * Adds the given config option to the cache. Does NOT save to file.
     * @param string $path
     * @param mixed $value
     * @return void
     */
    public function set($path, $value) {

        $pointer = $value;
        foreach (array_reverse(explode('.', $path)) as $part) {
            $pointer = array($part => $pointer);
        }

        $this->config = array_merge_recursive($this->config, $pointer);
    }

    /**
     * Checks to see if a given item exists in the cache.
     * @param string $path
     * @return bool
     */
    protected function existsInCache($path) {
        $pointer = $this->config;

        foreach (explode('.', $path) as $part) {
            if (!array_key_exists($part, $pointer)) {
                return false;
            }
        }
        return true;
    }

    /**
     * Gets a given item by its path from the cache.
     * @param string $path
     * @return mixed
     */
    protected function getFromCache($path) {
        $pointer = $this->config;

        foreach (explode('.', $path) as $part) {
            $pointer = $pointer[$part];
        }

        return $pointer;
    }
}
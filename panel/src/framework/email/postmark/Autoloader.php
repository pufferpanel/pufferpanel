<?php

namespace Postmark;

/**
 * Postmark PHP class
 *
 * Copyright 2011, Markus Hedlund, Mimmin AB, www.mimmin.com
 * Licensed under the MIT License.
 * Redistributions of files must retain the above copyright notice.
 *
 * @author Markus Hedlund (markus@mimmin.com) at mimmin (www.mimmin.com)
 * @copyright Copyright 2009 - 2011, Markus Hedlund, Mimmin AB, www.mimmin.com
 * @version 0.5
 * @license http://www.opensource.org/licenses/mit-license.php The MIT License
 */

class Autoloader {
	/**
	 * Register the autoloader
	 * 
	 * @return  void
	 */
	public static function register() {
		ini_set('unserialize_callback_func', 'spl_autoload_call');
		spl_autoload_register(array(new self, 'autoload'));
	}

	/**
	 * Autoloader
	 *
	 * @param   string
	 * @return  void
	 */
	public static function autoload( $class ) {
		if (0 !== strpos($class, 'Postmark\\')) {
			return;
		} else if (file_exists($file = dirname(__FILE__) . '/' . preg_replace('!^Postmark\\\!', '', $class) . '.php')) {
			require $file;
		}
    }
}
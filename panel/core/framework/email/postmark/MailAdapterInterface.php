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

interface MailAdapterInterface
{
	public static function getApiKey();
	public static function setupDefaults(Mail &$mail);
	public static function log($logData);
}
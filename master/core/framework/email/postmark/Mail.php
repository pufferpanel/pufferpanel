<?php

namespace Postmark;

use ReflectionClass, Mail_Postmark_Adapter, InvalidArgumentException, OverflowException, Exception, BadMethodCallException;

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

class Mail
{
	const DEBUG_OFF = 0;
	const DEBUG_VERBOSE = 1;
	const DEBUG_RETURN = 2;
	const TESTING_API_KEY = 'POSTMARK_API_TEST';
	const MAX_ATTACHMENT_SIZE = 10485760; // 10 MB
	const RECIPIENT_TYPE_TO = 'to';
	const RECIPIENT_TYPE_CC = 'cc';
	const RECIPIENT_TYPE_BCC = 'bcc';

	static $_mimeTypes = array('ai' => 'application/postscript', 'avi' => 'video/x-msvideo', 'doc' => 'application/msword', 'eps' => 'application/postscript', 'gif' => 'image/gif', 'htm' => 'text/html', 'html' => 'text/html', 'jpeg' => 'image/jpeg', 'jpg' => 'image/jpeg', 'mov' => 'video/quicktime', 'mp3' => 'audio/mpeg', 'mpg' => 'video/mpeg', 'pdf' => 'application/pdf', 'ppt' => 'application/vnd.ms-powerpoint', 'ps' => 'application/postscript', 'rtf' => 'application/rtf', 'tif' => 'image/tiff', 'tiff' => 'image/tiff', 'txt' => 'text/plain', 'xls' => 'application/vnd.ms-excel', 'csv' => 'text/comma-separated-values', 'docx' => 'application/vnd.openxmlformats-officedocument.wordprocessingml.document', 'flv' => 'video/x-flv', 'ics' => 'text/calendar', 'log' => 'text/plain', 'png' => 'image/png', 'pptx' => 'application/vnd.openxmlformats-officedocument.presentationml.presentation', 'psd' => 'image/photoshop', 'rm' => 'application/vnd.rn-realmedia', 'swf' => 'application/x-shockwave-flash', 'xlsx' => 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 'xml' => 'text/xml');

	private $_apiKey;
	private $_from;
	private $_to = array();
	private $_cc = array();
	private $_bcc = array();
	private $_replyTo;
	private $_subject;
	private $_tag;
	private $_messagePlain;
	private $_messageHtml;
	private $_headers = array();
	private $_attachments = array();
	private $_debugMode = self::DEBUG_OFF;

	/**
	 * Initialize
	 *
	 * @param   string  $apiKey Postmark server API key
	 * @return  void
	 */
	public function __construct($apiKey)
	{
		if (class_exists('Mail_Postmark_Adapter', false)) {
			$reflection = new ReflectionClass('Mail_Postmark_Adapter');

			if (!$reflection->implementsInterface('MailAdapterInterface')) {
				trigger_error('Mail_Postmark_Adapter must implement interface MailAdapterInterface', E_USER_ERROR);
			}

			$this->_apiKey = Mail_Postmark_Adapter::getApiKey();

			Mail_Postmark_Adapter::setupDefaults($this);

		} else {
			$this->_apiKey = $apiKey;
		}

		$this->messageHtml(null)->messagePlain(null);
	}

	/**
	 * Add a physical file as an attachment
	 * Options:
	 * - filenameAlias, use a different filename for the attachment
	 *
	 * @param string $filename Location of the file
	 * @param array $options An optional array with options
	 * @throws InvalidArgumentException If file doesn't exist
	 * @throws OverflowException If maximum attachment size has been reached
	 * @return Mail
	 */
	public function &addAttachment($filename, $options = array())
	{
		if (!is_file($filename)) {
			throw new InvalidArgumentException("File \"{$filename}\" does not exist");
		}

		$this->addCustomAttachment(
			isset($options['filenameAlias']) ? $options['filenameAlias'] : basename($filename),
			file_get_contents($filename),
			$this->_getMimeType($filename)
		);

		return $this;
	}

	/**
	 * Add a BCC address
	 *
	 * @param string $address E-mail address used in BCC
	 * @param string $name Optional. Name used in BCC
	 * @throws InvalidArgumentException On invalid address
	 * @throws OverflowException If there are too many email recipients
	 * @return Mail
	 */
	public function &addBcc($address, $name = null)
	{
		$this->_addRecipient(self::RECIPIENT_TYPE_BCC, $address, $name);
		return $this;
	}

	/**
	 * Add a CC address
	 *
	 * @param string $address E-mail address used in CC
	 * @param string $name Optional. Name used in CC
	 * @throws InvalidArgumentException On invalid address
	 * @throws OverflowException If there are too many email recipients
	 * @return Mail
	 */
	public function &addCc($address, $name = null)
	{
		$this->_addRecipient(self::RECIPIENT_TYPE_CC, $address, $name);
		return $this;
	}

	/**
	 * Add an attachment.
	 *
	 * @param string $filename What to call the file
	 * @param string $content Raw file data
	 * @param string $mimeType The mime type of the file
	 * @throws OverflowException If maximum attachment size has been reached
	 * @return Mail
	 */
	public function &addCustomAttachment($filename, $content, $mimeType)
	{
		$length = strlen($content);
		$lengthSum = 0;

		foreach ($this->_attachments as $file) {
			$lengthSum += $file['length'];
		}

		if ($lengthSum + $length > self::MAX_ATTACHMENT_SIZE) {
			throw new OverflowException("Maximum attachment size reached");
		}

		$this->_attachments[$filename] = array(
			'content' => base64_encode($content),
			'mimeType' => $mimeType,
			'length' => $length
		);

		return $this;
	}

	/**
	 * Add a custom header
	 *
	 * @param string $name Custom header name
	 * @param string $value Custom header value
	 * @return Mail
	 */
	public function &addHeader($name, $value)
	{
		$this->_headers[$name] = $value;
		return $this;
	}

	/**
	 * Add a receiver
	 *
	 * @param string $address E-mail address used in To
	 * @param string $name Optional. Name used in To
	 * @throws InvalidArgumentException On invalid address
	 * @throws OverflowException If there are too many email recipients
	 * @return Mail
	 */
	public function &addTo($address, $name = null)
	{
		$this->_addRecipient(self::RECIPIENT_TYPE_TO, $address, $name);
		return $this;
	}

	/**
	 * New e-mail
	 *
	 * @param   string  $apiKey Postmark server API key
	 * @return  Mail
	 */
	public static function compose($apiKey)
	{
		return new self($apiKey);
	}

	/**
	 * Turns debug output on
	 *
	 * @param int $mode One of the debug constants
	 * @return Mail
	 */
	public function &debug($mode = self::DEBUG_VERBOSE)
	{
		$this->_debugMode = $mode;
		return $this;
	}

	/**
	 * Specify sender. Overwrites default From. Note that the address
	 * must first be added in the Postmark app admin interface
	 *
	 * @param string $address E-mail address used in From
	 * @param string $name Optional. Name used in From
	 * @throws InvalidArgumentException On invalid address
	 * @return Mail
	 */
	public function &from($address, $name = null)
	{
		if (!$this->_validateAddress($address)) {
			throw new InvalidArgumentException("From address \"{$address}\" is invalid");
		}

		$this->_from = array('address' => $address, 'name' => $name);
		return $this;
	}

	/**
	 * Specify sender name. Overwrites default From name, but doesn't change address.
	 *
	 * @param string $name Name used in From
	 * @return Mail
	 */
	public function &fromName($name)
	{
		$this->_from['name'] = $name;
		return $this;
	}

	/**
	 * Add HTML message. Can be used in conjunction with messagePlain()
	 *
	 * @param string $message E-mail message
	 * @return Mail
	 */
	public function &messageHtml($message)
	{
		$this->_messageHtml = $message;
		return $this;
	}

	/**
	 * Add plaintext message. Can be used in conjunction with messageHtml()
	 * @param string $message E-mail message
	 * @return Mail
	 */
	public function &messagePlain($message)
	{
		$this->_messagePlain = $message;
		return $this;
	}

	/**
	 * Specify reply-to
	 *
	 * @param string $address E-mail address used in To
	 * @param string $name Optional. Name used in To
	 * @throws InvalidArgumentException On invalid address
	 * @return Mail
	 */
	public function &replyTo($address, $name = null)
	{
		if (!$this->_validateAddress($address)) {
			throw new InvalidArgumentException("Reply To address \"{$address}\" is invalid");
		}

		$this->_replyTo = array('address' => $address, 'name' => $name);
		return $this;
	}

	/**
	 * Sends the e-mail. Prints debug output if debug mode is turned on
	 *
	 * Options:
	 * 		returnMessageId
	 *
	 * @throws Exception If HTTP code 422, Exception with API error code and Postmark message, otherwise HTTP code.
	 * @throws BadMethodCallException If From address, To address or Subject is missing
	 * @return boolean - True if success and $returnID is false.
	 * @return string - if $returnID is true and one message is sent.
	 * @return array - if DEBUG_RETURN is enabled.
	 */
	public function send($options = array())
	{
		$this->_validateData();
		$data = $this->_prepareData();
		$headers = array(
			'Accept: application/json',
			'Content-Type: application/json',
			'X-Postmark-Server-Token: ' . $this->_apiKey
		);

		$ch = curl_init();
		curl_setopt($ch, CURLOPT_URL, 'https://api.postmarkapp.com/email');
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'POST');
		curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
		curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
		curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);
		curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, 2);
		curl_setopt($ch, CURLOPT_CAINFO, dirname(__FILE__) . '/Certificate/cacert.pem');

		$return = curl_exec($ch);
		$curlError = curl_error($ch);
		$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);

		$this->_log(array(
			'messageData' => $data,
			'return' => $return,
			'curlError' => $curlError,
			'httpCode' => $httpCode
		));

		if (($this->_debugMode & self::DEBUG_VERBOSE) === self::DEBUG_VERBOSE) {
			echo "JSON: " . json_encode($data)
				. "\nHeaders: \n\t" . implode("\n\t", $headers)
				. "\nReturn:\n{$return}"
				. "\nCurl error: {$curlError}"
				. "\nHTTP code: {$httpCode}";
		}

		if ($curlError !== '') {
			throw new Exception($curlError);
		}

		if (!$this->_isTwoHundred($httpCode)) {
			if ($httpCode == 422) {
				$return = json_decode($return);
				throw new Exception($return->Message, $return->ErrorCode);
			} else {
				throw new Exception("Error while mailing. Postmark returned HTTP code {$httpCode} with message \"{$return}\"", $httpCode);
			}
		}

		if (($this->_debugMode & self::DEBUG_RETURN) === self::DEBUG_RETURN) {
			return array(
				'json' => json_encode($data),
				'headers' => $headers,
				'return' => $return,
				'curlError' => $curlError,
				'httpCode' => $httpCode
			);
		}

		// Return the ID of the message sent if the option is set.
		if(!empty($options['returnMessageId'])) {
			$messageInformation = json_decode($return);
			return $messageInformation->MessageID;
		}

		return true;
	}

	/**
	 * Specify subject
	 *
	 * @param string $subject E-mail subject
	 * @return Mail
	 */
	public function &subject($subject)
	{
		$this->_subject = $subject;
		return $this;
	}

	/**
	 * You can categorize outgoing email using the optional Tag  property.
	 * If you use different tags for the different types of emails your
	 * application generates, you will be able to get detailed statistics
	 * for them through the Postmark user interface.
	 * Only 1 tag per email is supported.
	 *
	 * @param string $tag One tag
	 * @return Mail
	 */
	public function &tag($tag)
	{
		$this->_tag = $tag;
		return $this;
	}

	/**
	 * Specify receiver. Use addTo to add more.
	 *
	 * @deprecated Use addTo.
	 * @param string $address E-mail address used in To
	 * @param string $name Optional. Name used in To
	 * @return Mail
	 */
	public function &to($address, $name = null)
	{
		$this->_to = array();
		$this->addTo($address, $name);
		return $this;
	}

	/**
	 * @param string $type Either 'to', 'cc' or 'bcc'
	 * @param string $address
	 * @param string|null $name
	 * @throws InvalidArgumentException On invalid address
	 * @throws OverflowException If there are too many email recipients
	 */
	public function _addRecipient($type, $address, $name = null)
	{
		$address = trim($address);

		if (!$this->_validateAddress($address)) {
			throw new InvalidArgumentException("Address \"{$address}\" is invalid");
		}

		if (count($this->_to) + count($this->_cc) + count($this->_bcc) === 20) {
			throw new OverflowException('Too many email recipients');
		}

		$data = array('address' => $address, 'name' => $name);

		switch ($type) {
			case self::RECIPIENT_TYPE_TO:
				$this->_to[] = $data;
				break;

			case self::RECIPIENT_TYPE_CC:
				$this->_cc[] = $data;
				break;

			case self::RECIPIENT_TYPE_BCC:
				$this->_bcc[] = $data;
				break;
		}
	}

	private function _createAddress($address, $name = null)
	{
		if (isset($name)) {
			return '"' . str_replace('"', '', $name) . '" <' . $address . '>';
		} else {
			return $address;
		}
	}

	/**
	 * Try to detect the mime type
	 * @param $filename
	 * @return string
	 */
	private function _getMimeType($filename)
	{
		$extension = pathinfo($filename, PATHINFO_EXTENSION);

		if (isset(self::$_mimeTypes[$extension])) {
			return self::$_mimeTypes[$extension];

		} else if (function_exists('mime_content_type')) {
			return mime_content_type($filename);

		} else if (function_exists('finfo_file')) {
			$fh = finfo_open(FILEINFO_MIME);
			$mime = finfo_file($fh, $filename);
			finfo_close($fh);
			return $mime;

		} else if ($image = getimagesize($filename)) {
			return $image[2];

		} else {
			return 'application/octet-stream';
		}
	}

	/**
	 * If a number is 200-299
	 * @param $value
	 * @return bool
	 */
	private function _isTwoHundred($value)
	{
		return intval($value / 100) == 2;
	}

	/**
	 * Call the logger method, if one exists
	 *
	 * @param array $logData
	 */
	private function _log($logData)
	{
		if (class_exists('Mail_Postmark_Adapter', false)) {
			Mail_Postmark_Adapter::log($logData);
		}
	}

	/**
	 * Prepares the data array
	 * @return array
	 */
	private function _prepareData()
	{
		$data = array(
			'Subject' => $this->_subject
		);

		$data['From'] = $this->_createAddress($this->_from['address'], $this->_from['name']);
		$data['To'] = array();
		$data['Cc'] = array();
		$data['Bcc'] = array();

		foreach ($this->_to as $to) {
			$data['To'][] = $this->_createAddress($to['address'], $to['name']);
		}

		foreach ($this->_cc as $cc) {
			$data['Cc'][] = $this->_createAddress($cc['address'], $cc['name']);
		}

		foreach ($this->_bcc as $bcc) {
			$data['Bcc'][] = $this->_createAddress($bcc['address'], $bcc['name']);
		}

		$data['To'] = implode(', ', $data['To']);

		if (empty($data['Cc'])) {
			unset($data['Cc']);
		} else {
			$data['Cc'] = implode(', ', $data['Cc']);
		}

		if (empty($data['Bcc'])) {
			unset($data['Bcc']);
		} else {
			$data['Bcc'] = implode(', ', $data['Bcc']);
		}

		if ($this->_replyTo !== null) {
			$data['ReplyTo'] = $this->_createAddress($this->_replyTo['address'], $this->_replyTo['name']);
		}

		if ($this->_messageHtml !== null) {
			$data['HtmlBody'] = $this->_messageHtml;
		}

		if ($this->_messagePlain !== null) {
			$data['TextBody'] = $this->_messagePlain;
		}

		if ($this->_tag !== null) {
			$data['Tag'] = $this->_tag;
		}

		if (!empty($this->_headers)) {
			$data['Headers'] = array();

			foreach ($this->_headers as $name => $value) {
				$data['Headers'][] = array('Name' => $name, 'Value' => $value);
			}
		}

		if (!empty($this->_attachments)) {
			$data['Attachments'] = array();

			foreach ($this->_attachments as $filename => $file) {
				$data['Attachments'][] = array(
					'Name' => $filename,
					'Content' => $file['content'],
					'ContentType' => $file['mimeType']
				);
			}
		}

		return $data;
	}

	/**
	 * Validates an e-mail address
	 * @param $email
	 * @return bool
	 */
	private function _validateAddress($email)
	{
		// http://php.net/manual/en/function.filter-var.php
		// return filter_var($email, FILTER_VALIDATE_EMAIL) !== false;
		// filter_var proved to be unworthy (passed foo..bar@domain.com as valid),
		// and was therefore replace with
		$regex = "/^([\w\!\#$\%\&\'\*\+\-\/\=\?\^\`{\|\}\~]+\.)*[\w\!\#$\%\&\'\*\+\-\/\=\?\^\`{\|\}\~]+@((((([a-z0-9]{1}[a-z0-9\-]{0,62}[a-z0-9]{1})|[a-z])\.)+[a-z]{2,6})|(\d{1,3}\.){3}\d{1,3}(\:\d{1,5})?)$/i";
		// from http://fightingforalostcause.net/misc/2006/compare-email-regex.php
		return preg_match($regex, $email) === 1;
	}

	/**
	 * Validate that the email can be sent
	 *
	 * @throws BadMethodCallException If From address, To address or Subject is missing
	 */
	private function _validateData()
	{
		if ($this->_from['address'] === null) {
			throw new BadMethodCallException('From address is not set');
		}

		if (empty($this->_to)) {
			throw new BadMethodCallException('No To address is set');
		}

		if (!isset($this->_subject)) {
			throw new BadMethodCallException('Subject is not set');
		}
	}
}
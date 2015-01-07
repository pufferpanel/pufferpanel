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
use \ORM as ORM, \PufferPanel\Core\Config\DatabaseConfig;


/**
 * PufferPanel Core Email Sending Class
 */
class Email {

	use Components\Authentication;

	/**
	 * @param string $message
	 */
	protected $message;
	protected $masterurl;
	protected $settings;

	/**
	 * Constructor for email sending
	 *
	 * @return void
	 */
	public function __construct() {

		$this->settings = new DatabaseConfig('acp_settings', 'setting_ref', 'setting_val');
		$this->masterurl = ($this->settings->config('https') == 1) ? 'https:' . $this->settings->config('master_url') : 'http:' . $this->settings->config('master_url');

	}

	/**
	 * Sends an email that has been formatted.
	 *
	 * @param string $email The email address to send to.
	 * @param string $subject The subject of the email.
	 * @return void
	 */
	public function dispatch($email, $subject) {

		$this->getDispatchSystem = $this->getDispatchSystemFunct();
		if ($this->getDispatchSystem == 'php') {

			$headers = 'From: ' . $this->settings->config('sendmail_email') . "\r\n" .
					'Reply-To: ' . $this->settings->config('sendmail_email') . "\r\n" .
					'MIME-Version: 1.0' . "\r\n" .
					'Content-type: text/html; charset=iso-8859-1' . "\r\n" .
					'X-Mailer: PHP/' . phpversion();

			mail($email, $subject, $this->message, $headers);

		} else if ($this->getDispatchSystem == 'postmark') {

			\Postmark\Mail::compose($this->settings->config('postmark_api_key'))
					->from($this->settings->config('sendmail_email'), $this->settings->config('company_name'))
					->addTo($email, $email)
					->subject($subject)
					->messageHtml($this->message)
					->send();

		} else if ($this->getDispatchSystem == 'mandrill') {

			try {

				$mandrill = new \Mandrill($this->settings->config('mandrill_api_key'));
				$mandrillMessage = array(
					'html' => $this->message,
					'subject' => $subject,
					'from_email' => $this->settings->config('sendmail_email'),
					'from_name' => $this->settings->config('company_name'),
					'to' => array(
						array(
							'email' => $email,
							'name' => $email
						)
					),
					'headers' => array('Reply-To' => $this->settings->config('sendmail_email')),
					'important' => false
				);
				$async = true;
				$ip_pool = 'Main Pool';
				$mandrill->messages->send($mandrillMessage, $async, $ip_pool);

			} catch (\Mandrill_Error $e) {

				echo 'A mandrill error occurred: ' . get_class($e) . ' - ' . $e->getMessage();
				throw $e;

			}

		} else if ($this->getDispatchSystem == 'mailgun') {

			list(, $domain) = explode('@', $this->settings->config('sendmail_email'));

			$mail = new \Mailgun\Mailgun($this->settings->config('mailgun_api_key'));
			$mail->sendMessage($domain, array(
				'from' => $this->settings->config('company_name') . ' <' . $this->settings->config('sendmail_email') . '>',
				'to' => $email . ' <' . $email . '>',
				'subject' => $subject,
				'html' => $this->message
			));

		} else if ($this->getDispatchSystem == 'sendgrid') {

			/*
			 * Decrypt Key Information
			 */
			list($iv, $hash) = explode('.', $this->settings->config('sendgrid_api_key'));
			list($username, $password) = explode('|', Components\Authentication::decrypt($hash, $iv));

			$sendgrid = new \SendGrid($username, $password);
			$email = new \SendGrid\Email();

			$email->addTo($email)->
					setFrom($this->settings->config('sendmail_email'))->
					setSubject($subject)->
					setHtml($this->message);

			$sendgrid->send($email);

		} else {

			$headers = 'From: ' . $this->settings->config('sendmail_email') . "\r\n" .
					'Reply-To: ' . $this->settings->config('sendmail_email') . "\r\n" .
					'X-Mailer: PHP/' . phpversion();

			mail($email, $subject, $this->message, $headers);

		}

	}

	/**
	 * Gets the Email System to send with from the settings.
	 *
	 * @return string
	 */
	private function getDispatchSystemFunct() {

		return $this->settings->config('sendmail_method');

	}

	/**
	 * Finds and outputs a given email template.
	 *
	 * @param string $template
	 * @return string
	 */
	private function readTemplate($template) {

		$getTemplate = file_get_contents(APP_DIR . 'templates/email/' . $template . '.tpl');
		if (!$getTemplate) {
			throw new \Exception('Requested template `' . $template . '` could not be found.');
		} else {
			return $getTemplate;
		}

	}

	/**
	 * Generates a Login Notification Email (does not send it).
	 *
	 * @param string $type What type of login notification are we sending?
	 * @param array $vars
	 * @return void
	 * @deprecated This is just a more or less glorified buildEmail
	 */
	public function generateLoginNotification($type, $vars) {

		$find = array('{{ HOST_NAME }}', '{{ IP_ADDRESS }}', '{{ GETHOSTBY_IP_ADDRESS }}', '{{ DATE }}', '{{ MASTER_URL }}');
		$replace = array($this->settings->config('company_name'), $vars['IP_ADDRESS'], $vars['GETHOSTBY_IP_ADDRESS'], date('r', time()), $this->masterurl);

		if ($type == 'failed') {
			
			$this->message = str_replace($find, $replace, $this->readTemplate('login_failed'));
			
		} else if ($type == 'success') {

			$this->message = str_replace($find, $replace, $this->readTemplate('login_success'));

		} else {

			throw new \Exception('Invalid email template specified.');

		}

		return $this;

	}

	/**
	 * Reads an email template and compiles the necessary data into it.
	 *
	 * @param string $tpl The email template to use.
	 * @param array $data
	 * @return void
	 */
	public function buildEmail($tpl, $data = array()) {

		$this->message = str_replace(array('{{ HOST_NAME }}', '{{ MASTER_URL }}', '{{ DATE }}'), array($this->settings->config('company_name'), $this->masterurl, date('j/F/Y H:i', time())), $this->readTemplate($tpl));

		foreach ($data as $key => $val) {
			$this->message = str_replace('{{ ' . $key . ' }}', $val, $this->message);
		}

		return $this;

	}

}

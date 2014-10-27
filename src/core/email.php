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

/**
 * PufferPanel Core Email Sending Class
 */
class Email {

	use Components\Authentication;

	/**
	 * @param string $message
	 */
	private $message;

	/**
	 * Constructor for email sending
	 *
	 * @return void
	 */
	public function __construct()
		{

			$this->mysql = null;
			$this->settings = new Settings();

		}

	/**
	 * Sends an email that has been formatted.
	 *
	 * @param string $email The email address to send to.
	 * @param string $subject The subject of the email.
	 * @return void
	 */
	public function dispatch($email, $subject)
		{

			$this->getDispatchSystem = $this->getDispatchSystemFunct();
			if($this->getDispatchSystem == 'php')
				{

					$headers = 'From: '. $this->settings->get('sendmail_email') . "\r\n" .
					    'Reply-To: '. $this->settings->get('sendmail_email') . "\r\n" .
						'MIME-Version: 1.0' . "\r\n" .
						'Content-type: text/html; charset=iso-8859-1' . "\r\n" .
					    'X-Mailer: PHP/' . phpversion();

					mail($email, $subject, $this->message, $headers);

				}
			else if($this->getDispatchSystem == 'postmark')
				{

					Postmark\Mail::compose($this->settings->get('postmark_api_key'))
					    ->from($this->settings->get('sendmail_email'), $this->settings->get('company_name'))
					    ->addTo($email, $email)
					    ->subject($subject)
					    ->messageHtml($this->message)
					    ->send();

				}
			else if($this->getDispatchSystem == 'mandrill')
				{

					try {

					    $mandrill = new Mandrill($this->settings->get('mandrill_api_key'));
					    $mandrillMessage = array(
					        'html' => $this->message,
					        'subject' => $subject,
					        'from_email' => $this->settings->get('sendmail_email'),
					        'from_name' => $this->settings->get('company_name'),
					        'to' => array(
					            array(
					                'email' => $email,
					                'name' => $email
					            )
					        ),
					        'headers' => array('Reply-To' => $this->settings->get('sendmail_email')),
					        'important' => false
					    );
					    $async = true;
					    $ip_pool = 'Main Pool';
					    $result = $mandrill->messages->send($mandrillMessage, $async, $ip_pool);

					} catch(Mandrill_Error $e) {

					    echo 'A mandrill error occurred: ' . get_class($e) . ' - ' . $e->getMessage();
					    throw $e;

					}

				}
			else if($this->getDispatchSystem == 'mailgun')
				{

					list($name, $domain) = explode('@', $this->settings->get('sendmail_email'));

					$ch = curl_init();
					curl_setopt($ch, CURLOPT_URL, 'https://api.mailgun.net/v2/'.$domain.'/messages');
					curl_setopt($ch, CURLOPT_HTTPAUTH, CURLAUTH_BASIC);
					curl_setopt($ch, CURLOPT_USERPWD, 'api:'.$this->settings->get('mailgun_api_key'));
					curl_setopt($ch, CURLOPT_RETURNTRANSFER, 0);
					curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'POST');
					curl_setopt($ch, CURLOPT_POSTFIELDS, array(
						'from' => $this->settings->get('company_name').' <'.$this->settings->get('sendmail_email').'>',
						'to' => $email.' <'.$email.'>',
						'subject' => $subject,
						'html' => $this->message
					));

					curl_exec($ch);

					curl_close($ch);

				}
			else if($this->getDispatchSystem == 'sendgrid')
				{

					/*
					 * Decrypt Key Information
					 */
					list($iv, $hash) = explode('.', $this->settings->get('sendgrid_api_key'));
					list($username, $password) = explode('|', components::decrypt($hash, $iv));

					$sendgrid = new SendGrid($username, $password);
					$email = new SendGrid\Email();

					$email->addTo($email)->
					       setFrom($this->settings->get('sendmail_email'))->
					       setSubject($subject)->
					       setHtml($this->message);

					$sendgrid->send($email);

				}
			else
				{

					$headers = 'From: '. $this->settings->get('sendmail_email') . "\r\n" .
					    'Reply-To: '. $this->settings->get('sendmail_email') . "\r\n" .
					    'X-Mailer: PHP/' . phpversion();

					mail($email, $subject, $this->message, $headers);

				}

		}

	/**
	 * Gets the Email System to send with from the settings.
	 *
	 * @return string
	 */
	private function getDispatchSystemFunct()
		{

			$this->selectSystem = $this->mysql->prepare("SELECT * FROM `acp_settings` WHERE `setting_ref` = 'sendmail_method'");
			$this->selectSystem->execute();

				$this->selectRow = $this->selectSystem->fetch();

				return $this->selectRow['setting_val'];

		}

	/**
	 * Finds and outputs a given email template.
	 *
	 * @param string $template
	 * @return string
	 */
	private function readTemplate($template)
		{

			$this->getTemplate = @file_get_contents(APP_DIR.'templates/email/'.$template.'.tpl');
			if(!$this->getTemplate)
				die('Requested template `'.$template.'` could not be found.');
			else
				return $this->getTemplate;

		}

	/**
	 * Generates a Login Notification Email (does not send it).
	 *
	 * @param string $type What type of login notification are we sending?
	 * @param array $vars
	 * @return void
	 */
	public function generateLoginNotification($type, $vars)
		{

			if($type == 'failed')
				{

					$this->find = array('{{ HOST_NAME }}', '{{ IP_ADDRESS }}', '{{ GETHOSTBY_IP_ADDRESS }}', '{{ DATE }}', '{{ MASTER_URL }}');
					$this->replace = array($this->settings->get('company_name'), $vars['IP_ADDRESS'], $vars['GETHOSTBY_IP_ADDRESS'], date('r', time()), $this->settings->get('master_url'));

					$this->message = str_replace($this->find, $this->replace, $this->readTemplate('login_failed'));
					return $this;

				}
			else if($type == 'success')
				{

					$this->find = array('{{ HOST_NAME }}', '{{ IP_ADDRESS }}', '{{ GETHOSTBY_IP_ADDRESS }}', '{{ DATE }}', '{{ MASTER_URL }}');
					$this->replace = array($this->settings->get('company_name'), $vars['IP_ADDRESS'], $vars['GETHOSTBY_IP_ADDRESS'], date('r', time()), $this->settings->get('master_url'));

					$this->message = str_replace($this->find, $this->replace, $this->readTemplate('login_success'));
					return $this;

				}
			else
				{

					die('Invalid email template specified.');

				}

		}

	/**
	 * Reads an email template and compiles the necessary data into it.
	 *
	 * @param string $tpl The email template to use.
	 * @param array $data
	 * @return void
	 */
	public function buildEmail($tpl, $data = array())
		{

			$this->message = $this->readTemplate($tpl);
			$this->message = str_replace(array('{{ HOST_NAME }}', '{{ MASTER_URL }}', '{{ DATE }}'), array($this->settings->get('company_name'), $this->settings->get('master_url'), date('j/F/Y H:i', time())), $this->message);

				foreach($data as $key => $val)
					$this->message  = str_replace('{{ '.$key.' }}', $val, $this->message);

				return $this;

		}

}

?>

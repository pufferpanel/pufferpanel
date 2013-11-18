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
 * Core Email Sending Class
 */

class tplMail extends dbConn {

	public function __construct($settings)
		{
	
			$this->mysql = parent::getConnection();
			$this->settings = $settings;
	
		}
		
	/*
	 * Email Sending
	 */
	public function dispatch($email, $subject, $message)
		{
		
			$this->getDispatchSystem = $this->getDispatchSystemFunct();
			if($this->getDispatchSystem == 'php')
				{
				
					$headers = 'From: '. $this->settings->get('sendmail_email') . "\r\n" .
					    'Reply-To: '. $this->settings->get('sendmail_email') . "\r\n" .
					    'X-Mailer: PHP/' . phpversion();
	
					mail($email, $subject, $message, $headers);

				}
			else if($this->getDispatchSystem == 'postmark')
				{
				
					include('postmark/Mail.php');
					Postmark\Mail::compose($this->settings->get('postmark_api_key'))
					    ->from($this->settings->get('sendmail_email'), $this->settings->get('company_name'))
					    ->addTo($email, $email)
					    ->subject($subject)
					    ->messageHtml($message)
					    ->send();
				
				}
			else if($this->getDispatchSystem == 'mandrill')
				{
				
					include('mandrill/Mandrill.php');
					
					try {
					
					    $mandrill = new Mandrill($this->settings->get('mandrill_api_key'));
					    $mandrillMessage = array(
					        'html' => $message,
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
						'html' => $message
					));
					
					curl_exec($ch);
									
					curl_close($ch);
					
				}
			else 
				{
			
					$headers = 'From: '. $this->settings->get('sendmail_email') . "\r\n" .
					    'Reply-To: '. $this->settings->get('sendmail_email') . "\r\n" .
					    'X-Mailer: PHP/' . phpversion();
					
					mail($email, $subject, $message, $headers);
			
				}
		
		}
		
	private function getDispatchSystemFunct()
		{
		
			$this->selectSystem = $this->mysql->prepare("SELECT * FROM `acp_settings` WHERE `setting_ref` = 'sendmail_method'");
			$this->selectSystem->execute();
			
				$this->selectRow = $this->selectSystem->fetch();
				
				return $this->selectRow['setting_val'];
		
		}
	
	/*
	 * Email Creation & Templates
	 */

	private function readTemplate($template)
		{
		
			$this->tpl = $this->mysql->prepare("SELECT * FROM	`acp_email_templates` WHERE `tpl_name` = ?");
			$this->tpl->execute(array($template));
			
			if($this->tpl->rowCount() == 1)
				{
				
					$row = $this->tpl->fetch();
					return $row['tpl_content'];
				
				}
			else 
				{
				
					die('Requested template `'.$tpl.'` could not be found.');
					
				}
			
		
		}
	
	/*
	 * Generate the Email Message from the Templates for Login Notifications.
	 * This does not send the email, simply generates the email message.
	 */
	public function generateLoginNotification($type, $vars)
		{

			if($type == 'failed')
				{
								
					$this->find = array('<%HOST_NAME%>', '<%IP_ADDRESS%>', '<%GETHOSTBY_IP_ADDRESS%>', '<%DATE%>', '<%MASTER_URL%>');
					$this->replace = array($this->settings->get('company_name'), $vars['IP_ADDRESS'], $vars['GETHOSTBY_IP_ADDRESS'], date('r', time()), $this->settings->get('master_url'));
						
					return str_replace($this->find, $this->replace, $this->readTemplate('login_failed'));
				
				}
			else if($type == 'success')
				{
				
					$this->find = array('<%HOST_NAME%>', '<%IP_ADDRESS%>', '<%GETHOSTBY_IP_ADDRESS%>', '<%DATE%>', '<%MASTER_URL%>');
					$this->replace = array($this->settings->get('company_name'), $vars['IP_ADDRESS'], $vars['GETHOSTBY_IP_ADDRESS'], date('r', time()), $this->settings->get('master_url'));
						
					return str_replace($this->find, $this->replace, $this->readTemplate('login_success'));
				
				}
			else 
				{
				
					die('Invalid email template specified.');
					
				}
		
		}
		
	public function generateEmailChangedNotification($vars)
		{
		
			$this->find = array('<%HOST_NAME%>', '<%EMAIL_KEY%>', '<%IP_ADDRESS%>', '<%GETHOSTBY_IP_ADDRESS%>', '<%DATE%>', '<%MASTER_URL%>');
			$this->replace = array($this->settings->get('company_name'), $vars['EMAIL_KEY'], $vars['IP_ADDRESS'], $vars['GETHOSTBY_IP_ADDRESS'], date('r', time()), $this->settings->get('master_url'));
				
			return str_replace($this->find, $this->replace, $this->readTemplate('email_changed'));
		
		}
		
	public function generatePasswordChangedNotification($vars)
		{
		
			$this->find = array('<%HOST_NAME%>', '<%IP_ADDRESS%>', '<%GETHOSTBY_IP_ADDRESS%>', '<%DATE%>');
			$this->replace = array($this->settings->get('company_name'), $vars['IP_ADDRESS'], $vars['GETHOSTBY_IP_ADDRESS'], date('r', time()));
				
			return str_replace($this->find, $this->replace, $this->readTemplate('password_changed'));
		
		}
	
	public function generateForgottenPasswordEmail($vars)
		{
		
			$this->find = array('<%HOST_NAME%>', '<%PKEY%>', '<%IP_ADDRESS%>', '<%GETHOSTBY_IP_ADDRESS%>', '<%DATE%>', '<%MASTER_URL%>');
			$this->replace = array($this->settings->get('company_name'), $vars['PKEY'], $vars['IP_ADDRESS'], $vars['GETHOSTBY_IP_ADDRESS'], date('r', time()), $this->settings->get('master_url'));
				
			return str_replace($this->find, $this->replace, $this->readTemplate('password_reset'));
		
		}
		
	public function generateNewPasswordEmail($vars)
		{
		
			$this->find = array('<%HOST_NAME%>', '<%EMAIL%>', '<%NEW_PASS%>', '<%MASTER_URL%>');
			$this->replace = array($this->settings->get('company_name'), $vars['EMAIL'], $vars['NEW_PASS'], $this->settings->get('master_url'));
				
			return str_replace($this->find, $this->replace, $this->readTemplate('new_password'));
		
		}
		
	public function adminAccountCreated($vars)
		{
		
			$this->find = array('<%HOST_NAME%>', '<%EMAIL%>', '<%PASS%>', '<%MASTER_URL%>');
			$this->replace = array($this->settings->get('company_name'), $vars['EMAIL'], $vars['PASS'], $this->settings->get('master_url'));
				
			return str_replace($this->find, $this->replace, $this->readTemplate('admin_newaccount'));
		
		}
    
    public function generateSFTPPasswordUpdateEmail($vars)
        {
            
            $this->find = array('<%HOST_NAME%>', '<%PASS%>', '<%SERVER%>');
			$this->replace = array($this->settings->get('company_name'), $vars['PASS'], $vars['SERVER']);
				
			return str_replace($this->find, $this->replace, $this->readTemplate('admin_new_sftppass'));
        
        }

}

?>
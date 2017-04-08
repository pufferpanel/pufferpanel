<?php

/*
  PufferPanel - A Game Server Management Panel
  Copyright (c) 2015 Dane Everitt

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

use \Tracy\Debugger,
    \Exception,
    \SendGrid;

/**
 * PufferPanel Core Email Sending Class
 */
class Email {

    use Components\Authentication;

    /**
     * @param string $master_url
     */
    protected $master_url;

    /**
     * @param string $message
     */
    protected $message;

    /**
     * @param string $email
     */
    protected $email;

    /**
     * @param string $subject
     */
    protected $subject;

    /**
     * Sends an email that has been formatted using the method defined in settings.
     *
     * @param string $email The email address to send to.
     * @param string $subject The subject of the email.
     * @return void
     */
    public function dispatch($email, $subject) {

        $this->email = $email;
        $this->subject = $subject;

        switch (Settings::config()->transport_method) {

            case 'postmark':
                $this->_sendWithPostmark();
                break;
            case 'mandrill':
                $this->_sendWithMandrill();
                break;
            case 'mailgun':
                $this->_sendWithMailgun();
                break;
            case 'sendgrid':
                $this->_sendWithSendgrid();
                break;
            default:
                $this->_sendWithPHP();
        }
    }

    /**
     * Sends an email using the built-in PHP mail() function.
     *
     * @return void
     */
    protected function _sendWithPHP() {

        $headers = 'From: ' . Settings::config()->transport_email . "\r\n" .
                'Reply-To: ' . Settings::config()->transport_email . "\r\n" .
                'MIME-Version: 1.0' . "\r\n" .
                'Content-type: text/html; charset=iso-8859-1' . "\r\n" .
                'X-Mailer: PHP/' . phpversion();

        try {
            mail($this->email, $this->subject, $this->message, $headers);
        } catch (\Exception $e) {
            Debugger::log($e);
        }
    }

    /**
     * Sends an email using the Sendgrid Email API.
     *
     * @return void
     */
    protected function _sendWithSendgrid() {

        try {
            $sendgrid = new SendGrid(Settings::config()->transport_token);
            $to = new SendGrid\Email(null, $this->email);
            $from = new SendGrid\Email(null, Settings::config()->transport_email);
            $subject = $this->subject;
            $body = new SendGrid\Content('text/html', $this->message);

            $request = new SendGrid\Mail($from, $subject, $to, $body);

            $sendgrid->client->mail()->send()->post($request);
        } catch (\Exception $e) {
            Debugger::log($e);
        }
    }

    /**
     * Sends an email using the Postmark Email API.
     *
     * @return void
     */
    protected function _sendWithPostmark() {

        try {

            $client = new \Postmark\PostmarkClient(Settings::config()->transport_token);
            $client->sendEmail(
                    Settings::config()->transport_email, $this->email, $this->subject, $this->message
            );
        } catch (\Exception $e) {
            Debugger::log($e);
        }
    }

    /**
     * Sends an email using the Mailgun Email API.
     *
     * @return void
     */
    protected function _sendWithMailgun() {

        try {

            list(, $domain) = explode('@', Settings::config()->transport_email);

            $mail = new \Mailgun\Mailgun(Settings::config()->transport_token);
            $mail->sendMessage($domain, array(
                'from' => Settings::config()->company_name . ' <' . Settings::config()->transport_email . '>',
                'to' => $this->email . ' <' . $this->email . '>',
                'subject' => $this->subject,
                'html' => $this->message
            ));
        } catch (\Exception $e) {
            Debugger::log($e);
        }
    }

    /**
     * Sends an email using the Mandrill Email API.
     *
     * @return void
     */
    protected function _sendWithMandrill() {

        try {

            $mandrill = new \Mandrill(Settings::config()->transport_token);
            $mandrill->messages->send(array(
                'html' => $this->message,
                'subject' => $this->subject,
                'from_email' => Settings::config()->transport_email,
                'from_name' => Settings::config()->company_name,
                'to' => array(
                    array(
                        'email' => $this->email,
                        'name' => $this->email
                    )
                ),
                'headers' => array('Reply-To' => Settings::config()->transport_email),
                'important' => false
                    ), true, 'Main Pool');
        } catch (\Exception $e) {
            Debugger::log($e);
        }
    }

    /**
     * Finds and outputs a given email template.
     *
     * @param string $template
     * @return string
     */
    protected function _readTemplate($template) {

        $readTemplate = file_get_contents(APP_DIR . 'templates/email/' . $template . '.tpl');
        if (!$readTemplate) {
            throw new Exception('Requested template `' . $readTemplate . '` could not be found.');
        } else {
            return $readTemplate;
        }
    }

    /**
     * Generates a Login Notification Email (does not send it).
     *
     * @param string $type What type of login notification are we sending?
     * @param array $vars
     * @return Email
     * @deprecated This is just a more or less glorified buildEmail
     */
    public function generateLoginNotification($type, $vars) {

        $find = array('{{ HOST_NAME }}', '{{ IP_ADDRESS }}', '{{ GETHOSTBY_IP_ADDRESS }}', '{{ DATE }}', '{{ MASTER_URL }}');
        $replace = array(Settings::config()->company_name, $vars['IP_ADDRESS'], $vars['GETHOSTBY_IP_ADDRESS'], date('r', time()), Settings::config()->master_url);

        $this->message = str_replace($find, $replace, $this->_readTemplate('login_' . $type));

        return $this;
    }

    /**
     * Reads an email template and compiles the necessary data into it.
     *
     * @param string $tpl The email template to use.
     * @param array $data
     * @return Email
     */
    public function buildEmail($tpl, array $data) {

        $this->message = str_replace(array('{{ HOST_NAME }}', '{{ MASTER_URL }}', '{{ DATE }}'), array(Settings::config()->company_name, Settings::config()->master_url, date('j/F/Y H:i', time())), $this->_readTemplate($tpl));

        foreach ($data as $key => $val) {
            $this->message = str_replace('{{ ' . $key . ' }}', $val, $this->message);
        }

        return $this;
    }

}

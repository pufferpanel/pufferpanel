<?php

namespace PufferPanel\Core;

use \Klein\Response;

class PPResponse extends Response
{
    private $klein;
    private $altered = null;

    public function __construct($klein, $body = '', $status_code = null, array $headers = array())
    {
        parent::__construct($body, $status_code, $headers);
        $this->klein = $klein;
        $this->altered = false;
    }


    public function abort($code) {
        $this->setAltered();
        $this->klein->abort($code);
    }

    public function body($body = null)
    {
        if(!is_null($body)) {
            $this->setAltered();
        }
        return parent::body($body);
    }

    public function code($code = null)
    {
        if(!is_null($code)) {
            $this->setAltered();
        }
        return parent::code($code);
    }

    public function redirect($url, $code = 302)
    {
        $this->setAltered();
        if (strpos($url,'http') !== false)
            parent::redirect($url, $code);
        else
            parent::redirect(BASE_URL.$url, $code);
        $this->abort($code);
    }

    private function setAltered() {
        if (is_null($this->altered)) {
            return;
        }
        $this->altered = true;
    }

    public function isAltered() {
        return $this->altered;
    }
}
<?php

class getSettings extends dbConn {

	public function __construct()
		{
		
			$this->mysql = parent::getConnection();
		
			$this->query = $this->mysql->prepare("SELECT * FROM `acp_settings`");
			$this->query->execute();
				
				while($this->row = $this->query->fetch()){
					
					$this->_data[$this->row['setting_ref']] = $this->row['setting_val'];
					
				}
		
		}
		
	public function get($setting)
		{
		
			return (array_key_exists($setting, $this->_data)) ? $this->_data[$setting] : '_notfound_';
		
		}

}

?>
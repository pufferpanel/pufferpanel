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
 * PufferPanel Folder Management Class
 */

class files {

	public function formatSize($bytes, $decimals = 2) {

		  $sz = explode(',', 'B,KB,MB,GB');
		  $factor = floor((strlen($bytes) - 1) / 3);
		  
		  return sprintf("%.{$decimals}f", $bytes / pow(1024, $factor)) . ' '.$sz[$factor];
	    
	}
	
	public function format($size, $precision = 0){
	
		$base = log($size) / log(1024);
	
	    return round(pow(1024, $base - floor($base)), $precision);
	    
	}
	
	public function readLines($filename, $lines)
	{

		$r = '';
		$file = file($filename);
		
			if(count($file) < $lines){
			
				$lines = count($file);
				
			}
		
				for ($i = count($file)-$lines; $i < count($file); $i++) {
		  
		 			 $r .= $file[$i];
		
				}
				
		return $r;
	      
	}
	
	function last_lines($path, $line_count, $block_size = 512){
	    
	    if(!is_file($path)) {
	    	
	    	echo('Could not locate server.log!');
	    	return false;
	    	
	    }else{
	    
		    $lines = array();
			    $leftover = "";
		    $fh = fopen($path, 'r');
		    fseek($fh, 0, SEEK_END);
		    
		    do{
	
		        $can_read = $block_size;
		        if(ftell($fh) < $block_size){
		            $can_read = ftell($fh);
		        }
		
		        fseek($fh, -$can_read, SEEK_CUR);
		        $data = fread($fh, $can_read);
		        $data .= $leftover;
		        fseek($fh, -$can_read, SEEK_CUR);
		
		        $split_data = array_reverse(explode("\n", $data));
		        $new_lines = array_slice($split_data, 0, -1);
		        $lines = array_merge($lines, $new_lines);
		        $leftover = $split_data[count($split_data) - 1];
		        
		    }
		    while(count($lines) < $line_count && ftell($fh) != 0);
		    
		    if(ftell($fh) == 0){
		        $lines[] = $leftover;
		    }
		    
		    fclose($fh);
	
			/*
			 * Return in Readbale Format
			 */
			$a = array_reverse(array_slice($lines, 0, $line_count));
		    
		    	$output = '';
		    	$total = count($a);
		    	foreach($a as $id => $line){
		    	
		    		if($id == ($total - 2))
		    			$output .= $line;
					else if($id == ($total - 1))
						$output .= '';
		    		else
		    			$output .= $line."\n";
		    	
		    	}	    	
		    
		    return $output;
		    
		}
	    
	}
	
	function download($filename) { 
		
		$chunksize = 1*(1024*1024);
		$buffer = '';
		$handle = fopen($filename, 'rb');
		
			if ($handle === false){
			
				return false;
				
			}

			while (!feof($handle)){
			
				$buffer = fread($handle, $chunksize);
				print $buffer;
				
			}
			
		return fclose($handle);
		
	} 
	
}

?>
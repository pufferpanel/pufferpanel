<?php
class MinecraftRconException extends Exception
{
	// Exception thrown by MinecraftRcon class
}

class MinecraftRcon
{
	/*
	 * Class written by xPaw
	 *
	 * Website: http://xpaw.ru
	 * GitHub: https://github.com/xPaw/PHP-Minecraft-Query
	 *
	 * Protocol: https://developer.valvesoftware.com/wiki/Source_RCON_Protocol
	 */
	
	// Sending
	const SERVERDATA_EXECCOMMAND    = 2;
	const SERVERDATA_AUTH           = 3;
	
	// Receiving
	const SERVERDATA_RESPONSE_VALUE = 0;
	const SERVERDATA_AUTH_RESPONSE  = 2;
	
	private $Socket;
	private $RequestId;
	
	public function __destruct( )
	{
		$this->Disconnect( );
	}
	
	public function Connect( $Ip, $Port = 25575, $Password, $Timeout = 3 )
	{
		$this->RequestId = 0;
		
		if( $this->Socket = FSockOpen( $Ip, (int)$Port ) )
		{
			Socket_Set_TimeOut( $this->Socket, $Timeout );
			
			if( !$this->Auth( $Password ) )
			{
				$this->Disconnect( );
				
				throw new MinecraftRconException( "Authorization failed." );
			}
		}
		else
		{
			throw new MinecraftQueryException( "Can't open socket." );
		}
	}
	
	public function Disconnect( )
	{
		if( $this->Socket )
		{
			FClose( $this->Socket );
			
			$this->Socket = null;
		}
	}
	
	public function Command( $String )
	{
		if( !$this->WriteData( self :: SERVERDATA_EXECCOMMAND, $String ) )
		{
			return false;
		}
		
		$Data = $this->ReadData( );
		
		if( $Data[ 'RequestId' ] < 1 || $Data[ 'Response' ] != self :: SERVERDATA_RESPONSE_VALUE )
		{
			return false;
		}
		
		return $Data[ 'String' ];
	}
	
	private function Auth( $Password )
	{
		if( !$this->WriteData( self :: SERVERDATA_AUTH, $Password ) )
		{
			return false;
		}
		
		$Data = $this->ReadData( );
		
		return $Data[ 'RequestId' ] > -1 && $Data[ 'Response' ] == self :: SERVERDATA_AUTH_RESPONSE;
	}
	
	private function ReadData( )
	{
		$Packet = Array( );
		
		$Size = FRead( $this->Socket, 4 );
		$Size = UnPack( 'V1Size', $Size );
		$Size = $Size[ 'Size' ];
		
		// TODO: Add multiple packets (Source)
		
		$Packet = FRead( $this->Socket, $Size );
		$Packet = UnPack( 'V1RequestId/V1Response/a*String/a*String2', $Packet );
		
		return $Packet;
	}
	
	private function WriteData( $Command, $String = "" )
	{
		// Pack the packet together
		$Data = Pack( 'VV', $this->RequestId++, $Command ) . $String . "\x00\x00\x00"; 
		
		// Prepend packet length
		$Data = Pack( 'V', StrLen( $Data ) ) . $Data;
		
		$Length = StrLen( $Data );
		
		return $Length === FWrite( $this->Socket, $Data, $Length );
	}
}
?>
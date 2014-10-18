<html>
	<head>
		<title>{{ HOST_NAME }} - Added to Server</title>
	</head>
	<body>
		<center><h1>{{ HOST_NAME }} - Added to Server</h1></center>
		<p>Hello there! This email is to inform you that you have been invited to help manage the following server: <i>{{ SERVER }}</i>.</p>
		<p>If you already have an account please <a href="{{ MASTER_URL }}">login as normal</a>, and then go to <strong>Account Settings &rarr; Server Token</strong> and input the token below.
			If you do not have an account please <a href="{{ MASTER_URL }}register.php?token={{ TOKEN }}">click here to create an account</a>.</p>
		<p><b>Token:</b> {{ TOKEN }}</p>
		<p>Thanks!<br />{{ HOST_NAME }}</p>
	</body>
</html>
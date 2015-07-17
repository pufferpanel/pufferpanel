<html>
	<head>
		<title>{{ HOST_NAME }} - New Server Added</title>
	</head>
	<body>
		<center><h1>{{ HOST_NAME }} - New Server Added </h1></center>
		<p>Hello there! This email is to inform you that a new server ({{ NAME }}) has been created for you.</p>
		<p><strong>SFTP Connection:</strong> {{ SFTP }}<br />
		<p><strong>SFTP Username:</strong> {{ USER }}<br />
		<p><strong>SFTP Password:</strong> {{ PASS }}<br />
		<p>You can connect to your server after starting it from the panel with the following connection information: <strong>{{ SERVER_CONN }}</strong>.
		<p>Thanks!<br />{{ HOST_NAME }}</p>
	</body>
</html>
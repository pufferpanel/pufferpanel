<html>
	<head>
		<title>{{ HOST_NAME }} - New Server Added</title>
	</head>
	<body>
		<center><h1>{{ HOST_NAME }} - New Server Added </h1></center>
		<p>Hello there! This email is to inform you that a new server ({{ NAME }}) has been created for you.</p>
		<p><strong>FTP Connection:</strong> {{ FTP }}<br />
		<p><strong>FTP Username:</strong> {{ USER }}<br />
		<p><strong>FTP Password:</strong> {{ PASS }}<br />
		<p>You can connect to Minecraft after starting your server from the panel with the following connection information: <strong>{{ MINECRAFT }}</strong>.
		<p>Thanks!<br />{{ HOST_NAME }}</p>
	</body>
</html>
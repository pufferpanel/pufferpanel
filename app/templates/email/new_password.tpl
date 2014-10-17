<html>
	<head>
		<title>{{ HOST_NAME }} - New Password</title>
	</head>
	<body>
		<center><h1>{{ HOST_NAME }} - New Password</h1></center>
		<p>Hello there! You are receiving this email because you requested a new password for your {{ HOST_NAME }} account.</p>
		<p><strong>Login:</strong> <a href="{{ MASTER_URL }}">{{ MASTER_URL }}</a><br />
			<strong>Email:</strong> {{ EMAIL }}<br />
			<strong>Password:</strong> {{ NEW_PASS }}</p>
		<p>Thanks!<br />{{ HOST_NAME }}</p>
	</body>
</html>
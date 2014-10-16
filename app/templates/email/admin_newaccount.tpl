<html>
	<head>
		<title>{{ HOST_NAME }} - Account Created</title>
	</head>
	<body>
		<center><h1>{{ HOST_NAME }} - Account Created</h1></center>
		<p>Hello there! This email is to inform you that an account has been created for you on {{ HOST_NAME }}.</p>
		<p><strong>Login:</strong> <a href="{{ MASTER_URL }}">{{ MASTER_URL }}</a><br />
			<strong>Email:</strong> {{ EMAIL }}<br />
			<strong>Password:</strong> {{ PASS }}</p>
		<p>Thanks!<br />{{ HOST_NAME }}</p>
	</body>
</html>
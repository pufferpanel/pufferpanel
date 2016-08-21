<html>
	<head>
		<title>{{ HOST_NAME }} - Account Login Notification</title>
	</head>
	<body>
		<center><h1>{{ HOST_NAME }} - Account Login Notification</h1></center>
		<p>You are recieving this email as a part of our continuing efforts to improve our server security. You are recieving this email as a sucessful login was made with your account on {{ HOST_NAME }}.</strong></p>
		<p><strong>IP Address:</strong> {{ IP_ADDRESS }}<br />
			<strong>Hostname:</strong> {{ GETHOSTBY_IP_ADDRESS }}<br />
			<strong>Time:</strong> {{ DATE }}</p>
		<p>This email is intended to keep you aware of any possible malicious account activity. You can change your notification preferences by <a href="{{ MASTER_URL }}account">clicking here</a>.</p>
		<p>Thanks!<br />{{ HOST_NAME }}</p>
	</body>
</html>
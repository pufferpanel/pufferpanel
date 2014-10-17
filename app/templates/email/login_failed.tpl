<html>
	<head>
		<title>{{ HOST_NAME }} - Account Login Failure Notification</title>
	</head>
	<body>
		<center><h1>{{ HOST_NAME }} - Account Login Failure Notification</h1></center>
		<p>You are recieving this email as a part of our continuing efforts to improve our server security. <strong>An unsucessful login was made with your account on {{ HOST_NAME }}.</strong></p>
		<p><strong>IP Address:</strong> {{ IP_ADDRESS }}<br />
			<strong>Hostname:</strong> {{ GETHOSTBY_IP_ADDRESS }}<br />
			<strong>Time:</strong> {{ DATE }}</p>
		<p>At this time your account is still safe and sound in our system. This email is simply to let you know that someone tried to login to your account and failed. You can change your notification preferences by <a href="{{ MASTER_URL }}accounts.php">clicking here</a>.</p>
		<p>Thanks!<br />{{ HOST_NAME }}</p>
	</body>
</html>
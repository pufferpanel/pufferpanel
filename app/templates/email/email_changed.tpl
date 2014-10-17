<html>
	<head>
		<title>{{ HOST_NAME }} Email Change Notification</title>
	</head>
	<body>
		<center><h1>{{ HOST_NAME }} Email Change Notification</h1></center>
		<p>Hello there! You are receiving this email because you requested to change your account email on {{ HOST_NAME }}.</p>
		<p>Please click the link below to confirm that you wish to use this email address on {{ HOST_NAME }}. If you did not make this request, or do not wish to continue simply ignore this email and nothing will happen. <strong>This link will expire in 4 hours.</strong></p>
		<p><a href="{{ MASTER_URL }}account_actions.php?conf=email&key={{ EMAIL_KEY }}">{{ MASTER_URL }}account.php?conf=email&key={{ EMAIL_KEY }}</a></p>
		<p>This change was requested from {{ IP_ADDRESS }} ({{ GETHOSTBY_IP_ADDRESS }}) on {{ DATE }}. Please do not hesitate to contact us if you belive something is wrong.
		<p>Thanks!<br />{{ HOST_NAME }}</p>
	</body>
</html>
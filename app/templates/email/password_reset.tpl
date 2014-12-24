<html>
	<head>
		<title>{{ HOST_NAME }} Lost Password Recovery</title>
	</head>
	<body>
		<center><h1>{{ HOST_NAME }} Lost Password Recovery</h1></center>
		<p>Hello there! You are receiving this email because you requested a new password for your {{ HOST_NAME }} account.</p>
		<p>Please click the link below to confirm that you wish to change your password. If you did not make this request, or do not wish to continue simply ignore this email and nothing will happen. <strong>This link will expire in 4 hours.</strong></p>
		<p><a href="{{ MASTER_URL }}password/verify/{{ PKEY }}">{{ MASTER_URL }}password/verify/{{ PKEY }}</a></p>
		<p>This change was requested from {{ IP_ADDRESS }} ({{ GETHOSTBY_IP_ADDRESS }}) on {{ DATE }}. Please do not hesitate to contact us if you belive something is wrong.
		<p>Thanks!<br />{{ HOST_NAME }}</p>
	</body>
</html>
<html>
    <head>
        <title>{{ HOST_NAME }} - Admin Reset Password</title>
    </head>
    <body>
    <center><h1>{{ HOST_NAME }} - Admin Reset Password</h1></center>
    <p>Hello there! You are receiving this email because an admin has reset the password on your {{ HOST_NAME }} account.</p>
    <p><strong>Login:</strong> <a href="{{ MASTER_URL }}auth/login">{{ MASTER_URL }}auth/login</a><br />
        <strong>Email:</strong> {{ EMAIL }}<br />
        <strong>Password:</strong> {{ NEW_PASS }}</p>
    <p>Thanks!<br />{{ HOST_NAME }}</p>
</body>
</html>
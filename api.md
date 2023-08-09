### Data format

The API works off of the JSON spec. Most of the API expects and returns JSON data. There is no support to change the data format used.

### Successful calls

Successful calls will be any result which produce a 2xx result. This can include 200 (OK), 202 (Accepted), or 204 (No Content). The API will document which status code is returned. 

### User errors

Client errors are a result which stems from a bad user request, whether it's a bad URL or bad data. These will be a 400 (Bad Request), 401 (Unauthorized), 403 (Forbidden), or 404 (Not Found). 

`400`: Refer to the HTTP body for why the call failed.

`401`: Refer to the [Authentication session](#authentication) for information on how to authenticate. 

`403`: Your authorization does not permit accessing this resource. Please check the credentials you are using.

`404`: The URL you are calling does not exist. Please check the full URL.

### Server errors

A 500 may be returned, indicating something went wrong with thr request. This generally indicates something is not as expected, but your request was valid. Please refer to the panel logs and try your request again.

Errors will be returned in the following JSON structure:
```json
{
  "msg": "A human-legible error message, with optional metadata to fill in the message more, such as 'Hello {foo}' would be 'Hello World'",
  "code": "ErrMachineGoBrr",
  "metadata": {
    "foo": "World"
  }
}
```

### Authentication

PufferPanel uses the OAuth2 standard to handle both authentication and authorization of resources. To access the API, a user must first register an OAuth2 client within the panel. A client may be registered within a server's settings, which grants the client only to that server, or to the account, which grants the client access as though it was that user.

To authenticate, refer to the docs specified in the [/oauth2/token](#/default/post_oauth2_token) endpoint on how to request your access token. This follows the OAuth2 spec, which you may read about here: https://oauth.net/2/grant-types/client-credentials/

Once you have a token, you pass the following Header to all API calls. `Authorization: Bearer <token>`

Failure to pass this token when calling the API will result in a 403.
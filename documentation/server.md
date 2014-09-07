server
===============

PufferPanel Core Server management class.




* Class name: server
* Namespace: 
* Parent class: [user](user.md)





Properties
----------


### $_data

```
private mixed $_data
```





* Visibility: **private**


### $_ndata

```
private mixed $_ndata = array()
```





* Visibility: **private**


### $_s

```
private mixed $_s = true
```





* Visibility: **private**


### $_n

```
private mixed $_n = true
```





* Visibility: **private**


### $_l

```
private mixed $_l
```





* Visibility: **private**


### $db

```
protected mixed $db
```





* Visibility: **protected**
* This property is **static**.


### $salt

```
public mixed $salt
```





* Visibility: **public**
* This property is **static**.


Methods
-------


### \server::__construct()

```
void server::\server::__construct()(string $hash, integer $userid, integer $isroot)
```

Constructor class for building server data.



* Visibility: **public**

#### Arguments

* $hash **string** - &lt;p&gt;The server hash.&lt;/p&gt;
* $userid **integer** - &lt;p&gt;The ID of the user who is requesting the server information.&lt;/p&gt;
* $isroot **integer** - &lt;p&gt;The root administrator status of the user requesting the server information.&lt;/p&gt;



### \server::rebuildData()

```
void server::\server::rebuildData()(integer $id)
```

Re-runs the __construct() class with a defined ID for the admin control panel.



* Visibility: **public**

#### Arguments

* $id **integer** - &lt;p&gt;This value should be the ID of the server you are getting information for.&lt;/p&gt;



### \server::getData()

```
string|array|boolean server::\server::getData()(string $id)
```

Provides the corresponding value for the id provided from the MySQL Database.



* Visibility: **public**

#### Arguments

* $id **string** - &lt;p&gt;The column value for the data you need (e.g. server_name).&lt;/p&gt;



### \server::nodeData()

```
string|array|boolean server::\server::nodeData()(string|\nukk $id)
```

Returns data about the node in which the server selected is running.



* Visibility: **public**

#### Arguments

* $id **string|nukk** - &lt;p&gt;The column value for the data you need (e.g. sftp_ip).&lt;/p&gt;



### \server::nodeRedirect()

```
void server::\server::nodeRedirect()($hash, $userid, $rootAdmin)
```

Handles incoming requests to access a server and redirects to the correct location and sets a cookie.



* Visibility: **public**

#### Arguments

* $hash **mixed**
* $userid **mixed**
* $rootAdmin **mixed**



### \server::_rebuildData()

```
array|boolean server::\server::_rebuildData()(integer $userid)
```

Rebuilds server data using a specified ID. Useful for Admin CP applications.



* Visibility: **private**

#### Arguments

* $userid **integer** - &lt;p&gt;The server ID.&lt;/p&gt;



### \server::_buildData()

```
array|boolean server::\server::_buildData()(string $hash, integer $userid, integer $isroot)
```

Builds server data using a specified ID, Hash, and Root Administrator Status.



* Visibility: **private**

#### Arguments

* $hash **string** - &lt;p&gt;The server hash.&lt;/p&gt;
* $userid **integer** - &lt;p&gt;The ID of the user who is requesting the server information.&lt;/p&gt;
* $isroot **integer** - &lt;p&gt;The root administrator status of the user requesting the server information.&lt;/p&gt;



### \Page\components::redirect()

```
mixed server::\Page\components::redirect()($url)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $url **mixed**



### \Page\components::genRedirect()

```
mixed server::\Page\components::genRedirect()()
```





* Visibility: **public**
* This method is **static**.



### \Page\components::twigGET()

```
mixed server::\Page\components::twigGET()()
```





* Visibility: **public**
* This method is **static**.



### \Auth\auth::__construct()

```
mixed server::\Auth\auth::__construct()()
```





* Visibility: **public**



### \user::rebuildData()

```
void server::\user::rebuildData()(string $id)
```

Re-runs the __construct() class with a defined ID for the admin control panel.



* Visibility: **public**

#### Arguments

* $id **string** - &lt;p&gt;The ID of a user requested in the Admin CP.&lt;/p&gt;



### \user::getData()

```
string|array|boolean server::\user::getData()(string|null $id)
```

Provides the corresponding value for the id provided from the MySQL Database.



* Visibility: **public**

#### Arguments

* $id **string|null** - &lt;p&gt;The column value for the data you need (e.g. email).&lt;/p&gt;



### \Auth\auth::validateTOTP()

```
mixed server::\Auth\auth::validateTOTP()($key, $secret)
```





* Visibility: **public**

#### Arguments

* $key **mixed**
* $secret **mixed**



### \Auth\auth::verifyPassword()

```
mixed server::\Auth\auth::verifyPassword()($email, $raw)
```





* Visibility: **public**

#### Arguments

* $email **mixed**
* $raw **mixed**



### \Auth\auth::isLoggedIn()

```
mixed server::\Auth\auth::isLoggedIn()($ip, $session, $serverhash, $acp)
```





* Visibility: **public**

#### Arguments

* $ip **mixed**
* $session **mixed**
* $serverhash **mixed**
* $acp **mixed**



### \Database\database::buildConnection()

```
mixed server::\Database\database::buildConnection()()
```





* Visibility: **public**
* This method is **static**.



### \Database\database::connect()

```
mixed server::\Database\database::connect()()
```





* Visibility: **public**
* This method is **static**.



### \Auth\components::hash()

```
mixed server::\Auth\components::hash()($raw)
```





* Visibility: **public**

#### Arguments

* $raw **mixed**



### \Auth\components::password_compare()

```
mixed server::\Auth\components::password_compare()($raw, $hashed)
```





* Visibility: **private**

#### Arguments

* $raw **mixed**
* $hashed **mixed**



### \Auth\components::generate_iv()

```
mixed server::\Auth\components::generate_iv()()
```





* Visibility: **public**



### \Auth\components::encrypt()

```
mixed server::\Auth\components::encrypt()($raw, $iv, $method)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $raw **mixed**
* $iv **mixed**
* $method **mixed**



### \Auth\components::decrypt()

```
mixed server::\Auth\components::decrypt()($encrypted, $iv, $method)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $encrypted **mixed**
* $iv **mixed**
* $method **mixed**



### \Auth\components::gen_UUID()

```
mixed server::\Auth\components::gen_UUID()()
```





* Visibility: **public**
* This method is **static**.



### \Auth\components::keygen()

```
mixed server::\Auth\components::keygen()($amount)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $amount **mixed**



### \Auth\components::getCookie()

```
mixed server::\Auth\components::getCookie()($cookie)
```





* Visibility: **public**

#### Arguments

* $cookie **mixed**



### \user::__construct()

```
void user::\user::__construct()(string $ip, string|null $session, string|null $hash)
```

Constructor Class responsible for filling in arrays with the data from a specified user.



* Visibility: **public**
* This method is defined by [user](user.md)

#### Arguments

* $ip **string** - &lt;p&gt;The IP address of a user who is requesting the function, or if called from the Admin CP it is the user id.&lt;/p&gt;
* $session **string|null** - &lt;p&gt;The value of the pp_auth_token cookie.&lt;/p&gt;
* $hash **string|null** - &lt;p&gt;The server hash of the requesting user which is used when they are viewing node pages.&lt;/p&gt;



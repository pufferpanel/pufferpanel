user
===============

PufferPanel Core User Class File




* Class name: user
* Namespace: 
* Parent class: [Auth\auth](Auth-auth.md)





Properties
----------


### $_data

```
private mixed $_data
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


### \user::__construct()

```
void user::\user::__construct()(string $ip, string|null $session, string|null $hash)
```

Constructor Class responsible for filling in arrays with the data from a specified user.



* Visibility: **public**

#### Arguments

* $ip **string** - &lt;p&gt;The IP address of a user who is requesting the function, or if called from the Admin CP it is the user id.&lt;/p&gt;
* $session **string|null** - &lt;p&gt;The value of the pp_auth_token cookie.&lt;/p&gt;
* $hash **string|null** - &lt;p&gt;The server hash of the requesting user which is used when they are viewing node pages.&lt;/p&gt;



### \user::rebuildData()

```
void user::\user::rebuildData()(string $id)
```

Re-runs the __construct() class with a defined ID for the admin control panel.



* Visibility: **public**

#### Arguments

* $id **string** - &lt;p&gt;The ID of a user requested in the Admin CP.&lt;/p&gt;



### \user::getData()

```
string|array|boolean user::\user::getData()(string|null $id)
```

Provides the corresponding value for the id provided from the MySQL Database.



* Visibility: **public**

#### Arguments

* $id **string|null** - &lt;p&gt;The column value for the data you need (e.g. email).&lt;/p&gt;



### \Auth\auth::__construct()

```
mixed user::\Auth\auth::__construct()()
```





* Visibility: **public**



### \Auth\auth::validateTOTP()

```
mixed user::\Auth\auth::validateTOTP()($key, $secret)
```





* Visibility: **public**

#### Arguments

* $key **mixed**
* $secret **mixed**



### \Auth\auth::verifyPassword()

```
mixed user::\Auth\auth::verifyPassword()($email, $raw)
```





* Visibility: **public**

#### Arguments

* $email **mixed**
* $raw **mixed**



### \Auth\auth::isLoggedIn()

```
mixed user::\Auth\auth::isLoggedIn()($ip, $session, $serverhash, $acp)
```





* Visibility: **public**

#### Arguments

* $ip **mixed**
* $session **mixed**
* $serverhash **mixed**
* $acp **mixed**



### \Database\database::buildConnection()

```
mixed user::\Database\database::buildConnection()()
```





* Visibility: **public**
* This method is **static**.



### \Database\database::connect()

```
mixed user::\Database\database::connect()()
```





* Visibility: **public**
* This method is **static**.



### \Auth\components::hash()

```
mixed user::\Auth\components::hash()($raw)
```





* Visibility: **public**

#### Arguments

* $raw **mixed**



### \Auth\components::password_compare()

```
mixed user::\Auth\components::password_compare()($raw, $hashed)
```





* Visibility: **private**

#### Arguments

* $raw **mixed**
* $hashed **mixed**



### \Auth\components::generate_iv()

```
mixed user::\Auth\components::generate_iv()()
```





* Visibility: **public**



### \Auth\components::encrypt()

```
mixed user::\Auth\components::encrypt()($raw, $iv, $method)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $raw **mixed**
* $iv **mixed**
* $method **mixed**



### \Auth\components::decrypt()

```
mixed user::\Auth\components::decrypt()($encrypted, $iv, $method)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $encrypted **mixed**
* $iv **mixed**
* $method **mixed**



### \Auth\components::gen_UUID()

```
mixed user::\Auth\components::gen_UUID()()
```





* Visibility: **public**
* This method is **static**.



### \Auth\components::keygen()

```
mixed user::\Auth\components::keygen()($amount)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $amount **mixed**



### \Auth\components::getCookie()

```
mixed user::\Auth\components::getCookie()($cookie)
```





* Visibility: **public**

#### Arguments

* $cookie **mixed**



### \Page\components::redirect()

```
mixed user::\Page\components::redirect()($url)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $url **mixed**



### \Page\components::genRedirect()

```
mixed user::\Page\components::genRedirect()()
```





* Visibility: **public**
* This method is **static**.



### \Page\components::twigGET()

```
mixed user::\Page\components::twigGET()()
```





* Visibility: **public**
* This method is **static**.



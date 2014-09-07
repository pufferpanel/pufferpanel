Auth\auth
===============






* Class name: auth
* Namespace: Auth





Properties
----------


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


### \Auth\auth::__construct()

```
mixed Auth\auth::\Auth\auth::__construct()()
```





* Visibility: **public**



### \Auth\auth::validateTOTP()

```
mixed Auth\auth::\Auth\auth::validateTOTP()($key, $secret)
```





* Visibility: **public**

#### Arguments

* $key **mixed**
* $secret **mixed**



### \Auth\auth::verifyPassword()

```
mixed Auth\auth::\Auth\auth::verifyPassword()($email, $raw)
```





* Visibility: **public**

#### Arguments

* $email **mixed**
* $raw **mixed**



### \Auth\auth::isLoggedIn()

```
mixed Auth\auth::\Auth\auth::isLoggedIn()($ip, $session, $serverhash, $acp)
```





* Visibility: **public**

#### Arguments

* $ip **mixed**
* $session **mixed**
* $serverhash **mixed**
* $acp **mixed**



### \Database\database::buildConnection()

```
mixed Auth\auth::\Database\database::buildConnection()()
```





* Visibility: **public**
* This method is **static**.



### \Database\database::connect()

```
mixed Auth\auth::\Database\database::connect()()
```





* Visibility: **public**
* This method is **static**.



### \Auth\components::hash()

```
mixed Auth\auth::\Auth\components::hash()($raw)
```





* Visibility: **public**

#### Arguments

* $raw **mixed**



### \Auth\components::password_compare()

```
mixed Auth\auth::\Auth\components::password_compare()($raw, $hashed)
```





* Visibility: **private**

#### Arguments

* $raw **mixed**
* $hashed **mixed**



### \Auth\components::generate_iv()

```
mixed Auth\auth::\Auth\components::generate_iv()()
```





* Visibility: **public**



### \Auth\components::encrypt()

```
mixed Auth\auth::\Auth\components::encrypt()($raw, $iv, $method)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $raw **mixed**
* $iv **mixed**
* $method **mixed**



### \Auth\components::decrypt()

```
mixed Auth\auth::\Auth\components::decrypt()($encrypted, $iv, $method)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $encrypted **mixed**
* $iv **mixed**
* $method **mixed**



### \Auth\components::gen_UUID()

```
mixed Auth\auth::\Auth\components::gen_UUID()()
```





* Visibility: **public**
* This method is **static**.



### \Auth\components::keygen()

```
mixed Auth\auth::\Auth\components::keygen()($amount)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $amount **mixed**



### \Auth\components::getCookie()

```
mixed Auth\auth::\Auth\components::getCookie()($cookie)
```





* Visibility: **public**

#### Arguments

* $cookie **mixed**



### \Page\components::redirect()

```
mixed Auth\auth::\Page\components::redirect()($url)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $url **mixed**



### \Page\components::genRedirect()

```
mixed Auth\auth::\Page\components::genRedirect()()
```





* Visibility: **public**
* This method is **static**.



### \Page\components::twigGET()

```
mixed Auth\auth::\Page\components::twigGET()()
```





* Visibility: **public**
* This method is **static**.



tplMail
===============






* Class name: tplMail
* Namespace: 





Properties
----------


### $message

```
private mixed $message
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


### \tplMail::__construct()

```
mixed tplMail::\tplMail::__construct()()
```





* Visibility: **public**



### \tplMail::dispatch()

```
mixed tplMail::\tplMail::dispatch()($email, $subject)
```





* Visibility: **public**

#### Arguments

* $email **mixed**
* $subject **mixed**



### \tplMail::getDispatchSystemFunct()

```
mixed tplMail::\tplMail::getDispatchSystemFunct()()
```





* Visibility: **private**



### \tplMail::readTemplate()

```
mixed tplMail::\tplMail::readTemplate()($template)
```





* Visibility: **private**

#### Arguments

* $template **mixed**



### \tplMail::generateLoginNotification()

```
mixed tplMail::\tplMail::generateLoginNotification()($type, $vars)
```





* Visibility: **public**

#### Arguments

* $type **mixed**
* $vars **mixed**



### \tplMail::buildEmail()

```
mixed tplMail::\tplMail::buildEmail()($tpl, $data)
```





* Visibility: **public**

#### Arguments

* $tpl **mixed**
* $data **mixed**



### \Auth\components::hash()

```
mixed tplMail::\Auth\components::hash()($raw)
```





* Visibility: **public**

#### Arguments

* $raw **mixed**



### \Auth\components::password_compare()

```
mixed tplMail::\Auth\components::password_compare()($raw, $hashed)
```





* Visibility: **private**

#### Arguments

* $raw **mixed**
* $hashed **mixed**



### \Auth\components::generate_iv()

```
mixed tplMail::\Auth\components::generate_iv()()
```





* Visibility: **public**



### \Auth\components::encrypt()

```
mixed tplMail::\Auth\components::encrypt()($raw, $iv, $method)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $raw **mixed**
* $iv **mixed**
* $method **mixed**



### \Auth\components::decrypt()

```
mixed tplMail::\Auth\components::decrypt()($encrypted, $iv, $method)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $encrypted **mixed**
* $iv **mixed**
* $method **mixed**



### \Auth\components::gen_UUID()

```
mixed tplMail::\Auth\components::gen_UUID()()
```





* Visibility: **public**
* This method is **static**.



### \Auth\components::keygen()

```
mixed tplMail::\Auth\components::keygen()($amount)
```





* Visibility: **public**
* This method is **static**.

#### Arguments

* $amount **mixed**



### \Auth\components::getCookie()

```
mixed tplMail::\Auth\components::getCookie()($cookie)
```





* Visibility: **public**

#### Arguments

* $cookie **mixed**



### \Database\database::buildConnection()

```
mixed tplMail::\Database\database::buildConnection()()
```





* Visibility: **public**
* This method is **static**.



### \Database\database::connect()

```
mixed tplMail::\Database\database::connect()()
```





* Visibility: **public**
* This method is **static**.



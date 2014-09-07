settings
===============

PufferPanel Core Settings Class File




* Class name: settings
* Namespace: 





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


### \settings::__construct()

```
void settings::\settings::__construct()()
```

Constructor class for building settings data.



* Visibility: **public**



### \settings::get()

```
array|string settings::\settings::get()(string $setting)
```

Function to retrieve various panel settings.



* Visibility: **public**

#### Arguments

* $setting **string** - &lt;p&gt;The name of the setting for which you want the value.&lt;/p&gt;



### \settings::nodeName()

```
string settings::\settings::nodeName()(integer $id)
```

Convert a node ID into a name for the node.



* Visibility: **public**

#### Arguments

* $id **integer** - &lt;p&gt;The ID of the node you want the name for.&lt;/p&gt;



### \Database\database::buildConnection()

```
mixed settings::\Database\database::buildConnection()()
```





* Visibility: **public**
* This method is **static**.



### \Database\database::connect()

```
mixed settings::\Database\database::connect()()
```





* Visibility: **public**
* This method is **static**.



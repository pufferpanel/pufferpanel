# Contributing to PufferPanel

The following standards should be followed when contributing to PufferPanel. This allows the code to be fluid across all contributions.

## PSR-0 and PSR-4 Standards
PufferPanel follows [PSR-0](http://www.php-fig.org/psr/psr-0/), [PSR-1](http://www.php-fig.org/psr/psr-1/), and [PSR-4](http://www.php-fig.org/psr/psr-4/) standards for autoloading. For most basic contributions you shouldn't have any issues following these standards.

Following the **PSR-1** formatting is absolutely required for your pull requests to be accepted.

## Files
* There should not be a newline at the end of a file.
* Files that contain only PHP code should not end with a closing tag ``?>``

## Lines
* There **must not** be a hard limit on line length.
* There **must not** be trailing whitespace at the end of non-blank lines.
* Blank lines **should** be added to improve readability and to indicate related blocks of code. These lines **should not** be indented.

## Indenting & Spacing
* Indents should be tabs, **no spaces**.
* Statements should not have a space before the opening parenthesis (``if(!$this) {`` is good, ``if (!$this) {`` is bad).
* Curly brackets should **always** be on the same line as the statement they are being used for.
* All functions should have a space after the closing parenthesis and before the curly brackets ``{ }``. Additionally, any function using curly brackets should have a space between the function and the bracket.
* For single line ``if/else`` statements, curly brackets **are required**. This change was implemented after a lot of previous code was written, so do not simply copy what is already there.

### Good Code
```php
<?php

class My_Class {

	public function myFunction() {
	
	}

}

if(isset($this)) {

	echo $something;

} else {

	echo $foobar;

}
```

### Bad Code
```php

class My_Class{

	public function myFunction()
		{
		
		}

}

if ($this) {

}

if($this)
{

}

if(!$something)
	echo 'Foo';
else
	echo 'Bar';
```

## Naming Classes and Functions
* Classes should be named using ``UpperCamelCase`` and functions should be named using ``lowerCamelCase``.
* ``static`` function naming should be done after declaring the visibility of the function (``public static`` vs ``static public``).
* ``abstract`` and ``final`` function calls should be made before declaring visibility.
* Functions must be named according to what they do, and should be properly documented using the ``phpDocumentor`` syntax. Please see other functions in the panel for how to do this properly.
* Protected functions are prefered over private functions when needed. Please name protected functions with a preceding underscore (``_protectedFunction``), and private functions with two preceeding underscores (``__privateFunction``).

```php
<?php
namespace PufferPanel\Core;

class My_Foobar_Class {

	public function __construct() {

	}
	
	final public static function finalPubStatic() {

	}
	
	protected static function _protectedFunction() {

	}

}
```

## Switch Statements
* Code should be indented for each new level of the case, see the example below for acceptable code.

```php
switch($foo) {

	case 1:
		echo 'Case One';
		break;
	case 2:
		echo 'Case Two';
		break;
	default:
		echo 'Default Case';
		break;

}
```

## Method and Function Calls
* There **must not** be a space between the method or function name and the opening parenthesis
* There **must not** be a space after the opening parenthesis
* There **must not** be a space before the closing parenthesis.
* For argument lists, there **must not** be a space before each comma, and there **must** be one space after each comma.

```php
<?php

$class->myFunction($foo);
Class::myStaticFunction($bar, $zar);
```

## Argument Lists
* Argument lists should be split across multiple lines for readability when they get long. When doing so, the first item in the list must be on the next line, and there must be only one argument per line.

```php
<?php

$class->myLongFunction(
    $foo,
    $bar->barista->makeCoffee(),
    array($zar, $tar, $mar)
);
```

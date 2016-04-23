# Contributing to PufferPanel

We welcome any contributions to PufferPanel, whether in the form of bug reports, feature suggestions, and even code submissions. When submitting any of those, please follow this documentation to help speed up the process.

# Reporting Bugs

A bug is when the software does not behave in a way that is expected. It is **not** invalid configurations which render the panel broken.

If you believe you have located a bug, please report it to the [Bug Tracker](https://github.com/PufferPanel/PufferPanel/issues).

**Please make sure there is not an issue for your specific bug already!** If you find that someone else has reported a bug you have, please comment on that issue stating you have replicated that bug. Do not make a new issue.

When submitting those bugs, follow these standards:
* The title of the issue should **clearly** and **quickly** explain the issue. A good title would be "Cannot delete IPs from node if it has 2 or more ports".
* The description should contain the following information
  * A complete description of the problem. This should explain what you expect the panel to do and what the panel actually did.
  * Steps to reproduce the bug. It is hard to figure out what the bug truly is if we cannot do it ourselves.

# Submitting feature requests

If you have an idea for a new feature or enhancement, please suggest it on our [Community Forum](https://community.pufferpanel.com/forum/5-feature-requests/)

# How to Contribute

When submitting new code to the panel, you **must** follow both the standards outlined later in this documentation, along with the following:
* All PRs must contain a reference to an **existing** issue. If there is no issue for your PR to reference, then create a new issue, following the guidelines above.
* PRs may only contain **1** feature or enhancement. Kitchen sinks will be throw out the window.

# Standards

The following standards should be followed when contributing to PufferPanel.

## PSR-0 and PSR-4 Standards
PufferPanel follows [PSR-0](http://www.php-fig.org/psr/psr-0/), [PSR-1](http://www.php-fig.org/psr/psr-1/), and [PSR-4](http://www.php-fig.org/psr/psr-4/) standards for autoloading. For most basic contributions you shouldn't have any issues following these standards.

Following the **PSR-1** formatting is absolutely required for your pull requests to be accepted.

## Files
* There **should** be a newline at the end of a file.
* Files that contain only PHP code **should not** end with a closing tag ``?>``

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

## Requires and Includes
* Parenthesis should not be used for including or requiring outside files.
```php
<?php
include 'something.php';
require_once 'some/other/file.php';
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

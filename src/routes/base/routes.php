<?php
/*
PufferPanel - A Minecraft Server Management Panel
Copyright (c) 2013 Dane Everitt

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see http://www.gnu.org/licenses/.
*/

$klein->respond('GET', '/account', function() use ($core) {
});

$klein->respond('POST', '/account', function() use ($core) {
});

$klein->respond('GET', '/logout', function() use ($core) {
	return 'Called /logout';
});

$klein->respond('POST', '/logout', function() use ($core) {
});

$klein->respond('GET', '/password', function() use ($core) {
});

$klein->respond('POST', '/password', function() use ($core) {
});

$klein->respond('GET', '/register', function() use ($core) {
});

$klein->respond('POST', '/register', function() use ($core) {
});

$klein->respond('GET', '/servers', function() use ($core) {
});

$klein->respond('POST', '/servers', function() use ($core) {
});
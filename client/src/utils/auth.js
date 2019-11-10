/*
 * Copyright 2019 Padduck, LLC
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *          http://www.apache.org/licenses/LICENSE-2.0
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

import Cookies from 'js-cookie'

export function hasAuth () {
  return !((Cookies.get('puffer_auth') || '') === '')
}

export function isAdmin () {
  return hasScope('servers.admin')
}

export function hasScope (scope) {
  const scopes = localStorage.getItem('scopes') || '[]'
  return JSON.parse(scopes).indexOf(scope) !== -1
}

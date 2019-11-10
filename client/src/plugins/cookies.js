/*
 * Copyright 2019 Padduck, LLC
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *  	http://www.apache.org/licenses/LICENSE-2.0
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

import Vue from 'vue'

Vue.prototype.hasScope = function (scope, server = '') {
  let data = localStorage.scopes
  if (!data) {
    return false
  }

  let scopeMap = JSON.parse(data)
  if (!scopeMap) {
    return false
  }

  let scopes = scopeMap[server]
  if (!scopes) {
    return false
  }

  for (let i = 0; i < scopes.length; i++) {
    if (scopes[i] === scope) {
      return true
    }
  }

  return false
}
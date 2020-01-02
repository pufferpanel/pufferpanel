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

export default {
  validPassword: function (password) {
    return password.length >= 8
  },
  samePassword: function (pass1, pass2) {
    return pass1 && pass2 && pass1 === pass2
  },
  validUsername: function (username) {
    return username && /^([0-9A-Za-z_-]){3,}$/.test(username)
  },
  validEmail: function (email) {
    return email && /^\S+@\S+\.\S+$/.test(email)
  }
}

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

export const handleError = (ctx, overrides = {}) => error => {
  // eslint-disable-next-line no-console
  console.log('ERROR', error)
  let msg = 'errors.ErrUnknownError'
  if (error && error.response && error.response.data.error) {
    if (error.response.data.error.code) {
      msg = 'errors.' + error.response.data.error.code
    } else {
      msg = error.response.data.error.msg
    }
  }

  if (overrides[error.response.status] !== undefined) msg = overrides[error.response.status]

  ctx.$toast.error(ctx.$t(msg))
}

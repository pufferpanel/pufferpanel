<!--
  - Copyright 2019 Padduck, LLC
  -  Licensed under the Apache License, Version 2.0 (the "License");
  -  you may not use this file except in compliance with the License.
  -  You may obtain a copy of the License at
  -  	http://www.apache.org/licenses/LICENSE-2.0
  -  Unless required by applicable law or agreed to in writing, software
  -  distributed under the License is distributed on an "AS IS" BASIS,
  -  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  -  See the License for the specific language governing permissions and
  -  limitations under the License.
  -->

<template>
  <b-container>
    <b-card header-tag="header"
            footer-tag="footer">
      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.Username')"
                    label-for="input-sm">
        <b-form-input v-model="username" id="input-sm" size="sm"></b-form-input>
      </b-form-group>

      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.Email')" label-for="input-sm">
        <b-form-input v-model="email" id="input-sm" size="sm" type="password"></b-form-input>
      </b-form-group>

      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.ConfirmPassword')"
                    label-for="input-sm">
        <b-form-input id="input-sm" size="sm" type="password"></b-form-input>
      </b-form-group>

      <b-button slot="footer" variant="primary" size="sm" v-text="$t('common.Update')"></b-button>
    </b-card>


    <!-- -->

    <b-card header-tag="header"
            footer-tag="footer">
      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.OldPassword')"
                    label-for="input-sm">
        <b-form-input v-model="oldPassword" id="input-sm" size="sm" type="password"></b-form-input>
      </b-form-group>

      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.NewPassword')"
                    label-for="input-sm">
        <b-form-input v-model="newPassword" id="input-sm" size="sm" type="password"></b-form-input>
      </b-form-group>

      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.ConfirmPassword')"
                    label-for="input-sm">
        <b-form-input v-model="newPassword2" id="input-sm" size="sm" type="password"></b-form-input>
      </b-form-group>

      <b-button slot="footer" variant="primary" size="sm" v-text="$t('common.ChangePassword')"></b-button>
    </b-card>
  </b-container>
</template>

<script>
export default {
  data () {
    return {
      username: '',
      email: '',
      confirmPassword: '',
      oldPassword: '',
      newPassword: '',
      newPassword2: ''
    }
  },
  mounted () {

  }
}
</script>

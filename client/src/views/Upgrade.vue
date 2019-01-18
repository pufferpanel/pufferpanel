<template>
  <v-container>
    <v-layout justify-center>
      <v-flex
        xs12
        sm10
        md8
      >
        <material-card
          color="primary"
          title="Vue Material Dashboard PRO"
          text="Are you looking for more components? Please check our Premium Version of Vue Material Dashboard."
        >
          <table class="v-table">
            <thead>
              <tr>
                <th/>
                <th class="subheading font-weight-light">Free</th>
                <th class="subheading font-weight-light">PRO</th>
              </tr>
            </thead>
            <tbody class="text-xs-center">
              <tr>
                <th class="text-xs-left font-weight-light subheading">Components</th>
                <td>60</td>
                <td>200</td>
              </tr>
              <tr>
                <th class="text-xs-left font-weight-light subheading">Plugins</th>
                <td>2</td>
                <td>15</td>
              </tr>
              <tr>
                <th class="text-xs-left font-weight-light subheading">Example Pages</th>
                <td>3</td>
                <td>27</td>
              </tr>
              <tr>
                <th class="text-xs-left font-weight-light subheading">Login, Register, Pricing, Lock Pages</th>
                <td>
                  <v-icon color="error">mdi-close</v-icon>
                </td>
                <td>
                  <v-icon color="success">mdi-check</v-icon>
                </td>
              </tr>
              <tr>
                <th class="text-xs-left font-weight-light subheading">Premium Support</th>
                <td>
                  <v-icon color="error">mdi-close</v-icon>
                </td>
                <td>
                  <v-icon color="success">mdi-check</v-icon>
                </td>
              </tr>
              <tr>
                <th/>
                <td>Free</td>
                <td>Just $59</td>
              </tr>
              <tr>
                <th/>
                <td>
                  <v-btn
                    class="subheading white--text font-weight-light"
                    round
                    large
                    color="grey"
                    disabled
                  >Current Version</v-btn>
                </td>
                <td>
                  <v-btn
                    class="subheading white--text font-weight-light"
                    round
                    large
                    color="cyan"
                  >Upgrade to Pro</v-btn>
                </td>
              </tr>
            </tbody>
          </table>
        </material-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

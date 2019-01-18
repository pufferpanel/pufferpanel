<template>
  <v-container
    fill-height
    fluid
  >
    <v-layout
      justify-center
      align-center
    >
      <v-flex xs12>
        <material-card
          color="green"
          title="Material Dashboard Heading"
          text="Created using Roboto Font Family"
        >
          <v-card-text>
            <h2 class="font-weight-light mb-4">Typography</h2>

            <v-container
              class="pa-0"
              grid-list-xl
              fluid
            >
              <v-layout
                v-for="(t, i) in typography"
                :key="i"
                align-end
                wrap
              >
                <v-flex
                  xs1
                  md3>
                  <span
                    class="tim-note"
                    v-text="t[0]"
                  />
                </v-flex>
                <v-flex xs8>
                  <component
                    :is="t[2]"
                    :class="i"
                  >
                    <template v-if="i !== 'quote'">
                      {{ t[1] }}
                    </template>

                    <p v-if="i === 'quote'">{{ t[1] }}</p>
                    <small v-if="i === 'quote'">Kanye West, Musician</small>

                    <template v-if="i === 'small'">
                      <br>
                      <small>Use 'small' tag for the headers</small>
                    </template>
                  </component>
                </v-flex>
              </v-layout>
            </v-container>
          </v-card-text>
        </material-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
const leader = 'I will be the leader of a company that ends up being worth billions of dollars, because I got the answers. I understand culture. I am the nucleus. I think that’s a responsibility that I have, to push possibilities, to show people, this is the level that things could be at.'
const leaderShort = leader.slice(0, 105) + '...'
const material = 'The Life of Material Dashboard'
const small = 'Header with small subtitle'

export default {
  data: () => ({
    typography: {
      'heading-1': ['Header 1', material, 'h1'],
      'heading-2': ['Header 2', material, 'h2'],
      'heading-3': ['Header 3', material, 'h3'],
      'heading-4': ['Header 4', material, 'h4'],
      'heading-5': ['Header 5', material, 'h5'],
      'heading-6 text-uppercase': ['Header 6', material, 'h6'],
      '': ['Paragraph', leader, 'p'],
      'quote': ['Quote', leader, 'blockquote'],
      'text--disabled': ['Muted Text', leaderShort, 'p'],
      'text-primary': ['Primary Text', leaderShort, 'p'],
      'text-info': ['Info Text', leaderShort, 'p'],
      'text-success': ['Success Text', leaderShort, 'p'],
      'text-warning': ['Warning Text', leaderShort, 'p'],
      'text-danger': ['Danger Text', leaderShort, 'p'],
      'small': ['Small Tag', small, 'h2']
    }
  })
}
</script>

<style lang="scss">
  .tim-note {
    bottom: 10px;
    color: #c0c1c2;
    display: block;
    font-weight: 400;
    font-size: 13px;
    line-height: 13px;
    left: 0;
    margin-left: 20px;
    width: 260px;
  }
</style>

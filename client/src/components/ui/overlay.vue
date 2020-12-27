<template>
  <v-overlay
    v-model="value"
    :dark="isDark()"
    color="#434343"
  >
    <v-container
      fluid
      class="overlayContainer"
    >
      <v-row>
        <v-col
          cols="12"
          offset-md="1"
          md="10"
        >
          <v-card v-if="card">
            <v-card-title>
              <span v-text="title" />
              <v-spacer />
              <v-btn
                v-if="closable"
                icon
                @click="$emit('close'); $emit('input', false)"
              >
                <v-icon>mdi-close</v-icon>
              </v-btn>
            </v-card-title>
            <v-card-text class="overlayContent">
              <slot />
            </v-card-text>
            <v-card-actions>
              <slot name="actions" />
            </v-card-actions>
          </v-card>
          <v-sheet
            v-else
            class="overlayContent"
          >
            <slot />
          </v-sheet>
        </v-col>
      </v-row>
    </v-container>
  </v-overlay>
</template>

<style>
.overlayContainer {
  width: 100vw !important;
  max-width: 100vw !important;
  max-height: 100vh !important;
}

.overlayContent {
  max-height: calc(90vh - 68px);
  overflow-y: scroll;
}
</style>

<script>
import { isDark } from '@/utils/dark'

export default {
  props: {
    card: { type: Boolean, default: () => false },
    closable: { type: Boolean, default: () => false },
    title: { type: String, default: () => '' },
    value: { type: Boolean, default: () => false }
  },
  methods: {
    isDark
  }
}
</script>

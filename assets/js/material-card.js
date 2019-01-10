Vue.component('material-card', {
    template: '  <v-card v-bind="$attrs" :style="styles" v-on="$listeners">\n' +
        '    <helper-offset v-if="hasOffset" :inline="inline" :full-width="fullWidth" :offset="offset">\n' +
        '      <v-card v-if="!$slots.offset" :color="color" :class="`elevation-${elevation}`" class="v-card--material__header" dark>\n' +
        '        <slot v-if="!title && !text" name="header" />\n' +
        '        <span v-else>\n' +
        '          <v-img v-if="logo" :src="logo" height="34" contain />' +
        '          <h4 class="title font-weight-light mb-2" v-text="title" />\n' +
        '          <p class="category font-weight-thin" v-text="text" />\n' +
        '        </span>\n' +
        '      </v-card>\n' +
        '      <slot v-else name="offset" />\n' +
        '    </helper-offset>\n' +
        '\n' +
        '    <v-card-text>\n' +
        '      <slot />\n' +
        '    </v-card-text>\n' +
        '\n' +
        '    <v-divider v-if="$slots.actions" class="mx-3" />\n' +
        '\n' +
        '    <v-card-actions v-if="$slots.actions">\n' +
        '      <slot name="actions" />\n' +
        '    </v-card-actions>\n' +
        '  </v-card>',
    props: {
        color: {
            type: String,
            default: 'secondary'
        },
        elevation: {
            type: [Number, String],
            default: 10
        },
        inline: {
            type: Boolean,
            default: false
        },
        fullWidth: {
            type: Boolean,
            default: false
        },
        offset: {
            type: [Number, String],
            default: 24
        },
        title: {
            type: String,
            default: undefined
        },
        text: {
            type: String,
            default: undefined
        },
        logo: {
            type: String,
            default: undefined
        }
    },

    computed: {
        hasOffset () {
            return this.$slots.header ||
                this.$slots.offset ||
                this.title ||
                this.text
        },
        styles () {
            if (!this.hasOffset) return null

            return {
                marginBottom: `${this.offset}px`,
                marginTop: `${this.offset * 2}px`
            }
        }
    }
});
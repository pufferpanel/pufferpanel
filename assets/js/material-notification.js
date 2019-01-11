Vue.component('material-notification', {
    template: '<v-alert v-bind="$attrs" :class="[`elevation-${elevation}`]" :value="value" class="v-alert--notification" v-on="$listeners"> <slot /> </v-alert>',
    inheritAttrs: false,
    props: {
        elevation: {
            type: [Number, String],
            default: 6
        },
        value: {
            type: Boolean,
            default: true
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
Vue.component('helper-offset', {
    template: '<div :style="styles" :class="classes" class="v-offset"><slot /></div>',
    props: {
        fullWidth: {
            type: Boolean,
            default: false
        },
        offset: {
            type: [Number, String],
            default: 0
        }
    },
    computed: {
        classes() {
            return {
                'v-offset--full-width': this.fullWidth
            }
        },
        styles() {
            return {
                top: `-${this.offset}px`,
                marginBottom: `-${this.offset}px`
            }
        }
    }
});
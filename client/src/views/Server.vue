<template>
  <b-row v-if="loading">
    <b-col cols="5"/>
    <b-col cols="2">
      <b-spinner class="align-middle"/>
      <strong :text="$t('common.Loading')" v-text="$t('common.Loading')"></strong>
    </b-col>
  </b-row>
  <b-container v-if="!loading">

  </b-container>
</template>

<script>
export default {
  data () {
    return {
      loading: true,
      type: 'generic'
    }
  }
}
</script>
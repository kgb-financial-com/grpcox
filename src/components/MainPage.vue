<template>
  <div class="mainPage">
    <h1>Backend Services gRPC Frontend</h1>
    <div>Base URL: {{urlBase.href}}</div>
  </div>
</template>

<script>
export default {
  name: 'MainPage',
  props: {
    msg: String
  },
  computed: {
    urlBase() {
      let base = process.env.VUE_APP_URL_BASE;
      return base ? new URL(base) : new URL(".", document.baseURI);
    }
  }
}
</script>

<style>
</style>

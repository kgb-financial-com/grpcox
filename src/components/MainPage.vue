<template>
  <div class="mainPage">
    <h1>Backend Services gRPC Frontend</h1>
    <HostSelector/>
  </div>
</template>

<script>
import HostSelector from "@/components/HostSelector";
export default {
  name: 'MainPage',
  components: { HostSelector }
}
</script>

<style>
</style>

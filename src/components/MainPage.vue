<template>
  <div class="mainPage">
    <h1>Backend Services gRPC Frontend</h1>
    <HostAndServiceSelector/>
  </div>
</template>

<script>
import HostAndServiceSelector from "@/components/HostAndServiceSelector";
export default {
  name: 'MainPage',
  components: { HostAndServiceSelector }
}
</script>

<style>
</style>

<template>
  <div class="mainPage">
    <h1>Backend Services gRPC Frontend</h1>
    <HostAndServiceSelector/>
    <MethodContainer/>
  </div>
</template>

<script>
import HostAndServiceSelector from "@/components/HostAndServiceSelector";
import MethodContainer from "@/components/MethodContainer";

export default {
  name: 'MainPage',
  components: { HostAndServiceSelector, MethodContainer }
}
</script>

<style>
</style>

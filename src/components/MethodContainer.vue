<template>
  <div class="methodContainer">
    <Method v-for="method in methods" :key="method" :name="method"/>
  </div>
</template>

<script>
import axios from "axios";
import Method from "@/components/Method";

export default {
  name: "MethodContainer",

  data: function() {
    return {
      methods: []
    }
  },

  components: { Method },

  computed: {
    selectedHost() {
      return this.$store.state.selectedHost;
    },
    selectedService() {
      return this.$store.state.selectedService;
    }
  },

  watch:{
    selectedService() {
      this.refreshMethodList();
    }
  },

  created() {
    this.refreshMethodList();
  },

  methods: {

    refreshMethodList() {

      if (!this.selectedService) {
        return;
      }

      axios.get(this.$store.state.urlBase + "server/" + this.selectedHost + "/service/" + this.selectedService + "/functions")
          .then(data => {
            if (!data?.data?.data) {
              throw "Could not connect to host: " + this.selectedHost.split('|')[1];
            }
            this.methods = data.data.data;
          })
          .catch(err => this.errorMessage = "Could not connect to server " + this.selectedHost.split('|')[0] + "(" + this.selectedHost.split('|')[1] + "): " + err)

    }

  }
}
</script>

<style scoped>

</style>
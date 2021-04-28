<template>
  <div>Methods... {{ methods }}</div>
</template>

<script>
import axios from "axios";

export default {
  name: "MethodContainer",

  data: function() {
    return {
      methods: []
    }
  },

  computed: {
    selectedName() {
      return this.$store.state.selectedHost.name;
    },
    selectedHost() {
      return this.$store.state.selectedHost.host;
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
              throw "Could not connect to host: " + this.selectedHost;
            }
            this.methods = data.data.data;
          })
          .catch(err => this.errorMessage = "Could not connect to server " + this.selectedName + "(" + this.selectedHost + "): " + err)

    }

  }
}
</script>

<style scoped>

</style>
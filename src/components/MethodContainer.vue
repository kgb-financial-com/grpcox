<template>
  <div class="methodContainer">
    <Method v-for="method in methods" :key="method" :name="method"/>
    <div v-if="!!errorMessage" class="error">{{ errorMessage }}</div>
  </div>
</template>

<script>
import axios from "axios";
import Method from "@/components/Method";

export default {
  name: "MethodContainer",

  data: function() {
    return {
      methods: [],
      errorMessage: null
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
    },
    selectedHost() {
      this.refreshMethodList();
    }
  },

  created() {
    this.refreshMethodList();
  },

  methods: {

    refreshMethodList() {

      if (!this.selectedService || !this.selectedHost) {
        return;
      }
      this.errorMessage = null;

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

.error {
  border: solid 1px;
  padding: 10px;
  background-color: #ffd0d0;
}

</style>

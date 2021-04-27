<template>
  <div class="error" v-if="!!errorMessage">{{errorMessage}}</div>
  <div class="hostList" v-else>
    <select v-model="selectedFromList" multiple>
      <option v-for="option in availableHosts" :key="option">{{option}}</option>
    </select>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: "HostSelector",
  data() {
    return {
      availableHosts: [],
      errorMessage: null,
      selectedFromList: []
    }
  },
  computed: {
    selectedName() {
      return this.$store.state.selectedHost.name;
    },
    selectedHost() {
      return this.$store.state.selectedHost.host;
    }
  },
  created() {
    this.refreshHostList();
  },
  methods: {
    refreshHostList() {
      axios.get(this.$store.state.urlBase + "server/" + encodeURIComponent(this.selectedHost) + "/services")
          .then(data =>
              this.availableHosts = data.data.data
                  .map( s => s.substring(0, s.length - 1))
                  .filter(s => !s.startsWith("grpc."))
          )
          .catch(err => this.errorMessage = "Could not connect to server " + this.selectedName
              + "(" + this.selectedHost + "): " + err.errorMessage);
    }
  }
}
</script>

<style>

</style>
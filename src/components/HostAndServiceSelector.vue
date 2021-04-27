<template>

  <div class="error" v-if="!!errorMessage">{{ errorMessage }}</div>

  <div class="hostList" v-else>

    <div v-if="!!errorHostsMessage">{{ errorHostsMessage }}</div>
    <div v-else>
      <label for="hostDropDown">Host</label>
      <select id="hostDropDown" v-model="selectedHostInList">
        <option v-for="option in availableHosts" :key="option">{{ option }}</option>
      </select>
    </div>

    <div>
      <label for="serviceDropDown">Service</label>
      <select id="serviceDropDown" v-model="selectedServiceInList">
        <option v-for="option in currentServices" :key="option">{{ option }}</option>
      </select>
    </div>

  </div>

</template>

<script>
import axios from 'axios';

export default {

  name: "HostAndServiceSelector",

  data() {
    return {
      currentServices: [],
      availableHosts: [],
      errorMessage: null,
      errorHostsMessage: null,
      selectedHostInList: null,
      selectedServiceInList: null
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

    selectFirst(array) {
      return array && array.length ? array[0] : null;
    },

    refreshHostList() {
      axios.get(this.$store.state.urlBase + "server/" + encodeURIComponent(this.selectedHost) + "/services")
          .then(data => {
            if (!data?.data?.data) {
              throw "Could not connect to host: " + this.selectedHost;
            }
            this.currentServices = data.data.data
                .map(s => s.substring(0, s.length - 1))
                .filter(s => !s.startsWith("grpc."))
            this.selectedServiceInList = this.selectFirst(this.currentServices);
          })
          .catch(err => this.errorMessage = "Could not connect to server " + this.selectedName
              + "(" + this.selectedHost + "): " + err)
          .finally(() => axios.get(this.$store.state.urlBase + "active/get")
              .then(data => {
                this.availableHosts = data.data.data;
                this.selectedHostInList = this.selectFirst(this.availableHosts);
              })
              .catch(err => this.errorHostsMessage = "Could not retrieve list of available hosts: " + err.errorHostsMessage)
          )
    }

  }
}
</script>

<style scoped>

label {
  display: block;
  font-size: larger;
}

select {
  margin: 10px 0;
}

.hostList div {
  border: solid 1px;
  padding: 10px;
  margin-bottom: 10px;
}

.error {
  border: solid 1px;
  padding: 10px;
  background-color: #ffd0d0;
}

</style>
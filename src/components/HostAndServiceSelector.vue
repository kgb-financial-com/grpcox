<template>

  <div class="error" v-if="!!errorMessage">{{ errorMessage }}</div>

  <div class="hostList" v-else>

    <div class="flexGroup">
      <div v-if="!!errorHostsMessage">{{ errorHostsMessage }}</div>
      <div v-else>
        <label for="hostDropDown">Host</label>
        <select id="hostDropDown" v-model="selectedHostInList">
          <option v-for="option in availableHosts" :key="option">{{ option }}</option>
        </select>
      </div>
      <div class="rightContainer">
        <label for="addHostGroup">Add and Select New Host</label>
        <div id="addHostGroup">
          <input v-model="inputName" placeholder="instance name" @keyup.enter="selectByUserInput()"/>
          <input v-model="inputHost" placeholder="hostname:port" @keyup.enter="selectByUserInput()"/>
          <button v-on:click="selectByUserInput()" type="submit">Add</button>
        </div>
      </div>
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
      selectedServiceInList: null,
      inputName: "",
      inputHost: ""
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

  watch: {
    selectedHostInList(newValue) {
      this.$store.commit("selectHost", {name: "<unknown>", host: newValue})
    },
    selectedServiceInList(newValue) {
      this.$store.commit("selectService", newValue)
    }
  },

  created() {
    this.refreshHostList(this.selectedName, this.selectedHost);
  },

  methods: {

    selectFirst(array) {
      return array && array.length ? array[0] : null;
    },

    selectSpecified(array, toSelect) {
      if (!array || !array.length) {
        return null;
      }
      const result = array.find(x => x === toSelect)
      return result ? result : array[0];
    },

    selectByUserInput() {
      this.refreshHostList(this.inputName, this.inputHost);
    },

    refreshHostList(name, host) {

      axios.get(this.$store.state.urlBase + "server/" + encodeURIComponent(host) + "/services")
          .then(data => {
            if (!data?.data?.data) {
              throw "Could not connect to host: " + host;
            }
            this.currentServices = data.data.data
                .map(s => s.substring(0, s.length - 1))
                .filter(s => !s.startsWith("grpc."))
            this.selectedServiceInList = this.selectFirst(this.currentServices);
          })
          .catch(err => this.errorMessage = "Could not connect to server " + name + " (" + host + "): " + err)
          .finally(() => axios.get(this.$store.state.urlBase + "active/get")
              .then(data => {
                this.availableHosts = data.data.data;
                this.selectedHostInList = this.selectSpecified(this.availableHosts, host);
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
  flex: 1 1 0;
}

.hostList div.flexGroup {
  border: none;
  padding: 0;
  display: flex;
}

#addHostGroup {
  margin-top: 10px;
}

.rightContainer {
  margin-left: 10px;
}

#addHostGroup input {
  margin-right: 10px;
}

.error {
  border: solid 1px;
  padding: 10px;
  background-color: #ffd0d0;
}

</style>
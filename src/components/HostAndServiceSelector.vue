<template>

  <div class="error" v-if="!!errorMessage">{{ errorMessage }}</div>

  <div class="hostList">

    <div class="flexGroup">
      <div class="error" v-if="!!errorHostsMessage">{{ errorHostsMessage }}</div>
      <div v-else>
        <label for="hostDropDown">Host</label>
        <div class="no-frame" id="hostDropDown">
          <select v-if="availableHosts.length" v-model="selectedHostInList">
            <option v-for="option in availableHosts" v-bind:value="option.value" :key="option.value">{{
                option.text
              }}
            </option>
          </select>
          <div v-else class="error under-label">Please add a host.</div>
        </div>
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
      <div class="no-frame" id="serviceDropDown">
        <select v-if="currentServices.length" v-model="selectedServiceInList">
          <option v-for="option in currentServices" :key="option">{{ option }}</option>
        </select>
        <div v-else class="error under-label">Please select a host.</div>
      </div>
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
    selectedHost() {
      return this.$store.state.selectedHost;
    }
  },

  watch: {
    selectedHostInList(newValue) {
      this.$store.commit("selectHost", newValue)
    },
    selectedServiceInList(newValue) {
      this.$store.commit("selectService", newValue)
    }
  },

  created() {
    this.refreshHostListForCurrentHost()
  },

  methods: {

    selectFirst(array) {
      return array && array.length ? array[0] : null;
    },

    selectSpecified(array, toSelect) {
      if (!array || !array.length) {
        return null;
      }
      const result = array.find(x => x.value === toSelect)
      return result ? result.value : array[0].value;
    },

    selectByUserInput() {
      const fullHostName = (this.inputName?.length ? this.inputName : this.inputHost) + '|' + this.inputHost;
      this.refreshHostList(fullHostName);
    },

    refreshHostListForCurrentHost() {
      this.refreshHostList(this.selectedHost);
    },

    refreshHostList(host) {

      if (host.endsWith("|null")) {
        this.getActiveHostsFromServer(null, this.refreshHostListForCurrentHost);
        return;
      }

      axios.get(this.$store.state.urlBase + "server/" + encodeURIComponent(host) + "/services")
          .then(data => {
            if (!data?.data?.data) {
              throw "Could not connect to host: " + host.split('|')[1];
            }
            this.currentServices = data.data.data
                .map(s => s.substring(0, s.length - 1))
                .filter(s => !s.startsWith("grpc."))
            this.selectedServiceInList = this.selectFirst(this.currentServices);
          })
          .catch(err => this.errorMessage = "Could not connect to server " + host.split('|')[0] + " (" + host.split('|')[1] + "): " + err)
          .finally(() => this.getActiveHostsFromServer(host, null));
    },

    getActiveHostsFromServer(host, next) {
      axios.get(this.$store.state.urlBase + "servers/get")
          .then(data => {
            this.availableHosts = this.createHostList(data.data.data);
            this.selectedHostInList = this.selectSpecified(this.availableHosts, host);
          })
          .catch(err => this.errorHostsMessage = "Could not retrieve list of available hosts: " + err.errorHostsMessage)
          .finally(() => {if (next) next()});
    },

    createHostList(serverList) {
      const exists = new Set();
      const multi = new Set();
      const result = serverList.map(server => {
        const nameAndHost = server.split('|');
        const name = nameAndHost.length === 2 ? nameAndHost[0] : server;
        return {
          name,
          text: null,
          value: server
        };
      });
      result.forEach(entry => {
        const name = entry.name;
        if (exists.has(name)) {
          multi.add(name);
        }
        exists.add(name);
      });
      return result.map(entry => {
        return {
          name: entry.name,
          text: entry.text = multi.has(entry.name)
              ? entry.name + ' - ' + entry.value.substr(entry.value.indexOf('|') + 1)
              : entry.name,
          value: entry.value
        }
      }).sort((a, b) => a.value.localeCompare(b.value));
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
  margin-bottom: 0;
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

.under-label {
  margin-top: 10px;
}

.hostList div.no-frame {
  border: none;
  padding: 0;
  margin-bottom: 0;
}

</style>

<template>
  <div class="outer">
    <div class="error" v-if="!!errorMessage">{{ errorMessage }}</div>
    <div class="row">
      <button @click="executeMethod">Execute</button>
      <label :for="id('method')">Method <span class="bold">{{ shortName }}</span></label>
    </div>
    <div class="method" :id="id('method')">
      <div :id="id('editor')" class="formatContainer">
        <label :for="id('editor')">Edit Input Parameters</label>
        <div class="ace-container">
          <v-ace-editor
              v-model:value="methodDescription.template"
              lang="json"
              theme="chrome"
              style="height: inherit"
              :options="{ showLineNumbers: false, showPrintMargin: false, highlightActiveLine: false, highlightGutterLine: false }"
          />
        </div>
      </div>
      <div :id="id('schema')" class="formatContainer">
        <label :for="id('schema')">Schema</label>
        <div class="ace-container">
          <v-ace-editor
              v-model:value="methodDescription.schema"
              lang="json"
              theme="chrome"
              style="height: inherit"
              :options="{ showLineNumbers: false, showPrintMargin: false, highlightActiveLine: false, readOnly: true, highlightGutterLine: false }"
          />
        </div>
      </div>
    </div>
    <div :id="id('result')" class="formatContainer" v-if="!!receivedResult">
      <label :for="id('result')">
        Result <span class="result-title-details">
        received at {{ receivedDateString }},
        request took {{ requestMillis }} ms</span>
      </label>
      <div class="ace-container">
        <v-ace-editor
            v-model:value="receivedResult"
            lang="json"
            theme="chrome"
            style="height: inherit"
            :options="{ showLineNumbers: false, showPrintMargin: false, highlightActiveLine: false, readOnly: true, highlightGutterLine: false }"
        />
      </div>
    </div>
    <div class="error" v-if="!!receiveErrorMessage">{{ receiveErrorMessage }}</div>
  </div>
</template>

<script>
import axios from "axios";
import {VAceEditor} from 'vue3-ace-editor';
import 'ace-builds/src-noconflict/mode-json';
import 'ace-builds/src-noconflict/theme-chrome';

export default {

  name: "Method",

  props: {
    name: String
  },

  components: {VAceEditor},

  data: function () {
    return {
      methodDescription: {
        schema: "",
        template: ""
      },
      errorMessage: null,
      receivedResult: null,
      receiveErrorMessage: null,
      receivedDate: null,
      requestMillis: 0,
    }
  },

  computed: {
    shortName() {
      const words = this.name.split(".");
      return words[words.length - 1];
    },
    selectedHost() {
      return this.$store.state.selectedHost;
    },
    receivedDateString() {
      if (!this.receivedDate?.toLocaleDateString) {
        return null;
      }
      return this.receivedDate.toLocaleString("en-US", { hour: "2-digit", minute: "2-digit", second: "2-digit", hour12: false});
    }
  },

  watch: {
    selectedService() {
      this.refreshMethodDetails();
    },
  },

  created() {
    this.refreshMethodDetails();
  },

  methods: {

    refreshMethodDetails() {

      if (!this.selectedHost) {
        return;
      }

      axios.get(this.$store.state.urlBase + "server/" + this.selectedHost + "/function/" + this.name + "/describe")
          .then(data => {
            if (!data?.data?.data) {
              throw "Could not connect to host: " + this.selectedHost.split('|')[1];
            }
            this.methodDescription = data.data.data;
          })
          .catch(err => this.errorMessage = "Could not connect to server " + this.selectedName.split('|')[0] + "(" + this.selectedHost.split('|')[1] + "): " + err)

    },

    executeMethod() {
      this.receivedResult = null;
      this.receiveErrorMessage = null;
      const startDate = Date.now();
      axios.post(this.$store.state.urlBase + "server/" + this.selectedHost + "/function/" + this.name + "/invoke", this.methodDescription.template)
          .then(data => {
            if (data?.data?.error) {
              throw data.data.error;
            }
            if (!data?.data?.data?.result) {
              throw "No data received";
            }
            this.receivedResult = data.data.data.result;
            const now = Date.now();
            this.requestMillis = now - startDate;
            this.receivedDate = new Date(now);
          })
          .catch(err => this.receiveErrorMessage = "Could not execute method on server " + this.selectedName.split('|')[0] + " (" + this.selectedHost.split('|')[1] + "): " + err)
    },

    id(kind) {
      return this.name + "-" + kind;
    }

  }

}
</script>

<style scoped>

.outer {
  border: solid 1px;
  padding: 10px;
  margin: 10px 0;
}

.method {
  display: flex;
}

.formatContainer {
  border: solid 1px;
  padding: 10px;
  margin: 10px;
  flex: 1 1 0;
}

.formatContainer:first-child {
  margin-left: 0;
}

.formatContainer:last-child {
  margin-left: 0;
}

.ace-container {
  margin-top: 5px;
  padding-top: 10px;
  border-top: solid 1px;
  height: 250px;
}

label {
  font-size: larger;
}

.bold {
  font-weight: bold;
}

.row {
  display: inline;
}

.row button {
  margin: 1px 10px 0 0;
}

.error {
  border: solid 1px;
  padding: 10px;
  background-color: #ffd0d0;
}

</style>

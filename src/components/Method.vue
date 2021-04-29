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
              :options="{ showLineNumbers: false, showPrintMargin: false, highlightActiveLine: false }"
          />
        </div>
      </div>
      <div :id="id('schema')" class="formatContainer">
        <label :for="id('schema')">Schema</label>
        <pre ref="prettifiedMethodDescription">{{ methodDescription.schema }}</pre>
      </div>
    </div>
    <div v-if="!!receivedResult">
      <label :for="id('result')">Result <span class="result-title-details">received at ... , request took ... ms</span></label>
      <div class="result" :id="id('result')">
        {{ receivedResult }}
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

const prettyPrint = require('code-prettify');

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
      receiveErrorMessage: null
    }
  },

  computed: {
    shortName() {
      const words = this.name.split(".");
      return words[words.length - 1];
    },
    selectedName() {
      return this.$store.state.selectedHost.name;
    },
    selectedHost() {
      return this.$store.state.selectedHost.host;
    }
  },

  watch: {
    selectedService() {
      this.refreshMethodDetails();
    },
    methodDescription() {
      this.$refs.prettifiedMethodDescription.classList = ["prettyprint"];
      this.$nextTick(() => prettyPrint.prettyPrint());
    }
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
              throw "Could not connect to host: " + this.selectedHost;
            }
            this.methodDescription = data.data.data;
          })
          .catch(err => this.errorMessage = "Could not connect to server " + this.selectedName + "(" + this.selectedHost + "): " + err)

    },

    executeMethod() {
      this.receivedResult = null;
      this.receiveErrorMessage = null;
      axios.post(this.$store.state.urlBase + "server/" + this.selectedHost + "/function/" + this.name + "/invoke", this.methodDescription.template)
          .then(data => {
            if (data?.data?.error) {
              throw data.data.error;
            }
            if (!data?.data?.data) {
              throw "No data received";
            }
            this.receivedResult = data.data.data;
          })
          .catch(err => this.receiveErrorMessage = "Could not execute method on server " + this.selectedName + " (" + this.selectedHost + "): " + err)
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

<style>

.ace_editor {
  height: 250px;
}

.prettyprint {
  margin-top: 5px;
  padding-top: 10px;
  border-top: solid 1px;
}

</style>
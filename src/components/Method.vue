<template>
  <div class="outer">
    <label :for="idMethod">Method: {{ shortName }}</label>
    <div class="method" :id="idMethod">
      <div :id="idEditor" class="formatContainer">
        <label :for="idEditor">Edit Input Parameters</label>
        <div class="ace-container">
        <v-ace-editor
            v-model:value="methodDescription.template"
            lang="json"
            theme="chrome"
            :options="{ showLineNumbers: false, showPrintMargin: false, highlightActiveLine: false }"
        />
        </div>
      </div>
      <div :id="idSchema" class="formatContainer">
        <label :for="idSchema">Schema</label>
        <pre ref="prettifiedMethodDescription">{{ methodDescription.schema }}</pre>
      </div>
    </div>
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
      errorMessage: null
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
    },
    idMethod() {
      return this.name + "-method";
    },
    idEditor() {
      return this.name + "-editor";
    },
    idSchema() {
      return this.name + "-schema";
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

</style>

<style>

.ace_editor {
  height: 300px;
}

.prettyprint {
  margin-top: 5px;
  padding-top: 10px;
  border-top: solid 1px;
}

</style>
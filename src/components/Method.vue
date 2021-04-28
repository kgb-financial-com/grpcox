<template>
  <div> &gt;&gt; {{ shortName }}</div>
  <div>{{ methodDescription.schema }}</div>
  <div>
    <v-ace-editor
        v-model:value="methodDescription.template"
        lang="json"
        theme="chrome"
        style="height: 300px; width: 300px;"
    />
  </div>
</template>

<script>
import axios from "axios";
import { VAceEditor } from 'vue3-ace-editor';
import 'ace-builds/src-noconflict/mode-json';
import 'ace-builds/src-noconflict/theme-chrome';

export default {

  name: "Method",

  props: {
    name: String
  },

  components: { VAceEditor },

  data: function() {
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
  },

  watch:{
    selectedService() {
      this.refreshMethodDetails();
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

</style>
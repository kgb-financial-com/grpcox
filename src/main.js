import { createApp } from 'vue'
import App from './App.vue'
import Vuex from 'vuex'

const getUrlBase = function() {  // must strictly return a string with a single trailing slash. Multiple slashes will cause CORS problems.
    let base = process.env.VUE_APP_URL_BASE;
    return base ? new URL(base) : new URL(".", document.baseURI);
}

const urlParams = new URL(location.href).searchParams;

const store = new Vuex.Store({
    state: {
        urlBase: getUrlBase(),
        selectedHost: {
            name: urlParams.get('name'),
            host: urlParams.get('host')
        },
        selectedService: null
    },
    mutations: {
        selectHost(state, name, host) {
            state.selectedHost = { name, host };
        },
        selectService(state, name) {
            state.selectedService = name;
        }
    }
});

const app = createApp(App)
app.use(store)
app.mount('#app')

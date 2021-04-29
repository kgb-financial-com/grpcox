import { createApp } from 'vue'
import { createStore } from 'vuex'
import App from './App.vue'

const getUrlBase = function() {  // must strictly return a string with a single trailing slash. Multiple slashes will cause CORS problems.
    let base = process.env.VUE_APP_URL_BASE;
    return base ? new URL(base) : new URL(".", document.baseURI);
}

const urlParams = new URL(location.href).searchParams;

const store = createStore({
    state() {
        return {
            urlBase: getUrlBase(),
            selectedHost: urlParams.get('name') + '|' + urlParams.get('host'),
            selectedService: null
        }
    },
    mutations: {
        selectHost(state, host) {
            state.selectedHost = host;
        },
        selectService(state, name) {
            state.selectedService = name;
        }
    }
});

const app = createApp(App)
app.use(store)
app.mount('#app')

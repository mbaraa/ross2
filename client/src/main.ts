import 'vuetify/styles';
import {createApp} from 'vue';
import {createVuetify} from 'vuetify';
import App from './App.vue';
import router from './router';
import VueGapi from 'vue-gapi';
import config from './config';
import {PublicClientApplication} from "@azure/msal-browser";
import store from './store'

const app = createApp(App);

app.config.globalProperties.$msal = new PublicClientApplication(config.msalConfig);

app.use(router)
    .use(createVuetify())
    .use(VueGapi, {
        clientId: config.googleClientID,
        scope: 'https://www.googleapis.com/auth/userinfo.email',
    })
    .use(store)
    .mount('#app');

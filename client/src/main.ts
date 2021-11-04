import 'vuetify/styles';
import {createApp} from 'vue';
import {createVuetify} from 'vuetify';
import App from './App.vue';
import router from './router';
import VueGapi from 'vue-gapi';
import config from './config';

createApp(App)
    .use(router)
    .use(createVuetify())
    .use(VueGapi, {
        clientId: config.googleClientID,
        scope: 'https://www.googleapis.com/auth/userinfo.email',
    })
    .mount('#app');

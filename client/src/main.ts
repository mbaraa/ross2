import 'vuetify/styles';
import {createApp} from 'vue';
import {createVuetify} from 'vuetify';
import App from './App.vue';
import router from './router';
import VueGapi from 'vue-gapi';

createApp(App)
    .use(router)
    .use(createVuetify())
    .use(VueGapi, {
        clientId: 'CLIENT_ID',
        scope: 'https://www.googleapis.com/auth/userinfo.email',
    })
    .mount('#app');

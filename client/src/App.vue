<template>
    <v-app>
        <v-app-bar app class="bg-indigo">
            <v-btn @click="toggleDrawer"
                   icon
                   color="error"
                   class="drawerButton">
                <FontAwesomeIcon class="text-indigo" :icon="{prefix: 'fas', iconName: 'bars'}"/>
                <FontAwesomeIcon class="text-red" v-if="newNotification" :icon="{prefix:'fas', iconName:'bell'}"/>
            </v-btn>
            <label @click="goHome" class="title">&nbsp;Ross 2</label>
        </v-app-bar>

        <v-navigation-drawer
            v-model="showDrawer"
            temporary
            app>
            <br/>

            <v-list @click="toggleDrawer" v-for="link in links" :key="link">
                <v-list-item>
                    <v-list-item-title>
                        <router-link :to="link.page">
                            <FontAwesomeIcon :class="getClass(link.name)" :icon="link.icon"/>
                            &nbsp;{{ link.name }}
                        </router-link>
                    </v-list-item-title>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>

        <v-main class="main">
            <router-view/>
        </v-main>
    </v-app>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faBars, faBell, faInfoCircle, faTrophy, faUserCircle, faGavel} from "@fortawesome/free-solid-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";
import Notification from "@/models/Notification";

library.add(faBars, faInfoCircle, faTrophy, faUserCircle, faBell, faGavel);
export default defineComponent({
    name: 'App',
    components: {
        FontAwesomeIcon,
    },
    data() {
        return {
            showDrawer: false,
            links: [
                {page: '/', name: "Contests", icon: {prefix: 'fas', iconName: 'trophy'}},
                {page: '/notifications', name: "Notifications", icon: {prefix: 'fas', iconName: 'bell'}},
                {page: '/profile', name: "Profile", icon: {prefix: 'fas', iconName: 'user-circle'}},
                {page: 'organizer', name: "Organizer", icon: {prefix: 'fas', iconName: 'gavel'}},
                {page: '/about', name: "About", icon: {prefix: 'fas', iconName: 'info-circle'}},
            ],
            newNotification: false,
        };
    },
    async mounted() {
        this.newNotification = await Notification.checkNotifications();
    },
    methods: {
        toggleDrawer(): void {
            this.showDrawer = !this.showDrawer;
        },
        goHome() {
            this.$router.push("/");
        },
        getClass(name: string): string {
            return this.newNotification && name == "Notifications"? "text-red": "";
        }
    }
})
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Ropa+Sans&display=swap");

.title {
    font-family: 'Ropa Sans', sans-serif;
    color: white;
    font-weight: bold;
    font-size: 2em;
    cursor: pointer;
}

.drawerButton {
    background: aliceblue;
}

a {
    font-family: 'Ropa Sans', sans-serif;
    font-weight: bold;
    text-decoration: none;
    color: #2c3e50;
    font-size: 1.5em;
}

a.router-link-exact-active {
    color: #3f51b5;
}

.main {
    font-family: 'Ropa Sans', sans-serif;
    color: #3f51b5;
}
</style>

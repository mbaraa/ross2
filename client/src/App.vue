<template>
    <v-app class="app" >
        <v-app-bar elevation="0" style="border-bottom-width: 1px; border-color: #e4e4e4;" app>
            <img @click="goHome" src="https://mbaraa.fun/ross2/logo_250.png" alt="Ross 2" class="logo"/>

            <div class="links">
                <router-link :to="links[0].page">
                    <FontAwesomeIcon :class="getClass(links[0].name)" :icon="links[0].icon"/>
                        &nbsp;
                </router-link>
                <router-link :to="links[1].page">
                    <FontAwesomeIcon :class="getClass(links[1].name)" :icon="links[1].icon"/>
                        &nbsp;
                </router-link>
                <router-link :to="links[2].page">
                    <FontAwesomeIcon :class="getClass(links[2].name)" :icon="links[2].icon"/>
                        &nbsp;
                </router-link>
                <router-link :to="links[3].page">
                    <FontAwesomeIcon :class="getClass(links[3].name)" :icon="links[3].icon"/>
                        &nbsp;
                </router-link>
            </div>
        </v-app-bar>

        <v-navigation-drawer v-model="showDrawer" temporary app>
            <br/>

            <v-list @click="toggleDrawer" v-for="link in links" :key="link">
                <v-list-item>
                    <v-list-item-title>
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
import {faBars, faBell, faGavel, faInfoCircle, faTrophy, faUserCircle} from "@fortawesome/free-solid-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";
import NotificationRequests from "@/utils/requests/NotificationRequests";

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
                {page: '/about', name: "About", icon: {prefix: 'fas', iconName: 'info-circle'}},
            ],
            newNotification: false,
        };
    },
    async mounted() {
        this.newNotification = await NotificationRequests.checkNotifications();
    },
    methods: {
        toggleDrawer(): void {
            this.showDrawer = !this.showDrawer;
        },
        goHome() {
            this.$router.push("/");
        },
        getClass(name: string): string {
            return this.newNotification && name == "Notifications" ? "text-red" : "";
        }
    }
})
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Ropa+Sans&display=swap");

.title {
    font-family: 'Ropa Sans', sans-serif;
    color: #6a1b9a;
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
    color: #4A148C;
}

.main {
    font-family: "Ropa Sans", sans-serif;
    color: #311B92;
}

.app {
    
}

.logo {
    height: 3rem;
    width: 3rem;
    background-color: #E4E4E4;
    border: 0.5px;
    border-radius: 100%;
    margin-right: 0.5rem;
    cursor: pointer;
    display: inline;
}

.logo:hover {
    opacity: 0.6;
}

.links {
    position: absolute;
    right: 0.5rem;
}
</style>

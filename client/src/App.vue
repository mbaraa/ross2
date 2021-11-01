<template>
    <v-app>
        <v-app-bar app class="bg-indigo">
            <v-btn @click="toggleDrawer"
                   icon
                   color="error"
                   class="drawerButton">
                <FontAwesomeIcon class="text-indigo" :icon="{prefix: 'fas', iconName: 'bars'}"/>
            </v-btn>
            <a href="/#" class="title">&nbsp;Ross 2</a>
        </v-app-bar>

        <v-navigation-drawer
            v-model="showDrawer"
            temporary
            app>

            <v-list @click="toggleDrawer">
                <br/>
                <v-list-item>
                    <v-list-item-title>
                        <router-link to="/">
                            <FontAwesomeIcon :icon="{prefix: 'fas', iconName: 'trophy'}"/>
                            &nbsp;Contests
                        </router-link>
                    </v-list-item-title>
                </v-list-item>

                <v-list-item>
                    <v-list-item-title>
                        <router-link to="/notifications">
                            <FontAwesomeIcon :icon="{prefix: 'fas', iconName: 'bell'}"/>
                            &nbsp;Notifications
                        </router-link>
                    </v-list-item-title>
                </v-list-item>

                <v-list-item>
                    <v-list-item-title>
                        <router-link to="/profile">
                            <FontAwesomeIcon :icon="{prefix: 'fas', iconName: 'user-circle'}"/>
                            &nbsp;Profile
                        </router-link>
                    </v-list-item-title>
                </v-list-item>

                <v-list-item>
                    <v-list-item-title>
                        <router-link to="/about">
                            <FontAwesomeIcon :icon="{prefix: 'fas', iconName: 'info-circle'}"/>
                            &nbsp;About
                        </router-link>
                    </v-list-item-title>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>

        <v-main class="main">
            <router-view></router-view>
        </v-main>
    </v-app>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faBars, faInfoCircle, faTrophy, faUserCircle, faBell} from "@fortawesome/free-solid-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";

library.add(faBars, faInfoCircle, faTrophy, faUserCircle, faBell);
export default defineComponent({
    name: 'App',
    components: {
        FontAwesomeIcon,
    },
    data() {
        return {
            showDrawer: false,
        };
    },
    methods: {
        toggleDrawer(): void {
            this.showDrawer = !this.showDrawer;
        }
    }
})
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Ropa+Sans&display=swap');

.title {
    font-family: 'Ropa Sans', sans-serif;
    color: white;
    font-weight: bold;
    font-size: 2em;
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

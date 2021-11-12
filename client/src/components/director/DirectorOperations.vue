<template>
    <div class="main2" v-if="(director.roles & 1) !== 0">
        <div v-if="showOps">
            <div class="pagesLinks" v-for="link in links" :key="link">
                <router-link :to="{name: link.page}">
                    {{ link.name }}
                </router-link>
            </div>
            <div>
                <router-view/>
            </div>
        </div>
        <v-btn @click="show" v-if="!showOps">Show Operations</v-btn>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Organizer from "@/models/Organizer";

export default defineComponent({
    name: "DirectorOperations",
    props: {
        director: Organizer
    },
    data() {
        return {
            profile: {},
            showOps: false,
            links: [
                {page: 'contests', name: 'Contests'},
                {page: 'organizers', name: 'Organizers'},
                {page: 'other', name: 'Other'},
            ]
        }
    },
    methods: {
        show() {
            this.showOps = true;
            this.$router.push('/organizer/contests/')
        },
        checkDirector(): boolean {
            return (this.director.roles & 1) != 0;
        }
    }
});
</script>

<style scoped>
.main2 {
    padding: 10px;
    text-align: center;
}

.pagesLinks {
    border-radius: 5px 5px 0 0;
    padding: 0px;
    display: inline;
    /*padding-right: 10px;*/
}

/* inactive subpage */
.pagesLinks a {
    font-weight: bold;
    padding: 10px;

    text-decoration: none;
    background-color: #212121;
    color: #e0e0e0;
    border: #212121 solid;
    border-bottom: white solid 0;
}

/* active subpage */
.pagesLinks a.router-link-exact-active {
    background-color: #e0e0e0;
    color: #212121;
}
</style>

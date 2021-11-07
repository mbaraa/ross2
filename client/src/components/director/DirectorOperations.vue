<template>
    <div class="main">
        <span class="contestName"><b>{{ profile.name }}</b></span>
        <br/><br/>
        <div>
            <div class="pagesLinks">
                <router-link :to="{name: 'contests'}">
                    Contests
                </router-link>
                <router-link :to="{ name:'organizers'}">
                    Organizers
                </router-link>
                <router-link :to="{ name:'other'}">
                    Other
                </router-link>
            </div>
            <div class="subpage">
                <router-view/>
            </div>
        </div>
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
            profile: {}
        }
    }
});
</script>

<style scoped>
.main {
    padding: 10px;
    text-align: center;
}

.contestLogo {
    border-radius: 100%;
    width: 125px;
    height: 125px;
}

.contestName {
    font-size: 2em;
    color: #212121;
}

.pagesLinks {
    border-radius: 5px 5px 0 0;
    padding-left: 10px;
    padding-right: 10px;
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

.subpage {
    margin: 6px auto;
    border-radius: 0 0 5px 5px;
    padding-top: 10px;

    display: block;

    width: 90%;
    border: #212121 solid 2px;
    border-top: white solid 0;
    background-color: #e0e0e0;
}
</style>

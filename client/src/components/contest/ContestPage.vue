<template>
    <div class="main">
        <img class="contestLogo" alt="contest logo" :src="config.backendAddress+contest.logo_path"/>
        <br/>
        <span class="contestName"><b>{{ contest.name }}</b></span>
        <br/><br/>
        <div>
            <div class="pagesLinks">
                <router-link :to="{ name:'details', query: { id:getContestId()} }">
                    Details
                </router-link>
                <router-link v-if="!contest.teams_hidden" :to="{ name:'teams', query: { id:getContestId()} }">
                    Teams
                </router-link>
                <!--            <router-link :to="{ name:'teamless', query: { id:getContestId()} }">-->
                <!--                Teamless-->
                <!--            </router-link>-->
                <router-link :to="{ name:'support', query: { id:getContestId()} }">
                    Support
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
import Contest from "@/models/Contest";
import config from "@/config";

export default defineComponent({
    name: "ContestPage",
    data() {
        return {
            contest: {},
            config: config
        }
    },
    async mounted() {
        this.contest = await Contest.getContestFromServer(this.$route.query.id);
    },
    methods: {
        getContestId(): number {
            return this.contest.id;
        },
    },

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
}

</style>

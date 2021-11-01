<template>
    <table class="main" v-if="teamless.length > 0">
        <tr>
            <th>Contestant name</th>
            <th>University ID</th>
            <th></th>
        </tr>
        <tr v-for="contestant in teamless" :key="contestant" :class="getContClass(contestant)">
            <td>{{ contestant.name }}</td>
            <td>{{ contestant.university_id }}</td>
            <td>
                <v-btn @click="inviteContestant(contestant)">Invite</v-btn>
            </td>
        </tr>
    </table>
    <h1 v-else>There are no teamless contestants for this contest 🤗</h1>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Contest from "@/models/Contest";
import Contestant from "@/models/Contestant";

export default defineComponent({
    name: "ContestTeamless",
    data() {
        return {
            contest: {},
            teamless: []
        }
    },
    async mounted() {
        this.contest = await Contest.getContestFromServer(this.$route.query.id);
        this.teamless = this.contest.teamless_contestants;
    },
    methods: {
        async inviteContestant(cont: Contestant) {
            window.alert("not ready yet :(")
            return;
        },
        getContClass(cont: Contestant): string {
            return this.teamless.indexOf(cont) % 2 == 0 ? "cont1" : "cont2";
        }
    }
});
</script>

<style scoped>
.main {
    width: 100%;
    border-collapse: collapse;
}

td, th {
    padding: 5px;
    border: #212121 solid 2px;
}

.cont1 {
    background-color: #d0d0d0;
}

.cont2 {
    background-color: #a0a0a0;
}
</style>

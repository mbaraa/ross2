<template>
    <table class="main" v-if="teams.length > 0">
        <tr>
            <th>Team name</th>
            <th>Team members</th>
            <th></th>
        </tr>
        <tr v-for="team in teams" :key="team" :class="getTeamClass(team)">
            <td>{{ team.name }}</td>
            <td>{{ team.members.join(", ") }}</td>
            <td>
                <v-btn @click="joinTeam(team)">Join team</v-btn>
            </td>
        </tr>
    </table>
    <h1 v-else>No teams are registered for this contest so far!</h1>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Contest from "@/models/Contest";
import config from "@/config";
import Team from "@/models/Team";

export default defineComponent({
    name: "ContestTeams",
    data() {
        return {
            contest: {},
            teams: []
        }
    },
    async mounted() {
        this.contest = await Contest.getContestFromServer(this.$route.query.id);
        this.teams = this.contest.teams;
    },
    methods: {
        async joinTeam(team: Team) {
            const resp = await fetch(`${config.backendAddress}/contestant/req-join-team/`, {
                method: "POST",
                mode: "cors",
                headers: {
                    "Authorization": <string>localStorage.getItem("token")
                },
                body: JSON.stringify(team)
            })

            if (resp.ok) {
                window.alert("request sent successfully!");
            } else {
                window.alert(`${resp.status} ${resp.statusText}!`);
            }
        },
        getTeamClass(team: Team): string {
            return this.teams.indexOf(team) % 2 == 0 ? "team1" : "team2";
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

.team1 {
    background-color: #d0d0d0;
}

.team2 {
    background-color: #a0a0a0;
}
</style>

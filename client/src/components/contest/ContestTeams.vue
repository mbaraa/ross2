<template>
    <table class="main" v-if="teams.length > 0">
        <tr>
            <th>Team name</th>
            <th>Team members</th>
            <th></th>
        </tr>
        <tr v-for="team in teams" :key="team" :class="getTeamClass(team)">
            <td>{{ team.name }}</td>
            <td>{{ getMembersNames(team) }}</td>
            <td>
                <v-btn v-if="!team.inTeam" @click="joinTeam(team)">Join team</v-btn>
            </td>
        </tr>
    </table>
    <h1 v-else>No teams are registered for this contest so far!</h1>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Contest from "@/models/Contest";
import Team from "@/models/Team";
import Contestant from "@/models/Contestant";
import JoinRequest from "@/models/JoinRequest";

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
        this.processInTeam();
    },
    methods: {
        async joinTeam(team: Team) {
            const resp = await Contestant.requestJoinTeam(<JoinRequest> {
                requested_team: team,
                requested_team_id: team.id,
                request_message: "",
            })

            if (resp.ok) {
                window.alert("request sent successfully!");
                team.inTeam = true;
            } else {
                window.alert(`${resp.status} ${resp.statusText}!`);
            }
        },
        getMembersNames(team: Team): string {
            let names = "";
            team.members.forEach((cont: Contestant) => {
                names += cont.name + ", ";
            });
            return names.substring(0, names.length-2);
        },
        processInTeam() {
            this.teams.forEach(async (team: Team) => {
                team.inTeam = await Contestant.checkJoinedTeam(team);
            });
        },
        getTeamClass(team: Team): string {
            return this.teams.indexOf(team) % 2 == 0 ? "team1" : "team2";
        },
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

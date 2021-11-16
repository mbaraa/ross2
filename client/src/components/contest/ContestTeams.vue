<template>
    <div v-if="teams.length > 0">
        <div class="main bg-green-accent-4 team" v-for="team in teams" :key="team">
            <TeamCard :team="team"/>
            <v-btn v-if="!team.inTeam" @click="joinTeam(team)">Join team</v-btn>
        </div>
    </div>
    <div v-else>
        <h1>No teams are registered for this contest so far!</h1>
        <h2 title="go to the home page and click on `create team` under the contest description :)">
            Be the first to register 😁
        </h2>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Contest from "@/models/Contest";
import Team from "@/models/Team";
import Contestant from "@/models/Contestant";
import JoinRequest from "@/models/JoinRequest";
import TeamCard from "@/components/team/TeamCard.vue";

export default defineComponent({
    name: "ContestTeams",
    components: {TeamCard},
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
            const resp = await Contestant.requestJoinTeam(<JoinRequest>{
                requested_team: team,
                requested_team_id: team.id,
                request_message: "",
                requested_contest_id: this.contest.id,
                requested_contest: this.contest,
            })

            if (resp.ok) {
                window.alert("request sent successfully!");
                team.inTeam = true;
            } else {
                window.alert(`${resp.status} ${resp.statusText}!`);
            }
        },
        processInTeam() {
            this.teams.forEach(async (team: Team) => {
                team.inTeam = await Contestant.checkJoinedTeam(team);
            });
        },
    }
});
</script>

<style scoped>
.main {
    color: white;
    text-align: center;
    margin: 10px auto;
}

.team {
    display: inline-grid;
    padding: 20px;
    margin: 10px;
    border-radius: 5px;
}
</style>

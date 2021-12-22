<template>
    <div v-if="contest.teams_hidden">
        <h1>This contest's teams are private!</h1>
    </div>
    <div v-else-if="teams.length > 0">
        <div class="main bg-green-accent-4 team" v-for="team in teams" :key="team">
            <TeamCard :team="team"/>
            <v-btn v-if="checkTeam(team)" @click="joinTeam(team)">Join team</v-btn>
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
import JoinRequest from "@/models/JoinRequest";
import TeamCard from "@/components/team/TeamCard.vue";
import ContestantRequests from "@/utils/requests/ContestantRequests";

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
        this.contest = await Contest.getContestFromServer(this.$route.query["id"]);
        this.teams = this.contest.teams;
        this.processInTeam();
    },
    methods: {
        checkRegisterEnds(): boolean {
            const regOver = (new Date().getTime()) > this.contest.registration_ends;
            if (regOver) {
                window.alert("sorry, the registration for this contest is over!")
            }

            return regOver;
        },
        async joinTeam(team: Team) {
            if (this.checkRegisterEnds()) {
                return;
            }
            console.log("team", team)
            const resp = await ContestantRequests.requestJoinTeam(<JoinRequest>{
                requested_team: team,
                requested_team_id: team.id,
                requested_team_join_id: team.join_id,
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
                team.inTeam = await ContestantRequests.checkJoinedTeam(team);
            });
        },
        checkTeam(team: Team): boolean {
            return !team.inTeam && team.members.length < this.contest.participation_conditions.max_team_members;
        }
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

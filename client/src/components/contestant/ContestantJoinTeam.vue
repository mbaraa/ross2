<template>
    <v-dialog
        max-height="400"
        max-width="-40"
        transition="dialog-bottom-transition"
        v-model="dialog"
        scrollable>
        <template v-slot:activator="{on, attrs}">
            <div v-bind="attrs" v-on="on" @click="checkTokenForAction(openDialog)" style="display: inline">
                <v-btn color="warning" title="create team">
                    <FontAwesomeIcon class="text-white" :icon="{prefix:'fas', iconName:'users'}"/>&nbsp;
                    Join a team
                </v-btn>
            </div>
        </template>

        <v-card elevation="16" class="teamForm">
            <v-card-title>
                <span class="text-h4">Join team</span>
            </v-card-title>
            <p>Ask the member who created the team to give you the team id</p>
            <v-text-field label="Team id" v-model="team.join_id" autofocus/>

            <v-btn class="bg-red" @click="dialog = false">
                Close
            </v-btn>&nbsp;
            <v-btn class="bg-blue" @click="joinTeam">
                Join
            </v-btn>
        </v-card>
    </v-dialog>
</template>

<script lang="ts">

import {defineComponent} from "vue";
import Team from "@/models/Team";
import ContestantRequests from "@/utils/requests/ContestantRequests";
import JoinRequest from "@/models/JoinRequest";
import Contest from "@/models/Contest";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {library} from "@fortawesome/fontawesome-svg-core";
import {faUserPlus} from "@fortawesome/free-solid-svg-icons";
import ActionChecker from "@/utils/ActionChecker";

library.add(faUserPlus);

export default defineComponent({
    name: "ContestantJoinTeam",
    props: {
        contest: Contest,
    },
    components: {
        FontAwesomeIcon
    },
    data() {
        return {
            dialog: false,
            team: new Team(),
        };
    },
    methods: {
        async checkTokenForAction(fn: () => void) {
            await ActionChecker.checkContestant(fn);
        },
        openDialog() {
            this.dialog = true;
        },
        async joinTeam() {
            const resp = await ContestantRequests.requestJoinTeam(<JoinRequest>{
                requested_team: this.team,
                requested_team_id: +this.team.id,
                requested_team_join_id: this.team.join_id,
                request_message: "",
                requested_contest_id: +this.contest.id,
                requested_contest: this.contest,
            })

            if (resp.ok) {
                window.alert("request sent successfully, now wait for the team's leader to accept your request!");
                this.team.inTeam = true;
            } else {
                window.alert(`${resp.status} ${resp.statusText}!`);
            }
            window.location.reload();
        }
    }
});
</script>

<style scoped>
.teamForm {
    padding: 10px;
    margin: 0 auto;
    width: 400px;
    overflow-y: auto;
}
</style>

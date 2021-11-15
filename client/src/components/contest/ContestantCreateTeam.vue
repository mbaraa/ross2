<template>
    <v-dialog
        max-height="400"
        max-width="-40"
        transition="dialog-bottom-transition"
        v-model="dialog"
        scrollable>
        <template v-slot:activator="{on, attrs}">
            <div v-bind="attrs" v-on="on" @click="checkTokenForAction(openDialog)" style="display: inline">
                <v-btn icon color="error" title="create team">
                    <FontAwesomeIcon class="text-white" :icon="{prefix:'fas', iconName:'user-plus'}"/>
                </v-btn>
            </div>
        </template>

        <v-card elevation="16" class="teamForm">
            <v-card-title>
                <span class="text-h4">Create Team</span>
            </v-card-title>

            <v-text-field label="Team name" v-model="team.name"/>

            <v-btn class="bg-red" @click="dialog = false">
                Close
            </v-btn>&nbsp;
            <v-btn class="bg-blue" @click="createTeam">
                Create
            </v-btn>
        </v-card>
    </v-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {library} from "@fortawesome/fontawesome-svg-core";
import {faUserPlus} from "@fortawesome/free-solid-svg-icons";
import Contest from "@/models/Contest";
import Team from "@/models/Team";
import Contestant from "@/models/Contestant";
import {checkTokenForAction} from "@/utils";

library.add(faUserPlus);

export default defineComponent({
    name: "ContestantCreateTeam",
    components: {
        FontAwesomeIcon
    },
    data() {
        return {
            dialog: false,
            team: new Team(),
        }
    },
    props: {
        contest: Contest
    },
    methods: {
        async createTeam() {
            this.team.contests.push(this.contest);
            await Contestant.createTeam(this.team);
            this.dialog = false;

            window.alert(`your team "${this.team.name}" was created successfully ☺️`);
            window.location.reload();
        },
        checkTokenForAction(fn: () => void) {
            checkTokenForAction(fn);
        },
        openDialog() {
            this.dialog = true;
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

<template>
    <v-dialog
        max-height="400"
        max-width="-40"
        transition="dialog-bottom-transition"
        v-model="dialog"
        scrollable>
        <template v-slot:activator="{on, attrs}">
            <div v-bind="attrs" v-on="on" @click="checkTokenForAction(checkRegistrationEndsForAction(openDialog))" style="display: inline">
                <v-btn color="success" title="join as teamless">
                    <FontAwesomeIcon class="text-white" :icon="{prefix:'fas', iconName:'users-slash'}"/>&nbsp;
                    Join as teamless
                </v-btn>
            </div>
        </template>

        <v-card elevation="16" class="teamForm">
            <v-card-title>
                <span class="text-h4">Join as teamless</span>
            </v-card-title>

            <label style="font-size: 1.2em">Select your gender: </label>

            <br/>
            <input type="radio" id="male" value="true" v-model="contestantProfile.gender">
            <label for="male">Male</label>
            &nbsp;
            <input type="radio" id="female" value="false" v-model="contestantProfile.gender">
            <label for="female">Female</label>
            <br/>

            <label style="font-size: 1.2em">Do you mind participating with the other gender? </label>

            <br/>
            <input type="radio" id="yes" value="true" v-model="contestantProfile.participate_with_other">
            <label for="yes">Yes, I mind</label>
            &nbsp;
            <input type="radio" id="no" value="false" v-model="contestantProfile.participate_with_other">
            <label for="no">No, I don't mind</label>
            <br/><br/>

            <v-btn class="bg-red" @click="dialog = false">
                Close
            </v-btn>&nbsp;
            <v-btn class="bg-blue" @click="checkTokenForAction(checkRegistrationEndsForAction(joinAsTeamless))">
                Join as teamless
            </v-btn>
        </v-card>
    </v-dialog>
</template>

<script lang="ts">

import {defineComponent} from "vue";
import ContestantRequests from "@/utils/requests/ContestantRequests";
import Contest from "@/models/Contest";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {library} from "@fortawesome/fontawesome-svg-core";
import {faUserSlash} from "@fortawesome/free-solid-svg-icons";
import ActionChecker from "@/utils/ActionChecker";

library.add(faUserSlash);

export default defineComponent({
    name: "ContestantJoinTeamless",
    props: {
        contest: Contest,
    },
    components: {
        FontAwesomeIcon
    },
    data() {
        return {
            dialog: false,
            contestantProfile: null,
        };
    },
    async mounted() {
        this.contestantProfile = await ContestantRequests.getProfile();
    },
    methods: {
        async checkTokenForAction(fn: () => void) {
            await ActionChecker.checkContestant(fn);
        },
        openDialog() {
            this.dialog = true;
        },
        async joinAsTeamless() {
            if (window.confirm(`are you sure you want to join the contest "${this.contest.name}" as teamless?`)) {
                this.contestantProfile.gender = (this.contestantProfile.gender == "true");
                this.contestantProfile.participate_with_other = (this.contestantProfile.participate_with_other == "true");
                await ContestantRequests.joinAsTeamless({
                    contest: this.contest,
                    contestant: this.contestantProfile,
                });
                window.alert(`you have registered as teamless in "${this.contest.name}"`);
                this.dialog = false;
            }
        },
        checkRegisterEnds(): boolean {
            const regOver = (new Date()).getTime() > this.contest.registration_ends;
            if (regOver) {
                window.alert("sorry, the registration for this contest is over!")
            }

            return regOver;
        },
        checkRegistrationEndsForAction(fn: () => void): () => void {
            return !this.checkRegisterEnds() ? fn : () => {
                const _ = true
            };
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

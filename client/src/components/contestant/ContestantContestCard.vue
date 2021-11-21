<template>
    <div class="main bg-blue-darken-4">
        <ContestCard :contest="contest"/>
        <v-divider style="padding-bottom: 10px"/>
        <div class="buttons">
            <ContestantCreateTeam v-if="!hasTeam" class="delete" :contest="contest"/>
            &nbsp;
            <v-btn @click="checkTokenForAction(joinTeam)" icon color="warning" class="delete" title="join a team">
                <FontAwesomeIcon class="text-white" :icon="{prefix:'fas', iconName:'users'}"/>
            </v-btn>
            &nbsp;
            <v-btn @click="checkTokenForAction(checkRegistrationEndsForAction(joinAsTeamless))" icon color="success" class="delete"
                   title="join as teamless">
                <FontAwesomeIcon class="text-white" :icon="{prefix:'fas', iconName:'users-slash'}"/>
            </v-btn>
        </div>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faUserPlus, faUsers, faUsersSlash} from "@fortawesome/free-solid-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";
import ContestCard from "@/components/contest/ContestCard.vue";
import Contest from "@/models/Contest";
import ContestantCreateTeam from "@/components/contestant/ContestantCreateTeam.vue";
import {checkTokenForAction} from "@/utils";
import ContestantRequests from "@/utils/requests/ContestantRequests";

library.add(faUserPlus, faUsers, faUsersSlash);

export default defineComponent({
    name: "ContestantContestCard",
    props: {
        contest: Contest
    },
    components: {
        ContestantCreateTeam,
        ContestCard,
        FontAwesomeIcon
    },
    data() {
        return {
            contestant: null,
            hasTeam: false
        }
    },
    async mounted() {
        this.contestant = await ContestantRequests.login();
        this.hasTeam = ((this.contestant) != null && (this.contestant).team_id > 1);
    },
    methods: {
        checkRegisterEnds(): boolean {
            const regOver = (new Date()).getTime() > this.contest.registration_ends;
            if (regOver) {
                window.alert("sorry, the registration for this contest is over!")
            }

            return regOver;
        },
        async joinAsTeamless() {
            if (window.confirm(`are you sure you want to join the contest "${this.contest.name}" as teamless?`)) {
                await ContestantRequests.joinAsTeamless(this.contest);
                window.alert(`you have registered as teamless in "${this.contest.name}"`)
            }
        },
        joinTeam() {
            this.$router.push(`/contest/teams/?id=${this.contest.id}`)
        },
        checkTokenForAction(fn: () => void) {
            checkTokenForAction(fn);
        },
        checkRegistrationEndsForAction(fn: () => void): () => void {
            return !this.checkRegisterEnds()? fn: () => {const _ = true};
        }
    }
});
</script>

<style scoped>
.main {
    cursor: pointer;
    border-radius: 5px;
}

.buttons {
    padding-bottom: 10px;
}
</style>

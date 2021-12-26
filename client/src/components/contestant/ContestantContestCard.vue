<template>
    <div class="main bg-blue-darken-4">
        <ContestCard :contest="contest"/>
        <v-divider style="padding-bottom: 10px"/>
        <div>
            <div class="buttons" v-if="!hasTeam">
                <ContestantCreateTeam :contest="contest"/>
            </div>
            <br/>
            <div class="buttons" v-if="!hasTeam">
                <ContestantJoinTeam :contest="contest"/>
            </div>
            <div class="buttons" v-if="!hasTeam">
                <ContestantJoinTeamless :contest="contest"/>
            </div>
<!--            <v-btn v-if="!hasTeam" @click="checkTokenForAction(checkRegistrationEndsForAction(joinAsTeamless))"-->
<!--                   color="success"-->
<!--                   class="buttons"-->
<!--                   title="you will be put in a team at the end of registration">-->
<!--                <FontAwesomeIcon class="text-white" :icon="{prefix:'fas', iconName:'users-slash'}"/>&nbsp;-->
<!--                Join as teamless-->
<!--            </v-btn>-->
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
import ContestantRequests from "@/utils/requests/ContestantRequests";
import ContestantJoinTeam from "@/components/contestant/ContestantJoinTeam.vue";
import ActionChecker from "@/utils/ActionChecker";
import ContestantJoinTeamless from "@/components/contestant/ContestantJoinTeamless.vue";

library.add(faUserPlus, faUsers, faUsersSlash);

export default defineComponent({
    name: "ContestantContestCard",
    props: {
        contest: Contest
    },
    components: {
        ContestantJoinTeamless,
        ContestantJoinTeam,
        ContestantCreateTeam,
        ContestCard,
    },
    data() {
        return {
            contestant: null,
            hasTeam: false
        }
    },
    async mounted() {
        this.contestant = await ContestantRequests.getProfile();
        this.hasTeam = ((this.contestant) != null && (this.contestant.team_id > 1 || this.contestant.teamless_contest_id > 0));
    },
    methods: {
        checkRegisterEnds(): boolean {
            const regOver = (new Date()).getTime() > this.contest.registration_ends;
            if (regOver) {
                window.alert("sorry, the registration for this contest is over!")
            }

            return regOver;
        },
        joinTeam() {
            this.$router.push(`/contest/teams/?id=${this.contest.id}`)
        },
        async checkTokenForAction(fn: () => void) {
            await ActionChecker.checkContestant(fn);
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
.main {
    cursor: pointer;
    border-radius: 5px;
}

.buttons {
    display: inline-block;
    padding: 5px;
    margin-bottom: 5px;
}
</style>

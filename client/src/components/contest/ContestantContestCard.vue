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
            <v-btn @click="checkTokenForAction(joinAsTeamless)" icon color="success" class="delete" title="join as teamless">
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
import Contestant from "@/models/Contestant";
import ContestantCreateTeam from "@/components/contest/ContestantCreateTeam.vue";
import {checkTokenForAction} from "@/utils";

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
        this.contestant = await Contestant.login();
        this.hasTeam = ((this.contestant) != null && (this.contestant).team_id > 0) ;
    },
    methods: {
        async joinAsTeamless() {
            if (window.confirm(`are you sure you want to join the contest "${this.contest.name}" as teamless?`)) {
                await Contestant.joinAsTeamless(this.contest);
                window.alert(`you have registered as teamless in "${this.contest.name}"`)
            }
        },
        joinTeam() {
            this.$router.push(`/contest/teams/?id=${this.contest.id}`)
        },
        checkTokenForAction(fn: () => void) {
            checkTokenForAction(fn);
        },
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

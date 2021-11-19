<template>
    <div class="main bg-blue-darken-4">
        <ContestCard :contest="contest" @click="openContestDetails"/>
        <div class="buttons">
            <v-btn title="delete contest" @click="deleteContest" icon color="error">
                <FontAwesomeIcon class="text-white" :icon="{prefix:'fas', iconName:'trash'}"/>
            </v-btn>
            &nbsp;
            <v-btn title="generate teams for team less contestants" @click="generateTeams(contest)" icon
                   color="success">
                <FontAwesomeIcon class="text-white" :icon="{prefix:'fas', iconName:'cogs'}"/>
            </v-btn>
            &nbsp;
            <v-btn title="send society service forms to contestants" @click="finishContest(contest)" icon
                   color="info">
                <FontAwesomeIcon class="text-white" :icon="{prefix:'fas', iconName:'calendar-check'}"/>
            </v-btn>
        </div>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Contest from "@/models/Contest.ts";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faCogs, faTrash, faCalendarCheck} from "@fortawesome/free-solid-svg-icons"
import {library} from "@fortawesome/fontawesome-svg-core";
import ContestCard from "@/components/contest/ContestCard.vue";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";

library.add(faTrash, faCogs, faCalendarCheck);

export default defineComponent({
    name: "DirectorContestCard",
    props: {
        contest: Contest
    },
    components: {
        ContestCard,
        FontAwesomeIcon,
    },
    methods: {
        openContestDetails() {
            this.$router.push(`/contest/details/?id=${this.contest.id}`);
        },
        async deleteContest() {
            if (window.confirm(`Are you sure you want to delete the contest ${this.contest.name}?`)) {
                await OrganizerRequests.deleteContest(this.contest);
                window.location.reload();
            }
        },
        generateTeams(contest: Contest) {
            this.$router.push(`/organizer/other/?contest=${contest.name}`);
        },
        async finishContest(contest: Contest) {
            await OrganizerRequests.sendContestOverNotifications(contest);
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

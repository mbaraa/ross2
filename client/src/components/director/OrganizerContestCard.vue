<template>
    <div class="main bg-blue-darken-4">
        <ContestCard :contest="contest" @click="openContestDetails"/>
        <v-btn @click="deleteContest" icon color="error" class="delete">
            <FontAwesomeIcon class="text-white" :icon="{prefix:'fas', iconName:'trash'}"/>
        </v-btn>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Contest from "@/models/Contest.ts";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faTrash} from "@fortawesome/free-solid-svg-icons"
import {library} from "@fortawesome/fontawesome-svg-core";
import ContestCard from "@/components/contest/ContestCard.vue";
import Organizer from "@/models/Organizer";

library.add(faTrash);

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
                await Organizer.deleteContest(this.contest);
                window.location.reload();
            }
        }
    }
});
</script>
<style scoped>
.main {
    cursor: pointer;
    border-radius: 5px;
}

.delete {
    margin-bottom: 10px;
}
</style>

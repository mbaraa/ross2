<template>
    <div class="grid">
        <DirectorCreateContest/>
        <div v-if="contests.length > 0">
            <div v-for="contest in contests" :key="contest" class="grid">
                <DirectorContestCard :contest="contest"/>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import DirectorContestCard from "@/components/director/OrganizerContestCard.vue";
import DirectorCreateContest from "@/components/director/DirectorCreateContest.vue";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";

export default defineComponent({
    name: 'DirectorContests',
    components: {
        DirectorCreateContest,
        DirectorContestCard,
    },
    data() {
        return {
            contests: []
        }
    },
    async mounted() {
        this.contests = await OrganizerRequests.getContests();
    }
});
</script>

<style scoped>
.grid {
    display: inline-grid;
    padding: 20px;
}
</style>

<template>
    <div class="main" v-if="contests.length > 0">
        <div v-for="contest in contests" :key="contest" class="grid">
            <ContestantContestCard :contest="contest"/>
        </div>
    </div>
    <h1 v-else style="text-align: center">No contests are available at this time!</h1>
</template>

<script lang="ts">
import Contest from "@/models/Contest.ts";
import {defineComponent} from "vue";
import ContestantContestCard from "@/components/contestant/ContestantContestCard.vue";

export default defineComponent({
    name: 'Contests',
    components: {
        ContestantContestCard,
    },
    data() {
        return {
            contests: []
        }
    },
    async mounted() {
        this.contests = await Contest.getContestsFromServer();
    }
});
</script>

<style scoped>
.main {
    color: white;
    text-align: center;
    margin: 10px auto;

}

.grid {
    display: inline-grid;
    padding: 20px;
}

</style>

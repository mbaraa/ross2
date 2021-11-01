<template>
    <div class="main" v-if="contests.length > 0">
        <div v-for="contest in contests" :key="contest" class="grid">
            <ContestCard :contest="contest"/>
        </div>
    </div>
    <h1 v-else style="text-align: center">No contests are available at this time!</h1>
</template>

<script lang="ts">
import Contest from "@/models/Contest.ts";
import { defineComponent } from "vue";
import ContestCard from "@/components/ContestCard.vue";

export default defineComponent({
    name: 'Contests',
    components: {
        ContestCard
    },
    data() {
        return {
            contests: []
        }
    },
    async mounted() {
        this.contests = await Contest.getContestsFromServer();
        this.contests.push(...await Contest.getContestsFromServer())
        this.contests.push(...await Contest.getContestsFromServer())
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

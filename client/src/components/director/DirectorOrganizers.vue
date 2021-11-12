<template>
    <div class="grid">
        <DirectorCreateContest/>
        <div v-if="organizers.length > 0">
            <div v-for="org in organizers" :key="org" class="grid">
                <DirectorOrganizerCard :organizer="org"/>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import DirectorCreateContest from "@/components/director/DirectorCreateContest.vue";
import Organizer from "@/models/Organizer";
import DirectorOrganizerCard from "@/components/director/DirectorOrganizerCard.vue";

export default defineComponent({
    name: 'DirectorContests',
    components: {
        DirectorOrganizerCard,
        DirectorCreateContest,
    },
    data() {
        return {
            organizers: []
        }
    },
    async mounted() {
        this.organizers = await Organizer.getSubOrganizers();
    }
});
</script>

<style scoped>
.grid {
    display: inline-grid;
    padding: 20px;
}
</style>

<template>
    <div class="grid">
        <DirectorCreateOrganizer/>
        <div v-if="organizers.length > 0">
            <div v-for="org in organizers" :key="org" class="grid">
                <DirectorOrganizerCard :organizer="org"/>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Organizer from "@/models/Organizer";
import DirectorOrganizerCard from "@/components/director/DirectorOrganizerCard.vue";
import DirectorCreateOrganizer from "@/components/director/DirectorCreateOrganizer.vue";

export default defineComponent({
    name: 'DirectorContests',
    components: {
        DirectorCreateOrganizer,
        DirectorOrganizerCard,
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

<template>
    <div class="grid">
        <AdminCreateDirector/>
        <div v-if="directors != null && directors.length > 0">
            <div v-for="dir in directors" :key="dir" class="grid">
                <AdminDirectorCard :director="dir"/>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import AdminDirectorCard from "@/components/admin/AdminDirectorCard.vue";
import AdminRequests from "@/utils/requests/AdminRequests";
import AdminCreateDirector from "@/components/admin/AdminCreateDirector.vue";

export default defineComponent({
    name: 'AdminDirectors',
    components: {
        AdminCreateDirector,
        AdminDirectorCard,
    },
    data() {
        return {
            directors: []
        }
    },
    async mounted() {
        this.directors = await AdminRequests.getDirectors();
    }
});
</script>

<style scoped>
.grid {
    display: inline-grid;
    padding: 20px;
}
</style>

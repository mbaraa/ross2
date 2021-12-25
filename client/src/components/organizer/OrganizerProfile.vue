<template>
    <div v-if="organizerProfile !== null">
        <span class="organizerName" style="font-size: 1.5em"><b>
            Organizer Roles: {{
                organizerProfile.roles_names.join(", ")
            }}</b></span>
        <br/>

        <DirectorOperations :director="organizerProfile"/>
        <v-divider/>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import DirectorOperations from "@/components/director/DirectorOperations.vue";
import Organizer from "@/models/Organizer";

export default defineComponent({
    name: "OrganizerProfile",
    components: {
        DirectorOperations,
    },
    props: {
        organizerProfile: Organizer,
    },
    data() {
        return {}
    },
    methods: {}
});
</script>

<style scoped>

</style>

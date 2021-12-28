<template>
    <div v-if="organizerProfile !== null">
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
    mounted() {
        if (this.organizerProfile.user.id > 0) {
            this.$router.push("/profile/contests/");
        }
    },
    data() {
        return {}
    },
    methods: {}
});
</script>

<style scoped>

</style>

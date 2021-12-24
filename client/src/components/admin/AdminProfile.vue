<template>
    <div v-if="adminProfile !== null">
        <br/>

        <div v-if="(adminProfile.user.user_type_base & 16) !== 0">
            <AdminDirectors/>
        </div>
        <v-divider/>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Admin from "@/models/Admin";
import AdminDirectors from "@/components/admin/AdminDirectors.vue";

export default defineComponent({
    name: "AdminProfile",
    components: {AdminDirectors},
    props: {
        adminProfile: Admin,
    },
});
</script>

<style scoped>

</style>

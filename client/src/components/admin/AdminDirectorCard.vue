<template>
    <div class="main bg-red-darken-3">
        <OrganizerCard :organizer="director"/>
        <v-btn @click="deleteOrganizer()" icon color="error" class="delete">
            <FontAwesomeIcon class="text-white" :icon="{prefix:'fas', iconName:'trash'}"/>
        </v-btn>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faTrash} from "@fortawesome/free-solid-svg-icons"
import {library} from "@fortawesome/fontawesome-svg-core";
import OrganizerCard from "@/components/organizer/OrganizerCard.vue";
import Organizer from "../../models/Organizer";
import AdminRequests from "@/utils/requests/AdminRequests";

library.add(faTrash)

export default defineComponent({
    name: "AdminDirectorCard",
    components: {
        OrganizerCard,
        FontAwesomeIcon
    },
    props: {
        director: Organizer
    },
    methods: {
        async deleteOrganizer() {
            if (window.confirm("are you sure you want to delete this director?")) {
                await AdminRequests.deleteDirector(this.director);
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

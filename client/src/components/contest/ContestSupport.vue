<template>
    <div class="main" v-if="organizers.length > 0">
        <h2 class="text-black">
            <FontAwesomeIcon :icon="{prefix:'fas', iconName:'file-alt'}"/>&nbsp;Contact us for support 😁</h2>
        <div v-for="org in organizers" :key="org" class="grid">
            <OrganizerCard :organizer="org"/>
        </div>
    </div>
    <h1 v-else>Something went wrong 😥</h1>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Contest from "@/models/Contest";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faFileAlt} from "@fortawesome/free-solid-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";
import OrganizerCard from "@/components/organizer/OrganizerCard.vue";

library.add(faFileAlt);

export default defineComponent({
    name: "ContestSupport",
    components: {
        OrganizerCard,
        FontAwesomeIcon
    },
    data() {
        return {
            organizers: []
        }
    },
    async mounted() {
        const contest = await Contest.getContestFromServer(this.$route.query.id);
        this.organizers = contest.organizers;
    },
    methods: {}
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

<template>
    <div class="main bg-grey-darken-3" @click="openContestDetails">
        <br/>
        <img class="contestLogo" :alt="contest.name + ' logo'" :src="contest.logo_path"/>
        <h1>{{ contest.name }}</h1>
        <v-divider/>
        <p>{{ contest.description }}</p>
        <v-divider/>

        <p>
            <FontAwesomeIcon :icon="{prefix: 'fas', iconName: 'clock'}"/>&nbsp;
            {{ getLocaleTime(contest.starts_at) }}
        </p>
        <p>
            <FontAwesomeIcon :icon="{prefix: 'fas', iconName: 'map-marker-alt'}"/>&nbsp;{{ contest.location }}
        </p>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Contest from "@/models/Contest.ts";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faClock, faMapMarkerAlt} from "@fortawesome/free-solid-svg-icons"
import {library} from "@fortawesome/fontawesome-svg-core";
import {getLocaleTime} from "@/utils";

library.add(faClock, faMapMarkerAlt);

export default defineComponent({
    name: "ContestCard",
    props: {
        contest: Contest
    },
    components: {
        FontAwesomeIcon
    },
    methods: {
        openContestDetails() {
            this.$router.push(`/contest/details/?id=${this.contest.id}`);
        },
        getLocaleTime(time: Date): string {
            return getLocaleTime(time);
        }
    }
});
</script>

<style scoped>
.main {
    color: white;
    text-align: center;
    width: 350px;
    margin: 0 auto;
    height: auto;
    border-radius: 5px;
    padding: 5px;
    cursor: pointer;
}

.contestLogo {
    width: 75px;
    height: 75px;
    border-radius: 100%;
    background-color: white;
    padding: 5px;
}

h1 {
    font-size: 2em;
}

p {
    font-size: 1.3em;
}
</style>

<template>
    <div class="main bg-blue-darken-4" @click="openContestDetails">
        <br/>
        <img class="contestLogo" :alt="contest.name + ' logo'" :src="config.backendAddress+contest.logo_path"/>
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
import config from "@/config";

library.add(faClock, faMapMarkerAlt);

export default defineComponent({
    name: "ContestCard",
    props: {
        contest: Contest
    },
    components: {
        FontAwesomeIcon
    },
    data() {
        return {
            config: config,
        }
    },
    methods: {
        openContestDetails() {
            this.$router.push(`/contest/details/?id=${this.contest.id}`);
        },
        getLocaleTime(ts: number): string {
            return getLocaleTime(new Date(ts));
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
    width: 100px;
    height: 100px;
    border-radius: 100%;
}

h1 {
    font-size: 2em;
}

p {
    font-size: 1.3em;
}
</style>

<template>
    <div class="main">
        <div class="basic">
            <table class="details" style="width: 100%">
                <tr>
                    <td v-for="field in fields" :key="field" style="width: 100%">
                        <FontAwesomeIcon :icon="field.icon"/>&nbsp;<b>{{ field.name }}</b>
                        <br/>
                        <div v-if="field.name !== 'Registration ends in'">{{ field.value }}</div>
                        <div v-else><TimerCountdown :end-timestamp="field.value"/></div>
                    </td>
                </tr>
            </table>
        </div>
        <v-divider/>
        <div class="desc">
            <FontAwesomeIcon :icon="{prefix: 'fas', iconName: 'file-alt'}"/>&nbsp;<b>Description:</b><br/>
            <span>{{ contest.description }}</span>
        </div>
<!--        <v-divider/>-->
<!--        <div class="desc">-->
<!--            <FontAwesomeIcon :icon="{prefix: 'fas', iconName: 'file-alt'}"/>&nbsp;<b>Participation conditions:</b><br/>-->
<!--            <span>The contestant must be from one of these majors:</span>-->
<!--            <ul class="majors">-->
<!--                <li v-for="major in majors" :key="major">-->
<!--                    {{ major }}-->
<!--                </li>-->
<!--            </ul>-->
<!--        </div>-->
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Contest from "@/models/Contest";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faClock, faFileAlt, faMapMarkerAlt, faUsers} from "@fortawesome/free-solid-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";
import {formatDuration, getLocaleTime} from "@/utils";
import TimerCountdown from "@/components/TimerCountdown.vue";

library.add(faClock, faUsers, faMapMarkerAlt, faFileAlt);

export default defineComponent({
    name: "ContestDetails",
    components: {
        FontAwesomeIcon,
        TimerCountdown,
    },
    data() {
        return {
            contest: {},
            fields: [],
            majors: []
        }
    },
    async mounted() {
        this.contest = await Contest.getContestFromServer(this.$route.query.id);
        this.fields = [
            {name: "Location", icon: {prefix: "fas", iconName: "map-marker-alt"}, value: this.contest.location},
            {
                name: "Members limit",
                icon: {prefix: "fas", iconName: "users"},
                value: `Min: ${this.contest.participation_conditions.min_team_members}, Max: ${this.contest.participation_conditions.max_team_members}`
            },
            {
                name: "Starts at",
                icon: {prefix: "fas", iconName: "clock"},
                value: getLocaleTime(this.contest.starts_at)
            },
            {
                name: "Registration ends in",
                icon: {prefix: "fas", iconName: "clock"},
                value: this.contest.registration_ends
            },
            {
                name: "Duration",
                icon: {prefix: "fas", iconName: "clock"},
                value: formatDuration(this.contest.duration)
            },
        ];
        this.majors = this.contest.participation_conditions.majors == 14 ? // yeah it's kinda stupid using a literal, but it's better than "indexOf()"
            ["Any"] :
            this.contest.participation_conditions.majors_names;
    },
});
</script>

<style scoped>
.main {
    text-align: center;
    width: 100%;
    display: inline-grid;

    background-color: #B0BEC5;
    border-radius: 10px;
}

.basic {
    display: inline-grid;
    width: 100%;
}

table {
    table-layout: fixed;
    overflow-x: auto;
}

td {
    padding: 10px;

}

@media only screen and (max-width: 500px) {
    table {
        overflow-x: scroll;
        display: block;
    }
}

.desc {
    text-align: left;
    padding: 20px;
}

.majors {
    padding-left: 20px;
}

</style>

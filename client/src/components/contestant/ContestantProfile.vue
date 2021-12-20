<template>
    <div v-if="contestantProfile !== null">
            <span class="contestName" style="font-size: 1.5em">
                    <b>University ID: {{ contestantProfile.university_id }}</b>
                </span>
        <br/>

        <div v-if="checkTeam()">
            <v-divider/>

            <h3>
                <FontAwesomeIcon :icon="{ prefix: 'fas', iconName: 'file-alt' }"/>&nbsp;Team details:
            </h3>

            <ul v-if="team !== null">
                <li>Team name: {{ team.name }}</li>
                <li
                    title="share this id with team members you want to join this team"
                >Team ID: {{ team.id }}
                </li>
                <li>
                    Team members:
                    <ul v-for="member in team.members" :key="member">
                        <li>{{ member.name }}</li>
                    </ul>
                </li>
            </ul>

            <div class="buttons">
                <v-btn @click="leaveTeam" class="text-blue-darken-4">Leave team</v-btn>
                <v-btn
                    v-if="checkLeader()"
                    @click="deleteTeam"
                    class="text-red-darken-4"
                >Delete team
                </v-btn>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Contestant from "@/models/Contestant";
import ContestantRequests from "@/utils/requests/ContestantRequests";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faFileAlt} from "@fortawesome/free-solid-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";

library.add(faFileAlt);

export default defineComponent({
    name: "ContestantProfile",
    components: {
        FontAwesomeIcon,
    },
    props: {
        contestantProfile: Contestant,
    },
    data() {
        return {
            team: null,
        }
    },
    methods: {
        async leaveTeam() {
            if (window.confirm("Are you sure you want to leave your team?")) {
                await ContestantRequests.leaveTeam();
                window.location.reload();
            }
        },
        async deleteTeam() {
            if (window.confirm("Are you sure you want to delete your team :)")) {
                if (this.team == null || this.team.name.length == 0) {
                    window.alert("woah, something went wrong :(");
                    return;
                }
                await ContestantRequests.deleteTeam(this.team);
                window.location.reload();
            }
        },
        checkLeader(): boolean {
            return this.profile.id == this.team.leader_id;
        },
        checkTeam(): boolean {
            return this.team != null && this.team.id > 1;
        },

    }
});
</script>

<style scoped>

</style>

<template>
    <div v-if="contestantProfile !== null">
            <span class="contestName" style="font-size: 1.5em">
                    <b>University ID: {{ contestantProfile.university_id }}</b>
                </span>
        <br/>

        <div v-if="team !== null">
            <v-divider/>

            <div v-if="checkTeam()">
                <v-table style="width: 600px; margin: 0 auto; border-radius: 10px">
                    <caption>
                        <FontAwesomeIcon :icon="{ prefix: 'fas', iconName: 'file-alt' }"/>&nbsp;Team details:
                    </caption>
                    <tbody>
                    <tr>
                        <td>Team name</td>
                        <td>{{ team.name }}</td>
                    </tr>
                    <tr title="share this id with team members you want to join this team">
                        <td>Team ID</td>
                        <td>{{ team.join_id }}</td>
                    </tr>
                    <tr v-if="team.members.length > 0">
                        <td>Team members</td>
                        <td>{{ getMembersNames() }}</td>
                    </tr>
                    </tbody>
                </v-table>
                <div class="buttons">
                    <v-btn @click="leaveTeam" class="text-blue-darken-4">Leave team</v-btn>
                    &nbsp;
                    <v-btn
                        v-if="checkLeader()"
                        @click="deleteTeam"
                        class="text-red-darken-4"
                    >Delete team
                    </v-btn>
                </div>
            </div>

        </div>
        <v-divider/>
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
    async mounted() {
        this.team = await ContestantRequests.getTeam();
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
            return this.team != null && this.contestantProfile.user.id == this.team.leader_id;
        },
        checkTeam(): boolean {
            return this.team != null && this.team.id > 1;
        },
        getMembersNames(): string {
            let names = new Array<string>();
            this.team.members.forEach((c: Contestant) => {
                names.push(<string>c.user.name);
            });
            return names.join(", ");
        }
    }
});
</script>

<style scoped>

</style>

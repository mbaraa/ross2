<template>
    <br/>
    <div v-if="teams != null && teams.length > 0">
        <div class="teams" v-for="team in teams" :key="team">
            <DirectorTeamCard :team="team"/>
        </div>

        <br/>
        <v-btn class="bg-red-darken-4 text-white" @click="saveTeams()">Update teams</v-btn>

        <br/>
        <div v-if="removedMembers.length > 0">
            <h1 class="text-red">Removed members:</h1>
            <table class="removed">
                <tr class="bg-green">
                    <th>Name</th>
                    <th>ID</th>
                </tr>
                <tr :class="getContClass(member)" v-for="member in this.removedMembers" :key="member">
                    <td> {{ member.name }}</td>
                    <td> {{ member.id }}</td>
                </tr>
            </table>
        </div>
    </div>
    <h1 v-else>No teams are registered for this contest so far!</h1>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import DirectorTeamCard from "@/components/director/DirectorTeamCard.vue";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";
import Contestant from "@/models/Contestant";

export default defineComponent({
    name: "DirectorTeams",
    components: {DirectorTeamCard},
    data() {
        return {
            teams: [],
            removedMembers: this.$store.getters.getRemovedContestants,
        }
    },
    async mounted() {
        this.teams = await (await OrganizerRequests.getContest(+this.$route.query["contest"])).teams;
    },
    methods: {
        openContests() {
            this.$router.push(`/organizer/contests/?id=${this.$route.query["id"]}`)
        },
        async saveTeams() {
            console.log("lol", this.$store.getters.getModifiedTeams)
            if (this.$store.getters.getModifiedTeams.length == 0) {
                window.alert("no teams were modified!");
                return;
            }
            if (window.confirm("you are about to update some teams, continue?")) {
                await OrganizerRequests.updateTeams(
                    this.$store.getters.getModifiedTeams, this.$store.getters.getRemovedContestants);
            }
        },
        getContClass(contestant: Contestant): string {
            return this.removedMembers.indexOf(contestant) % 2 == 0 ? "cont1" : "cont2";
        }
    }
});
</script>

<style scoped>
.cont1 {
    background-color: #C8E6C9;
}

.cont2 {
    background-color: #A5D6A7;
}

tr, td {
    padding: 3px;
}

.removed {
    margin: 0 auto;
    width: 100%;
}

.teams {
    margin: 10px;
    display: inline-grid;
}
</style>

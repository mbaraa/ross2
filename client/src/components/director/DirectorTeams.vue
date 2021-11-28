<template>
    <br/>
    <div v-if="teams != null && teams.length > 0">
        <div class="grid">
            <div class="teams" v-for="team in teams" :key="team">
                <DirectorTeamCard :team="team"/>
            </div>

            <!-- it's fucked up, I know -->
            <v-dialog
                max-height="400"
                max-width="-40"
                transition="dialog-bottom-transition"
                v-model="dialog"
                scrollable>
                <template v-slot:activator="{on, attrs}">
                    <div v-bind="attrs" v-on="on" @click="openDialog()" style="display: inline; color: white">
                        <div v-bind="attrs" v-on="on" class="teamForm bg-grey-darken-3" @click="dialog = true"
                             title="just click it!">
                            <h1>Create Team</h1>
                            <v-divider/>
                            <FontAwesomeIcon style="font-size: 3em" :icon="{prefix:'fas', iconName:'plus'}"/>
                        </div>
                    </div>
                </template>

                <v-card elevation="16" class="teamForm">
                    <v-card-title>
                        <span class="text-h4">Create Team</span>
                    </v-card-title>

                    <v-text-field label="Team name" v-model="newTeam.name" autofocus/>

                    <v-btn class="bg-red" @click="dialog = false">
                        Close
                    </v-btn>&nbsp;
                    <v-btn class="bg-blue" @click="createTeam()">
                        Create
                    </v-btn>
                </v-card>
            </v-dialog>
            <!-- end of fuck-up :) -->

        </div>

        <div>
            <v-btn class="bg-red-darken-4 text-white" @click="saveTeams()">Update teams</v-btn>

            <br/>
            <div v-if="removedMembers.length > 0">
                <h1 class="text-red">Available contestants:</h1>
                <table class="removed">
                    <tr class="bg-green">
                        <th title="use it to add a contestant to a specific team">ID</th>
                        <th>Contestant name</th>
                        <th>University ID</th>
                        <th>Gender</th>
                        <th>Can participate with the other gender</th>
                    </tr>
                    <tr v-for="cont in removedMembers" :key="cont" :class="getContClass(cont)">
                        <td title="use it to add this contestant to a specific team">{{ cont.id }}</td>
                        <td>{{ cont.name }}</td>
                        <td>{{ cont.university_id }}</td>
                        <td>{{ cont.gender ? "Male" : "Female" }}</td>
                        <td>{{ cont.participate_with_other ? "Yes" : "No" }}</td>
                    </tr>
                </table>
            </div>
        </div>
    </div>
    <h1 v-else>No teams are registered for this contest so far!</h1>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import DirectorTeamCard from "@/components/director/DirectorTeamCard.vue";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";
import Contestant from "@/models/Contestant";
import Team from "@/models/Team";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {library} from "@fortawesome/fontawesome-svg-core";
import {faUserPlus} from "@fortawesome/free-solid-svg-icons";
import Contest from "@/models/Contest";


library.add(faUserPlus);

export default defineComponent({
    name: "DirectorTeams",
    components: {FontAwesomeIcon, DirectorTeamCard},
    data() {
        return {
            contest: new Contest(),
            teams: [],
            removedMembers: this.$store.getters.getRemovedContestants,
            newTeam: new Team(),
            dialog: false
        }
    },
    async mounted() {
        this.contest = await OrganizerRequests.getContest(+this.$route.query["contest"])
        this.teams = await this.contest.teams;
    },
    methods: {
        openContests() {
            this.$router.push(`/organizer/contests/?id=${this.$route.query["id"]}`)
        },
        async saveTeams() {
            if (this.$store.getters.getModifiedTeams.length == 0) {
                window.alert("no teams were modified!");
                return;
            }
            if (window.confirm("you are about to update some teams, continue?")) {
                await OrganizerRequests.updateTeams(
                    this.$store.getters.getModifiedTeams, this.$store.getters.getRemovedContestants);
                localStorage.removeItem("modifiedTeams");
            }
        },
        createTeam() {
            this.newTeam.id = Number.parseInt(`${Math.sqrt((new Date()).getTime() + 1)}`);

            this.newTeam.contests.push(<Contest>{id: this.contest.id});

            this.$store.dispatch("addTeam", this.newTeam);
            this.teams.push(this.newTeam);

            this.newTeam = new Team();
            this.dialog = false;
        },
        openDialog() {
            this.dialog = true;
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

.grid {
    display: inline-grid;
    padding: 20px;
}

.teamForm {
    padding: 10px;
    margin: 0 auto;
    width: 350px;
    overflow-y: auto;
    border-radius: 5px;
}
</style>

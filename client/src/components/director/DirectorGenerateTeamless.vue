<template>
    <div v-if="contests.length > 0">
        <h1>Generate teams for teamless contestants:</h1>
        <v-divider/>
        <table class="opts">
            <tr>
                <td>
                    <label for="contest">Select contest: </label>
                </td>
                <td>
                    <select style="background-color: #eeeeee" id="contest" v-model="selection">
                        <option v-for="contest in contests" :key="contest">
                            {{ contest.name }}
                        </option>
                    </select>
                </td>
            </tr>
            <tr>
                <td>
                    <label for="genType">Select generation type: </label>
                </td>
                <td>
                    <select style="background-color: #eeeeee" id="contest" v-model="genType">
                        <option value="numbered">Numbered</option>
                        <option value="random">Random teams names</option>
                    </select>
                </td>
            </tr>
            <tr>
                <td colspan="2">
                    <v-btn @click="generateTeams()">Generate Teams</v-btn>
                </td>
            </tr>
        </table>

        <!-- obaa -->
        <div v-if="generatedTeams.length > 0">
            <div class="teams" v-for="team in generatedTeams" :key="team">
                <TeamCard :team="team"/>
            </div>
            <br/>
            <v-btn class="bg-red-darken-4 text-white" @click="saveTeams()">Save teams</v-btn>
        </div>

        <!-- hmm -->
        <div v-if="leftTeamless.length > 0">
            <br/>
            <h2>Contestants left with no teams:</h2>
            <table class="tls">
                <tr>
                    <th>Contestant name</th>
                    <th>University ID</th>
                    <th>Gender</th>
                    <th>Can participate with the other gender</th>
                </tr>
                <tr v-for="cont in leftTeamless" :key="cont" :class="getContClass(cont)">
                    <td>{{ cont.name }}</td>
                    <td>{{ cont.university_id }}</td>
                    <td>{{ cont.gender? "Male": "Female"}}</td>
                    <td>{{ cont.participate_with_other? "Yes": "No"}}</td>
                </tr>
            </table>
        </div>

    </div>
    <div v-else>
        <h2>no contests were found 🙂</h2>
        <h3 @click="openContests()">you may want to create some!</h3>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Organizer from "@/models/Organizer";
import Contest from "@/models/Contest";
import TeamCard from "@/components/team/TeamCard.vue";
import Contestant from "@/models/Contestant";

export default defineComponent({
    name: "DirectorGenerateTeamless",
    components: {TeamCard},
    data() {
        return {
            contests: [],
            selection: "",
            generated: false,
            genType: "random",
            generatedTeams: [],
            leftTeamless: [],
        }
    },
    async mounted() {
        this.contests = await Organizer.getContests();
        this.selection = this.contests[0].name;
    },
    methods: {
        openContests() {
            this.$router.push(`/organizer/contests/?id=${this.$route.query["id"]}`)
        },
        selectContest(): Contest | null {
            for (const contest of this.contests) {
                if (contest.name == this.selection) {
                    return contest;
                }
            }
            return null
        },
        async generateTeams() {
            [this.generatedTeams, this.leftTeamless] = await Organizer.generateTeams(this.selectContest(), this.genType);

            this.generated = true;
        },
        async saveTeams() {
            await Organizer.saveTeams(this.generatedTeams);
        },
        getContClass(cont: Contestant): string {
            return this.leftTeamless.indexOf(cont) % 2 == 0 ? "cont1" : "cont2";
        }
    }
});
</script>

<style scoped>
.opts {
    margin: 0 auto;
}

tr, td {
    padding: 3px;
}

.teams {
    margin: 10px;
    display: inline-grid;
}

.cont1 {
    background-color: #d0d0d0;
}

.cont2 {
    background-color: #a0a0a0;
}

.tls {
    margin: 0 auto;
    width: 100%;
}
</style>

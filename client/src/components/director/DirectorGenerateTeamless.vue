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
            this.generatedTeams =
                await Organizer.generateTeams(this.selectContest(), this.genType);

            this.generated = true;
        },
        async saveTeams() {
            await Organizer.saveTeams(this.generatedTeams);
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
</style>

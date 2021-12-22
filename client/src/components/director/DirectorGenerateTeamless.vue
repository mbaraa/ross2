<template>
    <div v-if="contests.length > 0">
        <h1>Generate teams for teamless contestants:</h1>
        <v-divider/>
        <table class="opts">
            <tr>
                <td>
                    <label for="contest">Select contest:</label>
                </td>
                <td>
                    <select style="background-color: #eeeeee" id="contest" v-model="selection">
                        <option v-for="contest in contests" :key="contest">{{ contest.name }}</option>
                    </select>
                </td>
            </tr>
            <tr>
                <td>
                    <label for="genType">Select generation type:</label>
                </td>
                <td>
                    <select style="background-color: #eeeeee" id="contest" v-model="genType">
                        <option value="ordered" @click="closeNamesFileUpload()">Numbered</option>
                        <option value="random" @click="closeNamesFileUpload()">Random teams names</option>
                        <option value="given" @click="openNamesFileUpload()">Given teams names</option>
                    </select>
                </td>
            </tr>
            <tr v-if="!hideNamesFileUpload">
                <td colspan="2" style="color: #212121">
                    upload a file with the wanted teams names
                    <b>separated by a comma(,)</b>
                    <br/>eg: name1,name2,name3...
                </td>
            </tr>
            <tr v-if="!hideNamesFileUpload">
                <td colspan="2">
                    <v-file-input show-size label="Names file" prepend-icon @change="selectFile"/>
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
                <DirectorTeamCard :team="team"/>
            </div>
            <br/>
            <v-btn class="bg-red-darken-4 text-white" @click="saveTeams()">Save teams</v-btn>
        </div>

        <!-- hmm -->
        <div v-if="leftTeamless != null && leftTeamless.length > 0">
            <br/>
            <h2>Contestants left with no teams:</h2>
            <table class="tls">
                <tr class="bg-green">
                    <th title="use it to add a contestant to a specific team">ID</th>
                    <th>Contestant name</th>
                    <th>University ID</th>
                    <th>Gender</th>
                    <th>Can participate with the other gender</th>
                </tr>
                <tr v-for="cont in leftTeamless" :key="cont" :class="getContClass(cont)">
                    <td title="use it to add this contestant to a specific team">{{ cont.user.id }}</td>
                    <td>{{ cont.user.name }}</td>
                    <td>{{ cont.university_id }}</td>
                    <td>{{ cont.gender ? "Male" : "Female" }}</td>
                    <td>{{ cont.participate_with_other ? "Yes" : "No" }}</td>
                </tr>
            </table>
        </div>

        <h1
            class="text-red"
            v-if="noTeamless"
        >no teamless contestants were found for this contest :)</h1>
    </div>
    <div v-else>
        <h2>no contests were found 🙂</h2>
        <h3 @click="openContests()">you may want to create some!</h3>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Contest from "@/models/Contest";
import Contestant from "@/models/Contestant";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";
import DirectorTeamCard from "@/components/director/DirectorTeamCard.vue";

export default defineComponent({
    name: "DirectorGenerateTeamless",
    components: {DirectorTeamCard},
    data() {
        return {
            contests: [],
            selection: "",
            generated: false,
            genType: "random",
            generatedTeams: [],
            leftTeamless: [],
            noTeamless: false,
            hideNamesFileUpload: true,
            namesFile: undefined,
        }
    },
    async mounted() {
        this.contests = await OrganizerRequests.getContests();
        const name = this.$route.query["contest"];
        this.selection = name ?? this.contests[0].name;
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
        checkNamesFile(): boolean {
            return (this.genType != "given") || (this.genType == "given" && this.namesFile !== undefined);
        },
        async readNamesFile(): Promise<string[]> {
            return (await this.namesFile.text()).replace("\n", "").split(",");
        },
        async generateTeams() {
            if (!this.checkNamesFile()) { // had to do it this way :(
                window.alert("select file to upload!");
                return;
            }
            [this.generatedTeams, this.leftTeamless] =
                await OrganizerRequests.generateTeams(this.selectContest(), this.genType, this.namesFile !== undefined ? (await this.readNamesFile()) : []);

            if (this.generatedTeams.length == 0 && this.leftTeamless == null) {
                this.noTeamless = true;
                return;
            }

            this.generated = true;
            this.noTeamless = false;
            if (this.leftTeamless.length > 0) { // so left contestants can be assigned to any team later :)
                for (const c of this.leftTeamless) {
                    await this.$store.dispatch("addContestantToRemoved", c);
                }
            }
        },
        async saveTeams() {
            if (window.confirm("are you sure of the teams you are about to register?")) {
                const selectedContest = this.selectContest();

                for (const team of this.generatedTeams) {
                    team.contests = [selectedContest];
                }

                await OrganizerRequests.saveTeams(this.generatedTeams);
                window.location.reload();
            }
        },
        openNamesFileUpload() {
            this.hideNamesFileUpload = false;
        },
        closeNamesFileUpload() {
            this.hideNamesFileUpload = true;
        },
        selectFile(file: any) {
            this.namesFile = file.target.files[0];
            if (this.namesFile.type != "text/plain") {
                window.alert("file must be of text type!");
                this.namesFile = undefined;
            }
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

tr,
td {
    padding: 3px;
}

.teams {
    margin: 10px;
    display: inline-grid;
}

.cont1 {
    background-color: #c8e6c9;
}

.cont2 {
    background-color: #a5d6a7;
}

.tls {
    margin: 0 auto;
    width: 100%;
}
</style>

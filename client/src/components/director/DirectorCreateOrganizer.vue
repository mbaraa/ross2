<template>
    <v-dialog
        max-height="400"
        max-width="-40"
        transition="dialog-bottom-transition"
        v-model="dialog"
        scrollable>
        <template v-slot:activator="{on, attrs}">
            <div v-bind="attrs" v-on="on" class="main3 bg-grey-darken-3" @click="dialog = true" title="just click it!">
                <h1>Add Organizer!</h1>
                <v-divider/>
                <FontAwesomeIcon style="font-size: 3em" :icon="{prefix:'fas', iconName:'plus'}"/>
            </div>
        </template>

        <v-card elevation="16" class="contestForm">
            <v-card-title>
                <span class="text-h4">Add Organizer</span>
            </v-card-title>

            <v-text-field label="Name" v-model="newOrganizer.name" autofocus/>
            <v-text-field label="Gmail" v-model="newOrganizer.email"/>

            <div v-if="contests.length > 0">
                <h4>Set contest for organizer</h4>
                <select v-model="selectedContest" style="background-color: #eeeeee">
                    <option v-for="contest in contests" :key="contest">
                        {{ contest.name }}
                    </option>
                </select>
                <h4>Set roles:</h4>
                <div v-for="role in roles" :key="role">
                    <input type="checkbox" :id="role.name" name="role" :value="role.name" :checked="role.checked"
                           v-model="role.checked"/>
                    <label :for="role.name">&nbsp;{{ role.name }}</label>
                </div>
            </div>

            <v-btn class="bg-red" @click="dialog = false">
                Close
            </v-btn>&nbsp;
            <v-btn class="bg-blue" @click="createOrganizer">
                Create
            </v-btn>
        </v-card>
    </v-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {library} from "@fortawesome/fontawesome-svg-core";
import {faPlus} from "@fortawesome/free-solid-svg-icons";
import Organizer from "@/models/Organizer";
import Contest from "@/models/Contest";

library.add(faPlus);

export default defineComponent({
    name: "DirectorCreateOrganizer",
    components: {
        FontAwesomeIcon
    },
    data() {
        return {
            dialog: false,
            newOrganizer: new Organizer(),
            contests: [],
            selectedContest: "",
            roles: [
                {name: "Core Organizer", checked: false},
                {name: "Chief Judge", checked: false},
                {name: "Judge", checked: false},
                {name: "Media", checked: false},
                {name: "Balloons", checked: false},
                {name: "Technical", checked: false},
                {name: "Coordinator", checked: false},
                {name: "Food", checked: false},
                {name: "Receptionist", checked: false},
            ],
        }
    },
    async mounted() {
        this.contests = await Organizer.getContests();
        this.selectedContest = this.contests[0].name;
    },
    methods: {
        async createOrganizer() {
            this.setContest();
            this.setRoles();
            await Organizer.createOrganizer(this.newOrganizer);

            window.alert("organizer was created successfully!");
            window.location.reload();
        },
        setContest() {
            let contest = new Contest();
            for (const c of this.contests) {
                if (c.name == this.selectedContest) {
                    contest = c;
                }
            }
            this.newOrganizer.contests.push(contest);
        },
        setRoles() {
            this.newOrganizer.roles = 0;
            for (let i = 0; i <= 8; i++) {
                if (this.roles[i].checked) {
                    this.newOrganizer.roles |= (1 << (i + 1)); // i+1, because the roles can't start from director but the roles array starts from 0 :)
                    console.log("roles: ", (1 << (i + 1)));
                }
            }
        }
    }
});
</script>

<style scoped>
.main3 {
    color: white;
    text-align: center;
    width: 350px;
    margin: 0 auto;
    height: auto;
    border-radius: 5px;
    padding: 5px;
    cursor: pointer;
    font-family: 'Ropa Sans', sans-serif;
}

.contestForm {
    padding: 10px;
    margin: 0 auto;
    width: 400px;
    overflow-y: auto;
}
</style>

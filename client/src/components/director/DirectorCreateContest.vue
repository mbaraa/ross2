<template>
    <v-dialog
        max-height="400"
        max-width="-40"
        transition="dialog-bottom-transition"
        v-model="dialog"
        scrollable>
        <template v-slot:activator="{on, attrs}">
            <div v-bind="attrs" v-on="on" class="main3 bg-grey-darken-3" @click="dialog = true" title="just click it!">
                <h1>Create Contest!</h1>
                <v-divider/>
                <FontAwesomeIcon style="font-size: 3em" :icon="{prefix:'fas', iconName:'plus'}"/>
            </div>
        </template>

        <v-card
            elevation="16" class="contestForm">
            <v-card-title>
                <span class="text-h4">Create Contest</span>
            </v-card-title>

            <div class="list">
                <v-text-field label="Contest Title" v-model="contest.name" autofocus=""/>

                <label for="starts">Starts at:</label>
                <input id="starts" class="starts" type="datetime-local" v-model="startsAt" required/>

                <label for="ends">Registration ends at:</label>
                <input id="ends" class="starts" type="datetime-local" v-model="regEndsAt" required/>

                <v-text-field label="Duration (in minutes)" v-model="contest.duration" required/>
                <v-text-field label="Location" v-model="contest.location" required/>

                <label style="font-size: 1.2em">Show registered teams for everyone? </label>
                <br/>
                <input type="radio" id="yes" value="true" v-model="contest.teams_hidden">
                <label for="yes"> Yes</label>
                &nbsp;
                <input type="radio" id="no" value="false" v-model="contest.teams_hidden">
                <label for="no"> No</label>
                <br/><br/>

                <v-file-input
                    show-size
                    label="Logo file"
                    prepend-icon=""
                    @change="selectFile"
                ></v-file-input>
                <v-text-field label="Description" v-model="contest.description" required/>
                <!--            <v-text-field label="Allowed Majors" v-model="contest."/>-->
                <v-text-field label="Minimum team members" v-model="contest.participation_conditions.min_team_members"
                              required/>
                <v-text-field label="Maximum team members" v-model="contest.participation_conditions.max_team_members"
                              required/>

                <v-btn class="bg-red" @click="hideDialog()">
                    Close
                </v-btn>&nbsp;
                <v-btn class="bg-blue" @click="createContest">
                    Create
                </v-btn>
            </div>
        </v-card>
    </v-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {library} from "@fortawesome/fontawesome-svg-core";
import {faPlus} from "@fortawesome/free-solid-svg-icons";
import Contest from "@/models/Contest";
import config from "@/config";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";

library.add(faPlus);

export default defineComponent({
    name: "DirectorCreateContest",
    components: {
        FontAwesomeIcon
    },
    data() {
        return {
            dialog: false,
            contest: new Contest(),
            startsAt: new Date(),
            regEndsAt: new Date(),
            logoFile: undefined,
        }
    },
    methods: {
        checkRegAndStartDate(): boolean {
            return (<number>this.contest.registration_ends < <number>this.contest.starts_at);
        },
        async createContest() {
            this.contest.starts_at = (new Date(this.startsAt)).getTime();
            this.contest.registration_ends = (new Date(this.regEndsAt)).getTime();

            this.contest.duration = +this.contest.duration;
            this.contest.participation_conditions.min_team_members = +this.contest.participation_conditions.min_team_members;
            this.contest.participation_conditions.max_team_members = +this.contest.participation_conditions.max_team_members;
            this.contest.teams_hidden = this.contest.teams_hidden == "true";

            const errMsg = await this.uploadLogo();
            if (errMsg.length > 0) {
                window.alert(errMsg)
                return;
            }

            if (!this.checkRegAndStartDate()) {
                window.alert("woah... start date should be after end of registration date!");
                return;
            }

            this.contest.logo_path = '/' + this.logoFile.name;
            await OrganizerRequests.createContest(this.contest);

            this.dialog = false;
            window.alert("contest was created successfully!");
            window.location.reload();
        },
        selectFile(file: File) {
            this.logoFile = file;
        },
        async uploadLogo(): Promise<string> {
            if (!this.logoFile) {
                return "select logo file to upload :)";
            }

            let res = "";
            const formData = new FormData();
            this.logoFile = this.logoFile.target.files[0];
            formData.append("file", this.logoFile);

            await fetch(`${config.backendAddress}/organizer/upload-contest-logo-file/`, {
                method: "POST",
                mode: "cors",
                headers: {
                    // "Content-Type": `multipart/form-data", content type and boundary is calculated by the browser
                    "Authorization": <string>localStorage.getItem("token"),
                },
                body: formData,
            })
                .then(resp => {
                    res = <string><unknown>resp.text()
                })
                .catch(err => {
                    res = err.message;
                })

            return res;
        },
        hideDialog() {
            this.dialog = false;
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
    overflow-y: scroll;
}

.starts {
    color: #606060;
    padding-bottom: 20px;
    padding-left: 10px;
    margin-bottom: 35px;
    padding-top: 15px;

    background-color: #f0f0f0;
    width: 100%;

    border-radius: 5px 5px 0 0;

    border-bottom: #a0a0a0 solid 1px;
}

.list {
    overflow: hidden;
    overflow-y: scroll;
    height: 75vh;
}
</style>

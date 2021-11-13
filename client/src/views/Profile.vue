<template>
    <div class="main" v-if="profile != null">
        <img class="contestantLogo" alt="contestant picture" :src="profile.avatar_url"/>
        <br/>
        <span class="contestName"><b>{{ profile.name }}</b></span>
        <br/>
        <span class="contestName" style="font-size: 1.5em"><b>{{ profile.university_id }}</b></span>
        <br/><br/>
        <v-divider/>
        <div class="buttons">
            <v-btn @click="logout" class="text-red-darken-4">Logout</v-btn>&nbsp;
            <v-btn @click="deleteAccount" class="text-red-darken-4">Delete account</v-btn>
        </div>
        <div v-if="team.name.length > 0">
            <v-divider/>
            <FontAwesomeIcon :icon="{prefix:'fas', iconName:'file-alt'}"/>
            <b>Team details:</b>
            <br/>
            <ul>
                <li>Team name: {{ team.name }}</li>
                <li>Team members:
                    <ul v-for="member in team.members" :key="member">
                        <li>{{ member.name }}</li>
                    </ul>
                </li>
            </ul>
            <div class="buttons">
                <v-btn @click="leaveTeam" class="text-blue-darken-4">Leave team</v-btn>
                <v-btn @click="deleteTeam" class="text-red-darken-4">Delete team</v-btn>
            </div>
        </div>
    </div>
    <div v-else style="padding-top: 20px; text-align: center;">
        <h1 style="font-size: 3em">Oops! you're not logged in</h1>
        <v-btn @click="login" class="bg-red" style="font-size: 2em; padding: 20px">
            <FontAwesomeIcon :icon="{prefix:'fab', iconName:'google'}"/>&nbsp;Login with Google
        </v-btn>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faFileAlt} from "@fortawesome/free-solid-svg-icons";
import {faGoogle} from "@fortawesome/free-brands-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";
import Contestant from "@/models/Contestant";

library.add(faGoogle, faFileAlt);

export default defineComponent({
    name: "Profile",
    components: {
        FontAwesomeIcon
    },
    data() {
        return {
            profile: null,
            team: null,
        }
    },
    async mounted() {
        this.profile = await this.tokenLogin();
        this.team = this.profile.team;
    },
    methods: {
        async login() {
            await Contestant.googleLogin(
                (await this.$gapi.login()).currentUser
            );
            window.location.reload();
        },
        async logout() {
            await this.$gapi.logout();
            await Contestant.logout();
            window.location.reload();
        },
        async deleteAccount() {
            if (window.confirm("Are you sure you want to delete your account?")) {
                await this.$gapi.logout();
                await Contestant.deleteUser();
                window.location.reload();
            }
        },
        async leaveTeam() {
            if (window.confirm("Are you sure you want to leave your team?")) {
                await Contestant.leaveTeam();
                window.location.reload();
            }
        },
        async deleteTeam() {
            if (window.confirm("Are you sure you want to delete your team :)")) {
                if (this.team == null || this.team.name.length == 0) {
                    window.alert("woah, something went wrong :(");
                    return;
                }
                await Contestant.deleteTeam(this.team);
                window.location.reload();
            }
        },
        async tokenLogin(): Promise<Contestant | null> {
            const contestant = await Contestant.login();
            if ((await contestant) != null && !(await contestant)?.profile_finished) {
                await this.$router.push("/finish-profile/");
            }

            return contestant;
        }
    }

});
</script>

<style scoped>
.main {
    text-align: center;
    width: 80%;
    /*display: inline-grid;*/
    margin: 20px auto;
    border: #212121 solid 1px;
    border-radius: 5px;
}

.contestName {
    font-size: 2em;
    color: #212121;
}

.contestantLogo {
    width: 125px;
    height: 125px;
    border-radius: 100%;
    background-color: white;
    padding: 5px;
}

.buttons {
    padding: 5px;
}
</style>

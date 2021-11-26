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
            <v-btn @click="logout()" class="text-red-darken-4">Logout</v-btn>&nbsp;
            <v-btn @click="deleteAccount" class="text-red-darken-4">Delete account</v-btn>
        </div>
        <div v-if="checkTeam()">
            <v-divider/>

            <h3>
                <FontAwesomeIcon :icon="{prefix:'fas', iconName:'file-alt'}"/>
                &nbsp;Team details:
            </h3>

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
                <v-btn v-if="checkLeader()" @click="deleteTeam" class="text-red-darken-4">Delete team</v-btn>
            </div>
        </div>
    </div>
    <div v-else style="padding-top: 20px; text-align: center;">
        <h1 style="font-size: 3em">Oops! you're not logged in</h1>
        <v-btn @click="loginGoogle()" class="bg-red" style="font-size: 2em; padding: 20px">
            <FontAwesomeIcon :icon="{prefix:'fab', iconName:'google'}"/>&nbsp;Login with Google
        </v-btn>
        <br/><br/>
        <v-btn @click="loginMS()" class="bg-grey" style="font-size: 2em; padding: 20px">
            <FontAwesomeIcon :icon="{prefix:'fab', iconName:'microsoft'}"/>&nbsp;Login with ASU Account
        </v-btn>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faFileAlt} from "@fortawesome/free-solid-svg-icons";
import {faGoogle, faMicrosoft} from "@fortawesome/free-brands-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";
import Contestant from "@/models/Contestant";
import ContestantRequests from "@/utils/requests/ContestantRequests";

library.add(faGoogle, faFileAlt, faMicrosoft);

export default defineComponent({
    name: "Profile",
    components: {
        FontAwesomeIcon
    },
    data() {
        return {
            profile: new Contestant(),
            team: null,
        }
    },
    async mounted() {
        this.profile = await this.tokenLogin();
        this.team = await ContestantRequests.getTeam();
    },
    methods: {
        async loginGoogle() {
            await ContestantRequests.googleLogin(
                (await this.$gapi.login()).currentUser
            );
            window.location.reload();
        },
        async $logout() {
            if (this.profile.email.indexOf("@gmail") > -1) {
                await this.$gapi.logout();
            } else {
                await this.$msal.logoutPopup();
            }
        },
        async logout() {
            await this.$logout();

            await ContestantRequests.logout();
            window.location.reload();
        },
        async loginMS() {
            await this.$msal.loginPopup({
                scopes: ["openid", "profile", "User.Read"]
            })
                .then((resp: any) => ContestantRequests.microsoftLogin(resp));
            window.location.reload();
        },
        async deleteAccount() {
            if (window.confirm("Are you sure you want to delete your account?")) {
                await this.$logout();

                await ContestantRequests.deleteUser();
                window.location.reload();
            }
        },
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
        async tokenLogin(): Promise<Contestant | null> {
            const contestant = await ContestantRequests.login();
            if ((await contestant) != null && !(await contestant)?.profile_finished) {
                await this.$router.push("/finish-profile/");
            }

            return contestant;
        },
        checkTeam(): boolean {
            return this.team != null && this.team.id > 1;
        },
        checkLeader(): boolean {
            return this.profile.id == this.team.leader_id;
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

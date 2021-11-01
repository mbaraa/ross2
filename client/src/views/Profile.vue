<template>
    <div class="main" v-if="profile != null">
        <img class="contestLogo" alt="contestant logo" :src="profile.avatar_url"/>
        <br/>
        <span class="contestName"><b>{{ profile.name }}</b></span>
        <br/><br/>
        <v-divider/>
        <v-btn @click="logout">Logout</v-btn>&nbsp;
        <v-btn @click="deleteAccount">Delete account</v-btn>
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
        </div>
    </div>
    <div v-else style="padding-top: 20px; text-align: center;">
        <v-btn @click="login" class="bg-red" style="font-size: 2em; padding: 20px">
            <FontAwesomeIcon :icon="{prefix:'fab', iconName:'google'}"/>&nbsp;Google Login
        </v-btn>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faFileAlt} from "@fortawesome/free-solid-svg-icons";
import {faGoogle} from "@fortawesome/free-brands-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";
import config from "../config";
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
            const user = (await this.$gapi.login()).currentUser;

            await fetch(`${config.backendAddress}/gauth/login/`, {
                method: "POST",
                mode: "cors",
                headers: {
                    "Authorization": user.Zb.id_token
                },
                body: JSON.stringify({
                    name: user.it.Se,
                    avatar_url: user.it.SJ,
                    email: user.it.Tt,
                })
            })
                .then(resp => resp.json())
                .then(data => {
                    localStorage.setItem("token", <string>data["token"]);
                });
            window.location.reload();
        },
        async logout() {
            await this.$gapi.logout();
            await Contestant.logout();
            window.location.reload();
        },
        async deleteAccount() {
            await this.$gapi.logout();
            await Contestant.deleteUser();
            window.location.reload();
        },
        async leaveTeam() {
            this()
        },
        async tokenLogin(): Promise<Contestant> {
            const contestant = await Contestant.login();
            if (!contestant.profile_finished) {
                await this.$router.push("/finish-profile");
            }

            return contestant
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
}

.contestName {
    font-size: 2em;
    color: #212121;
}

</style>

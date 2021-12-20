<template>
    <div class="main" v-if="profile.id !== 0">
        <br/>
        <img class="contestantLogo" alt="contestant picture" :src="profile.avatar_url"/>
        <br/>
        <span class="contestName">
            <b>{{ profile.name }}</b>
        </span>
        <br/>
        <v-divider/>

        <!-- contestant stuff -->
        <ContestantProfile :contestantProfile="contestantProfile"/>
        <v-divider/>

        <!-- organizer stuff-->
        <OrganizerProfile :organizerProfile="organizerProfile"/>
        <v-divider/>

        <div class="buttons">
            <v-btn
                v-if="contestantProfile === null"
                @click="registerAsContestant()"
                class="text-blue-darken-4"
            >Register as contestant
            </v-btn>&nbsp;
            <v-btn @click="logout()" class="text-red-darken-4">Logout</v-btn>&nbsp;
            <!--            <v-btn @click="deleteAccount" class="text-red-darken-4">Delete account</v-btn>-->
        </div>
    </div>
    <div v-else style="padding-top: 20px; text-align: center;">
        <h1 style="font-size: 3em">Oops! you're not logged in</h1>
        <v-btn @click="loginGoogle()" class="bg-red" style="font-size: 2em; padding: 20px">
            <FontAwesomeIcon :icon="{ prefix: 'fab', iconName: 'google' }"/>&nbsp;Login with Google
        </v-btn>
        <br/><br/>
        <v-btn @click="loginMS()" class="bg-grey" style="font-size: 2em; padding: 20px">
            <FontAwesomeIcon :icon="{ prefix: 'fab', iconName: 'microsoft' }"/>&nbsp;Login with ASU Account
        </v-btn>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faGoogle, faMicrosoft} from "@fortawesome/free-brands-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";
import ContestantRequests from "@/utils/requests/ContestantRequests";
import GoogleLogin from "@/utils/requests/GoogleLogin";
import MicrosoftLogin from "@/utils/requests/MicrosoftLogin";
import User, {ProfileStatus, UserType} from "@/models/User";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";
import ContestantProfile from "@/components/contestant/ContestantProfile.vue";
import OrganizerProfile from "@/components/organizer/OrganizerProfile.vue";

library.add(faGoogle, faMicrosoft);

export default defineComponent({
    name: "Profile",
    components: {
        OrganizerProfile,
        ContestantProfile,
        FontAwesomeIcon,
    },
    data() {
        return {
            profile: new User(),
            contestantProfile: null,
            organizerProfile: null,
            adminProfile: null,
            team: null,
        }
    },
    async mounted() {
        this.profile = await this.tokenLogin();

        if (this.checkUserType(UserType.Contestant)) {
            this.contestantProfile = await ContestantRequests.getProfile(await this.profile);
            this.team = await ContestantRequests.getTeam();
        }

        if (this.checkUserType(UserType.Organizer)) {
            this.organizerProfile = await OrganizerRequests.getProfile(await this.profile);
        }

        if ((this.profile.user_type_base & UserType.Organizer) != 0 &&
            (this.profile.profile_status & ProfileStatus.OrganizerFinished) == 0) {
            await this.$store.dispatch("setCurrentOrganizer", await this.organizerProfile);
            await this.$router.push("/finish-org-profile/");
        }
    },
    methods: {
        checkUserType(type: UserType): boolean {
            return (type & this.profile.user_type_base) != 0;
        },
        async loginGoogle() {
            await GoogleLogin.login(
                (await this.$gapi.login()).currentUser
            );
            window.location.reload();
        },
        async $logout() {
            if (this.profile.email.indexOf("@gmail") > -1) {
                await this.$gapi.logout();
                await GoogleLogin.logout(this.profile);
            } else {
                await this.$msal.logoutPopup();
                await MicrosoftLogin.logout(this.profile);
            }
        },
        async logout() {
            await this.$logout();
            window.location.reload();
        },
        async loginMS() {
            await this.$msal.loginPopup({
                scopes: ["openid", "profile", "User.Read"]
            })
                .then((resp: any) => MicrosoftLogin.login(resp));
            window.location.reload();
        },
        async tokenLogin(): Promise<User> {
            return await GoogleLogin.loginWithToken();
        },
        async registerAsContestant() {
            await this.$store.dispatch("setCurrentUser", await this.profile);
            await this.$router.push("/register-contestant/");
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
    /* border: #212121 solid 1px; */
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
    border: 2px #212121 solid;
}

.buttons {
    padding: 5px;
}
</style>

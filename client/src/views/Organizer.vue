<template>
    <div v-if="!notOrg">
        <div class="main" v-if="profile != null">
            <br/>
            <img class="organizerLogo" alt="organizer picture" :src="profile.avatar_url"/>
            <br/>
            <span class="organizerName"><b>{{ profile.name }}</b></span>
            <br/>
            <span class="organizerName" style="font-size: 1.5em"><b>{{ profile.roles_names.join(", ") }}</b></span>
            <br/><br/>
            <v-divider/>
            <div class="buttons">
                <v-btn @click="logout" class="text-red-darken-4">Logout</v-btn>&nbsp;
            </div>

            <div v-if="(profile.roles & 1) !== 0">
                <v-divider/>
                <DirectorOperations :director="profile"/>
            </div>
        </div>

        <div v-else style="padding-top: 20px; text-align: center;">
            <h1>Are you really an organizer? 🙂</h1>
            <v-btn @click="login" class="bg-red" style="font-size: 2em; padding: 20px">
                <FontAwesomeIcon :icon="{prefix:'fab', iconName:'google'}"/>&nbsp;Prove it
            </v-btn>
        </div>
    </div>
    <h1 v-else>Only if I had half of your spirit!</h1>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faFileAlt} from "@fortawesome/free-solid-svg-icons";
import {faGoogle} from "@fortawesome/free-brands-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";
import Organizer from "@/models/Organizer";
import DirectorOperations from "@/components/director/DirectorOperations.vue";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";

library.add(faGoogle, faFileAlt);

export default defineComponent({
    name: "Organizer",
    components: {
        DirectorOperations,
        FontAwesomeIcon
    },
    data() {
        return {
            profile: null,
            team: null,
            notOrg: false
        }
    },
    async mounted() {
        this.profile = await this.tokenLogin();
        this.team = this.profile.team;
    },
    methods: {
        async login() {
            await OrganizerRequests.googleLogin(
                (await this.$gapi.login()).currentUser
            )
                .catch(() => {
                    this.notOrg = true;
                });
            window.location.reload();
        },
        async logout() {
            await this.$gapi.logout();
            await OrganizerRequests.logout();
            window.location.reload();
        },
        async tokenLogin(): Promise<Organizer> {
            const organizer = await OrganizerRequests.login();
            if (!organizer.profile_finished) {
                await this.$router.push("/finish-org-profile");
            }

            return organizer
        },
        checkDirector(): boolean {
            console.log("roles", this.profile.roles)
            return (this.profile.roles & 1) != 0;
        }
    }

});
</script>

<style scoped>
.main {
    text-align: center;
    width: 90%;
    margin: 20px auto;
    border: #212121 solid 1px;
    border-radius: 5px;
}

.organizerName {
    font-size: 2em;
    color: #212121;
}

.organizerLogo {
    width: 125px;
    height: 125px;
    border-radius: 100%;
    border: 2px #212121 solid;
}

.buttons {
    padding: 5px;
}
</style>

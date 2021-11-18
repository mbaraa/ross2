<template>
    <div class="main">
        <div v-if="!profile.profile_finished">
            <h1>Finish your profile data:</h1>

            <h3 style="text-align: left" class="text-blue">Fill at least one field<br/>Contact Info:</h3>
            <!--            <v-text-field label="Facebook profile URL" id="fb"/>-->
            <!--            <v-text-field label="Telegram number" id="tg"/>-->
            <!--            <v-text-field label="Whatsapp number" id="wa"/>-->
            <v-text-field label="Facebook profile URL" v-model="contactInfo.facebook_url"/>
            <v-text-field label="Telegram number" v-model="contactInfo.telegram_number"/>
            <v-text-field label="Whatsapp number" v-model="contactInfo.whatsapp_number"/>

            <v-btn @click="finishProfile">Finish profile
            </v-btn>
        </div>

        <h1 v-else>What do you think you're doing? 🙂</h1>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {ContactInfo} from "@/models/User";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";

export default defineComponent({
    name: "OrganizerSignup",
    data() {
        return {
            profile: {},
            contactInfo: new ContactInfo(),
        }
    },
    async mounted() {
        this.profile = await OrganizerRequests.login();
        if (this.profile.user?.profile_finished) {
            await this.$router.push("/organizer");
        }
    },
    methods: {
        async finishProfile() {
            this.profile.profile_finished = true;
            this.profile.contact_info = this.contactInfo;

            await OrganizerRequests.finishProfile(this.profile)

            await this.$router.push("/organizer");
        }
    }
});
</script>

<style scoped>
.main {
    width: 500px;
    margin: 0 auto;
}
</style>

<template>
    <div class="main">
        <div v-if="!organizerProfile.user.profile_finished">
            <h1>Finish your profile data:</h1>

            <h3 style="text-align: left" class="text-blue">Fill at least one field<br/>Contact Info:</h3>
            <!--            <v-text-field label="Facebook profile URL" id="fb"/>-->
            <!--            <v-text-field label="Telegram number" id="tg"/>-->
            <!--            <v-text-field label="Whatsapp number" id="wa"/>-->
            <v-text-field label="Facebook profile URL" v-model="contactInfo.facebook_url"/>
            <v-text-field label="Telegram URL" v-model="contactInfo.telegram_number"/>
            <!--            <v-text-field label="Whatsapp number" v-model="contactInfo.whatsapp_number"/>-->

            <v-btn @click="finishProfile()">Finish profile
            </v-btn>
        </div>

        <h1 v-else>What do you think you're doing? 🙂</h1>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {ContactInfo, ProfileStatus} from "@/models/User";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";

export default defineComponent({
    name: "OrganizerSignup",
    data() {
        return {
            organizerProfile: this.$store.getters.getCurrentOrganizer,
            contactInfo: new ContactInfo(),
        }
    },
    mounted() {
        if (this.organizerProfile.user != null &&
            (this.organizerProfile.user.profile_status & ProfileStatus.OrganizerFinished) != 0) {
            this.$router.push("/profile");
        }
    },
    methods: {
        async finishProfile() {
            this.organizerProfile.user.contact_info = this.contactInfo;

            await OrganizerRequests.finishProfile(this.organizerProfile)

            await this.$router.push("/profile");
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

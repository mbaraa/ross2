<template>
    <div class="main">
        <div v-if="!profile.profile_finished">
            <h1>Finish your profile data:</h1>

            <!--            <h3 style="text-align: left" class="text-red">Mandatory Field</h3>-->
            <v-text-field label="University ID" v-model="profile.university_id"/>

            <!--            <h3 style="text-align: left" class="text-blue">Optional Fields<br/>Contact Info:</h3>-->
            <!--            <v-text-field label="Facebook profile URL" :v-bind="contactInfo.facebook_url"/>-->
            <!--            <v-text-field label="Telegram number" :v-bind="contactInfo.telegram_number"/>-->
            <!--            <v-text-field label="Whatsapp number" :v-bind="contactInfo.whatsapp_number" required/>-->

            <v-btn @click="finishProfile">Finish profile
            </v-btn>
        </div>

        <h1 v-else>What do you think you're doing? 🙂</h1>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {ContactInfo} from "@/models/User";
import Contestant from "@/models/Contestant";

export default defineComponent({
    name: "ContestantSignup",
    data() {
        return {
            profile: {},
            contactInfo: new ContactInfo(),
        }
    },
    async mounted() {
        this.profile = await Contestant.login();
        if (this.profile.profile_finished) {
            await this.$router.push("/profile");
        }
    },
    methods: {
        async finishProfile() {
            if (this.profile.university_id.length == 0) {
                window.alert("wrong input value 🙂")
                return;
            }
            this.profile.profile_finished = true;
            await Contestant.signup(this.profile)
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

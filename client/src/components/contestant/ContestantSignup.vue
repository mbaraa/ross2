<template>
    <div class="main">
        <div v-if="checkStatus()">
            <h1>Finish your contestant profile data:</h1>

            <!--            <h3 style="text-align: left" class="text-red">Mandatory Field</h3>-->
            <v-text-field label="University ID" v-model="contestantProfile.university_id"/>

            <label style="font-size: 1.2em">Select your gender: </label>

            <br/>
            <input type="radio" id="male" value="true" v-model="contestantProfile.gender">
            <label for="male">Male</label>
            &nbsp;
            <input type="radio" id="female" value="false" v-model="contestantProfile.gender">
            <label for="female">Female</label>
            <br/>

            <label style="font-size: 1.2em">Do you mind participating with the other gender? </label>

            <br/>
            <input type="radio" id="yes" value="true" v-model="contestantProfile.participate_with_other">
            <label for="yes">Yes, I mind</label>
            &nbsp;
            <input type="radio" id="no" value="false" v-model="contestantProfile.participate_with_other">
            <label for="no">No, I don't mind</label>
            <br/><br/>
            <!--            <h3 style="text-align: left" class="text-blue">Optional Fields<br/>Contact Info:</h3>-->
            <!--            <v-text-field label="Facebook profile URL" :v-bind="contactInfo.facebook_url"/>-->
            <!--            <v-text-field label="Telegram number" :v-bind="contactInfo.telegram_number"/>-->
            <!--            <v-text-field label="Whatsapp number" :v-bind="contactInfo.whatsapp_number" required/>-->

            <v-btn @click="finishProfile()">Finish profile
            </v-btn>
        </div>

        <h1 v-else>What do you think you're doing? 🙂</h1>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {ContactInfo, ProfileStatus} from "@/models/User";
import Contestant from "@/models/Contestant";
import ContestantRequests from "@/utils/requests/ContestantRequests";

export default defineComponent({
    name: "ContestantSignup",
    data() {
        return {
            profile: this.$store.getters.getCurrentUser,
            contestantProfile: new Contestant(),
            contactInfo: new ContactInfo(),
        }
    },
    mounted() {
        this.contestantProfile.user = this.profile;
        if (this.contestantProfile.user != null &&
            (this.contestantProfile.user.profile_status & ProfileStatus.ContestantFinished) != 0) {
            this.$router.push("/profile");
        }
    },
    methods: {
        async finishProfile() {
            if (this.contestantProfile.university_id.length == 0) {
                window.alert("wrong input value 🙂")
                return;
            }

            this.setRadioValues();
            this.contestantProfile.user.profile_status |= ProfileStatus.ContestantFinished;
            await ContestantRequests.register(this.contestantProfile)
                .catch(err => window.alert(err));

            await this.$router.push("/profile");
        },
        setRadioValues() {
            this.contestantProfile.gender = (this.contestantProfile.gender == "true");
            this.contestantProfile.participate_with_other = (this.contestantProfile.participate_with_other == "true");
        },
        checkStatus(): boolean {
            return this.contestantProfile.user.profile_status !== ProfileStatus.ContestantFinished;
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

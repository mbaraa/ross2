<template>
    <div class="main">
        <div v-if="profile !== null && !profile.profile_finished">
            <h1>Finish your profile data:</h1>

            <!--            <h3 style="text-align: left" class="text-red">Mandatory Field</h3>-->
            <v-text-field label="University ID" v-model="profile.university_id"/>

            <label style="font-size: 1.2em">Select your gender: </label>

            <input type="radio" id="male" value="true" v-model="profile.gender">
            <label for="male">Male</label>
            &nbsp;
            <input type="radio" id="female" value="false" v-model="profile.gender">
            <label for="female">Female</label>
            <br/>

            <label style="font-size: 1.2em">Do you mind participating with the other gender? </label>

            <input type="radio" id="yes" value="true" v-model="profile.participate_with_other">
            <label for="yes">Yes</label>
            &nbsp;
            <input type="radio" id="no" value="false" v-model="profile.participate_with_other">
            <label for="no">No</label>
            <br/><br/>
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
            profile: null,
            contactInfo: new ContactInfo(),
        }
    },
    async mounted() {
        this.profile = await Contestant.login();
        if ((await this.profile) != null && (await this.profile).profile_finished) {
            await this.$router.push("/profile");
        }
    },
    methods: {
        async finishProfile() {
            if (this.profile.university_id.length == 0) {
                window.alert("wrong input value 🙂")
                return;
            }
            this.setRadioValues();
            this.profile.profile_finished = true;
            await Contestant.signup(this.profile)
            await this.$router.push("/profile");
        },
        setRadioValues() {
            this.profile.gender = (this.profile.gender == "true");
            this.profile.participate_with_other = (this.profile.participate_with_other == "true");
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

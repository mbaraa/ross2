<template>
    <form @submit.prevent="signup">
        <input type="text" placeholder="uni_id" v-bind="profile.university_id"/>

        <input type="submit"/>
    </form>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import config from "@/config";

export default defineComponent({
    name: "ContestantSignup",
    data() {
        return {
            profile: {}
        }
    },
    methods: {
        async signup() {
            await fetch(`${config.backendAddress}/contestant/signup/`, {
                method: "POST",
                mode: "cors",
                body: JSON.stringify(this.profile)
            })
        }
    }
});
</script>

<style scoped>

</style>

<template>
    <div class="main bg-red-darken-3">
        <br/>
        <img class="contestLogo" :alt="organizer.name + ' profile pic'" :src="organizer.avatar_url"/>
        <h1>{{ organizer.name }}</h1>
        <v-divider/>
        <p>{{ organizer.roles_names.join(", ") }}</p>
        <v-divider/>

        <p>
            <a class="text-blue-darken-4" :href="contact.facebook_url"><FontAwesomeIcon :icon="{prefix: 'fab', iconName: 'facebook'}"/></a>&nbsp;
            <a class="text-blue" :href="contact.telegram_number"><FontAwesomeIcon :icon="{prefix: 'fab', iconName: 'telegram'}"/></a>&nbsp;
            <a class="text-green" :href="contact.whatsapp_number"><FontAwesomeIcon :icon="{prefix: 'fab', iconName: 'whatsapp'}"/></a>&nbsp;
        </p>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faFacebook, faTelegram, faWhatsapp} from "@fortawesome/free-brands-svg-icons"
import {library} from "@fortawesome/fontawesome-svg-core";
import Organizer from "@/models/Organizer";

library.add(faFacebook, faTelegram, faWhatsapp);

export default defineComponent({
    name: "OrganizerCard",
    props: {
        organizer: Organizer
    },
    components: {
        FontAwesomeIcon
    },
    data() {
        return {
            contact: {}
        }
    },
    mounted() {
        this.contact = this.organizer.contact_info;
    }
});
</script>

<style scoped>
.main {
    color: white;
    text-align: center;
    width: 350px;
    margin: 0 auto;
    height: auto;
    border-radius: 5px;
    padding: 5px;
    cursor: pointer;
}

.contestLogo {
    width: 75px;
    height: 75px;
    border-radius: 100%;
    background-color: white;
    padding: 5px;
}

h1 {
    font-size: 2em;
}

p {
    font-size: 1.3em;
}
</style>

<template>
    <div v-if="nots.length > 0">
        <div v-for="not in nots" :key="not">
            <NotificationCard :notification="not"/>
        </div>
        <div class="clear">
            <v-btn @click="clearNotifications()" icon class="bg-red">
                <FontAwesomeIcon :icon="{prefix:'fas', iconName:'trash'}"/>
            </v-btn>
        </div>
    </div>
    <div v-else>
        <h1 style="text-align: center">look closer you might see some!</h1>
        <h2 style="text-align: center; cursor: pointer" @click="refresh">or just refresh the page 😕</h2>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Notification from "@/models/Notification";
import NotificationCard from "@/components/NotificationCard.vue";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faTrash} from "@fortawesome/free-solid-svg-icons";
import {library} from "@fortawesome/fontawesome-svg-core";

library.add(faTrash)

export default defineComponent({
    name: "Notifications",
    components: {NotificationCard, FontAwesomeIcon},
    data() {
        return {
            nots: []
        }
    },
    async mounted() {
        this.nots = await Notification.getNotifications();
    },
    methods: {
        refresh() {
            window.location.reload();
        },
        async clearNotifications() {
            if (window.confirm("are you sure you want to delete your notifications?")) {
                await Notification.clearNotifications();
                window.location.reload();
            }
        }
    }
});
</script>

<style scoped>
.clear {
    padding: 10px;
    font-size: 2em;
    text-align: right;
}
</style>

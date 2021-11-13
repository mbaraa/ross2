<template>
    <div v-for="not in nots" :key="not">
        <NotificationCard :notification="not"/>
    </div>
    <h1 style="text-align: center">look closer you might see some!</h1>
    <h2 style="text-align: center; cursor: pointer" @click="refresh">or just refresh the page 😕</h2>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Notification from "@/models/Notification";
import NotificationCard from "@/components/NotificationCard.vue";

export default defineComponent({
    name: "Notifications",
    components: {NotificationCard},
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
        }
    }
});
</script>

<style scoped>

</style>

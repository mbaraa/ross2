<template>
    <div class="main bg-purple">
        <!-- non-request -->
        <div class="content" v-html="filterContent()"/>
        <!-- request -->
        <div v-if="isRequest()">
            <v-btn color="success" @click="acReq">Accept</v-btn>
            &nbsp;
            <v-btn color="error" @click="waReq">Reject</v-btn>
        </div>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Notification from "@/models/Notification";
import ContestantRequests from "@/utils/requests/ContestantRequests";

export default defineComponent({
    name: "NotificationCard",
    props: {
        notification: Notification,
    },
    methods: {
        isRequest(): boolean {
            return this.notification != null &&
                this.notification.content.length > 0 &&
                this.notification.content.startsWith("_REQ");
        },
        filterContent(): string {
            const lastUnderscore = this.notification.content.lastIndexOf("_");
            const content = this.notification.content;

            return (this.isRequest() ? content.substring(4, lastUnderscore) : content);
        },
        async acReq() {
            await ContestantRequests.acceptJoinRequest(this.notification);
        },
        async waReq() {
            await ContestantRequests.rejectJoinRequest(this.notification);
            window.location.reload();
        },
    }
});
</script>

<style scoped>
.main {
    text-align: center;
    border-radius: 5px;
    padding: 10px;
    margin: 10px;
    color: white;
}

.content {
    font-size: 1.8em;
}
</style>

<template>
    <div class="main2 bg-purple-accent-4">
        <FontAwesomeIcon class="icon" :icon="{prefix:'fas', iconName:'users'}"/>
        <h1>{{ user.name }}</h1>
        <v-divider/>
        <h2>{{ user.email }}</h2>
    </div>
</template>

<script lang="ts">

import {defineComponent} from "vue";
import User from "@/models/User";

export default defineComponent({
    name: "UserCard",
    props: {
        user: User
    }
});
</script>

<style scoped>
.main2 {
    color: white;
    text-align: center;
    width: 350px;
    margin: 0 auto;
    height: auto;
    border-radius: 5px;
    padding: 5px;
}

h1 {
    font-size: 2em;
}

.icon {
    font-size: 5em;
}
</style>

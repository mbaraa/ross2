<template>
    <DirectorGenerateTeamless/>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import DirectorGenerateTeamless from "@/components/director/DirectorGenerateTeamless.vue";

export default defineComponent({
    name: "DirectorOther",
    components: {DirectorGenerateTeamless},
});
</script>

<style scoped>

</style>

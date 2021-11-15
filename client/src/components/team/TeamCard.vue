<template>
    <div class="main2 bg-green-accent-4">
        <FontAwesomeIcon class="icon" :icon="{prefix:'fas', iconName:'users'}"/>
        <h1>{{ team.name }}</h1>
        <v-divider/>
        <h2>{{ team.members.length }} member(s)</h2>
        <h3>{{ getMembersNames() }}</h3>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Team from "@/models/Team";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import {faUsers} from "@fortawesome/free-solid-svg-icons"
import {library} from "@fortawesome/fontawesome-svg-core";
import Contestant from "@/models/Contestant";

library.add(faUsers);

export default defineComponent({
    name: "TeamCard",
    props: {
        team: Team,
    },
    components: {
        FontAwesomeIcon,
    },
    methods: {
        getMembersNames(): string {
            let names = "";
            this.team.members.forEach((member: Contestant) => {
                names += `${member.name}, `
            })
            return names.substring(0, names.length - 2);
        }
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

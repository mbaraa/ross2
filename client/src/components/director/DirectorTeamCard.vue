<template>
    <div class="main2" :class="getGenderClass()">
        <FontAwesomeIcon class="icon" :icon="{ prefix: 'fas', iconName: 'users' }"/>
        <br/>
        <v-text-field label="Team name" class="text-blue-darken-4 name" v-model="newTeam.name"
                      @change="updateTeam()" :class="getGenderBoxClass()"/>
        <h2>Gender: {{ getGender() }}</h2>
        <v-divider/>

        <h2>{{ team.members.length }} member(s)</h2>
        <div v-for="member in newTeam.members" :key="member">
            {{ member.name }}
            <v-btn :title="'remove '.concat(member.name, ' from this team')" size="6" icon color="error"
                   @click="removeMember(member)">&cross;
            </v-btn>
        </div>
        <v-divider/>
        <v-text-field class="text-white name" label="New member ID" v-model="addedContID" :class="getGenderBoxClass()"/>
        <v-btn @click="addContestantToTeam()">Add Contestant</v-btn>
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
    name: "DirectorTeamCard",
    props: {
        team: Team,
    },
    components: {
        FontAwesomeIcon,
    },
    data() {
        return {
            newTeam: this.team,
            removedMembers: this.$store.getters.getRemovedContestants,
            addedContID: "",
        }
    },
    methods: {
        getMembersNames(): string {
            let names = "";
            this.team.members.forEach((member: Contestant) => {
                names += `${member.name}, `
            })
            return names.substring(0, names.length - 2);
        },
        getGender(): string {
            if (this.team.members.length > 0) {
                const firstMember = this.team.members[0];
                return firstMember.participate_with_other ? "Any" : firstMember.gender ? "Males" : "Females";
            }
            return "";
        },
        removeMember(contestant: Contestant) {
            this.$store.dispatch("addContestantToRemoved", contestant);
            this.newTeam.members.splice(this.newTeam.members.findIndex(
                (c: Contestant) => c.id == contestant.id
            ), 1);
            this.updateTeam();
        },
        addContestantToTeam() {
            this.addedContID = +this.addedContID;

            const cont = this.removedMembers.find((c: Contestant) => c.id == this.addedContID);
            if (cont === undefined) {
                window.alert("contestant doesn't exist or was already added to a team!");
                return;
            }

            if (this.newTeam.members.find((c: Contestant) => c.id == cont.id) === undefined) { // prevent duplicates from stored conts
                this.newTeam.members.push(cont);
                this.updateTeam();
            }
            this.$store.dispatch("delContestantFromRemoved", this.addedContID);
        },
        updateTeam() {
            this.$store.dispatch("addTeam", this.newTeam);
        },
        getGenderClass(): string {
            return this.getGender() == "Males"? "males": this.getGender() == "Females"? "females": "any";
        },
        getGenderBoxClass(): string {
            return this.getGender() == "Males"? "malesBox": this.getGender() == "Females"? "femalesBox": "anyBox";
        },
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

.icon {
    font-size: 5em;
}

.name {
    font-size: 1.5em;
    width: auto;
    margin-bottom: 1px;
}

.males {
    background-color: #1E88E5;
}

.females {
    background-color: #D81B60;
}

.any {
    background-color: #43A047;
}

.malesBox {
    background-color: #1565C0;
}

.femalesBox {
    background-color: #AD1457;
}

.anyBox {
    background-color: #2E7D32;
}
</style>

<template>
    <br/>
    <h2>
        <label for="contestName">Select contest: </label>
        <select id="contestName">
            <option v-for="contest in contests" :key="contest" @click="selectContest(contest)">
                {{ contest.name }}
            </option>
        </select>

    </h2>
    <div v-if="selected !== null">
        <div>
            <v-text-field v-model="searchQuery" class="search" label="University ID"/>
            <v-btn @click="listView = !listView">
                {{ getViewType() }}
            </v-btn>
        </div>
        <div v-if="users.length > 0">
            <br/>
            <div v-if="!listView">
                <div class="grid" v-for="user in filterUsers()" :key="user">
                    <div class="bg-purple-accent-4 card">
                        <UserCard :user="user"/>
                        <v-btn @click="checkAttended(user)">Attended</v-btn>
                    </div>
                </div>
            </div>
            <div v-else>
                <v-table class="bg-purple-accent-4 text-white font-weight-bold tt">
                    <thead>
                    <tr>
                        <td>Name</td>
                        <td>University ID</td>
                        <td></td>
                    </tr>
                    </thead>
                    <tbody v-for="user in filterUsers()" :key="user">
                    <tr>
                        <td>{{ user.name }}</td>
                        <td>{{ user.email.split("@")[0] }}</td>

                        <td>
                            <v-btn @click="checkAttended(user)">Attended</v-btn>
                        </td>
                    </tr>
                    </tbody>
                </v-table>
            </div>
        </div>
        <h1 class="text-red">no users were found!</h1>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";
import Contest from "@/models/Contest";
import UserCard from "@/components/user/UserCard.vue";
import User from "@/models/User";

export default defineComponent({
    name: "AttendanceList",
    components: {UserCard},
    data() {
        return {
            searchQuery: "",
            users: [],
            contests: [],
            selected: null,
            listView: true,
        }
    },
    async mounted() {
        this.contests = await OrganizerRequests.getContests();
        this.selected = this.contests[0];
        this.users = await OrganizerRequests.getParticipantsList(this.selected);
    },
    methods: {
        filterUsers(): User[] {
            return this.users
                .filter((user: User) => user.email?.includes(this.searchQuery))
                .slice(0, this.users.length >= 10 ? 10 : this.users.length);
        },
        async checkAttended(user: User) {
            await OrganizerRequests.markParticipantAsPresent(user, this.selected);
            this.users.splice(this.users.indexOf(user), 1);
        },
        selectContest(contest: Contest) {
            this.selected = contest;
        },
        getViewType(): string {
            return "Switch to" + (this.listView ? " Grid " : " List ") + "View";
        }
    }
});
</script>

<style scoped>
.grid {
    display: inline-grid;
    padding: 20px;
}

.card {
    border-radius: 5px;
    padding: 10px;
}

.search {
    margin: 0 auto;
    width: 500px;
}

.tt {
    border: #212121 solid 1px;
}
</style>

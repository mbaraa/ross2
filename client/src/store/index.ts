import {createStore} from 'vuex'
import Contestant from "@/models/Contestant";
import Team from "@/models/Team";

export default createStore({
    state: {
        removedContestants:
            <Contestant[]>JSON.parse(<string>localStorage.getItem("removedContestants")) ?? new Array<Contestant>(),
        modifiedTeams:
            <Team[]>JSON.parse(<string>localStorage.getItem("modifiedTeams")) ?? new Array<Team>(),
    },
    mutations: {
        ADD_CONTESTANT_TO_REMOVED(state, ...contestant: Contestant[]) {
            state.removedContestants.push(...contestant);

            localStorage.setItem("removedContestants", JSON.stringify(state.removedContestants));
        },
        DEL_CONTESTANT_FROM_REMOVED(state, contestantID: number) {
            const contIndex = state.removedContestants.findIndex((c) => c.id == contestantID);
            state.removedContestants.splice(contIndex, 1);

            localStorage.setItem("removedContestants", JSON.stringify(state.removedContestants));
        },
        ADD_TEAM(state, team: Team) {
            const index = state.modifiedTeams.indexOf(team);
            if (index == -1) {
                state.modifiedTeams.push(team);
            } else {
                state.modifiedTeams[index] = team;
            }

            localStorage.setItem("modifiedTeams", JSON.stringify(state.modifiedTeams));
        }
    },
    actions: {
        addContestantToRemoved({commit}, ...contestant: Contestant[]) {
            commit("ADD_CONTESTANT_TO_REMOVED", ...contestant);
        },
        delContestantFromRemoved({commit}, contestant: Contestant) {
            commit("DEL_CONTESTANT_FROM_REMOVED", contestant)
        },
        addTeam({commit}, team: Team) {
            commit("ADD_TEAM", team);
        },
    },
    getters: {
        getRemovedContestants(state) {
            return state.removedContestants;
        },
        getModifiedTeams(state) {
            return state.modifiedTeams;
        }
    },
    modules: {}
})

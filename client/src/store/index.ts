import {createStore} from 'vuex'
import Contestant from "@/models/Contestant";
import Team from "@/models/Team";
import User from '@/models/User';
import Organizer from "@/models/Organizer";

export default createStore({
    state: {
        removedContestants:
            <Contestant[]>JSON.parse(<string>localStorage.getItem("removedContestants")) ?? new Array<Contestant>(),
        modifiedTeams:
            <Team[]>JSON.parse(<string>localStorage.getItem("modifiedTeams")) ?? new Array<Team>(),
        currentTeams: new Array<Team>(),
        currentUser: new User(),
        currentOrganizer: new Organizer(),
    },
    mutations: {
        ADD_CONTESTANT_TO_REMOVED(state, contestant: Contestant) {
            const index = state.removedContestants.findIndex((c) => c.user.id == contestant.user.id);
            if (index == -1) {
                state.removedContestants.push(contestant);
            } else {
                state.removedContestants[index] = contestant;
            }

            localStorage.setItem("removedContestants", JSON.stringify(state.removedContestants));
        },
        DEL_CONTESTANT_FROM_REMOVED(state, contestantID: number) {
            const contIndex = state.removedContestants.findIndex((c) => c.user.id == contestantID);
            state.removedContestants.splice(contIndex, 1);

            localStorage.setItem("removedContestants", JSON.stringify(state.removedContestants));
        },
        ADD_TEAM(state, team: Team) {
            const index = state.modifiedTeams.findIndex((t) => t.id == team.id);
            if (index == -1) {
                state.modifiedTeams.push(team);
            } else {
                state.modifiedTeams[index] = team;
            }

            localStorage.setItem("modifiedTeams", JSON.stringify(state.modifiedTeams));
        },
        ADD_TEAM_TO_CURRENT(state, team: Team) {
            state.modifiedTeams.push(team);
        },
        SET_CURRENT_USER(state, user: User) {
            state.currentUser = user;
        },
        SET_CURRENT_ORGANIZER(state, org: Organizer) {
            state.currentOrganizer = org;
        }
    },
    actions: {
        addContestantToRemoved({commit}, contestant: Contestant) {
            commit("ADD_CONTESTANT_TO_REMOVED", contestant);
        },
        delContestantFromRemoved({commit}, contestant: Contestant) {
            commit("DEL_CONTESTANT_FROM_REMOVED", contestant)
        },
        addTeam({commit}, team: Team) {
            commit("ADD_TEAM", team);
        },
        addTeamToCurrent({commit}, team: Team) {
            commit("ADD_TEAM_TO_CURRENT", team);
        },
        setCurrentUser({commit}, user: User) {
            commit("SET_CURRENT_USER", user);
        },
        setCurrentOrganizer({commit}, org: Organizer) {
            commit("SET_CURRENT_ORGANIZER", org)
        }
    },
    getters: {
        getRemovedContestants(state) {
            return state.removedContestants;
        },
        getModifiedTeams(state) {
            return state.modifiedTeams;
        },
        getCurrentTeams(state) {
            return state.currentTeams
        },
        getCurrentUser(state) {
            return state.currentUser;
        },
        getCurrentOrganizer(state) {
            return state.currentOrganizer;
        }
    },
    modules: {}
})

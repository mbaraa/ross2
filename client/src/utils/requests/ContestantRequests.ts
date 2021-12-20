import Team from "@/models/Team";
import Notification from "@/models/Notification";
import Contest from "@/models/Contest";
import JoinRequest from "@/models/JoinRequest";
import RequestsManager, { UserType } from "@/utils/requests/RequestsManager";
import Contestant from "@/models/Contestant";
import User from "@/models/User";

class ContestantRequests {
    public static async getTeam(): Promise<Team> {
        let team = new Team();
        await RequestsManager.makeAuthGetRequest("get-team", UserType.Contestant)
            .then(resp => resp.json())
            .then(jResp => {
                team = <Team>jResp;
                return team;
            });

        return team;
    }

    public static async acceptJoinRequest(notification: Notification): Promise<void> {
        try {
            await RequestsManager.makeAuthPostRequest("accept-join-request", UserType.Contestant, notification)
                .then(resp => resp.json())
                .then(resp => {
                    window.alert(<string>resp['err']);
                });
        } finally {
            window.location.reload();
        }
    }

    public static async rejectJoinRequest(notification: Notification): Promise<void> {
        await RequestsManager.makeAuthPostRequest("reject-join-request", UserType.Contestant, notification);
    }

    public static async checkJoinedTeam(team: Team): Promise<boolean> {
        let inTeam = false;
        await RequestsManager.makeAuthPostRequest("check-joined-team", UserType.Contestant, team)
            .then(resp => resp.json())
            .then(resp => {
                inTeam = <boolean>resp["team_status"];

                return inTeam
            });

        return inTeam;
    }

    public static async createTeam(team: Team): Promise<void> {
        await RequestsManager.makeAuthPostRequest("create-team", UserType.Contestant, team);
    }

    public static async joinAsTeamless(contest: Contest): Promise<void> {
        await RequestsManager.makeAuthPostRequest("register-as-teamless", UserType.Contestant, contest);
    }

    public static async requestJoinTeam(jr: JoinRequest): Promise<Response> {
        return await RequestsManager.makeAuthPostRequest("req-join-team", UserType.Contestant, jr)
    }

    public static async register(profile: Contestant): Promise<void> {
        await RequestsManager.makeAuthPostRequest("register", UserType.Contestant, profile);
    }

    public static async getProfile(user: User): Promise<Contestant> {
        let c = new Contestant();
        await RequestsManager.makeAuthPostRequest("profile", UserType.Contestant, user)
        .then(resp => resp.json())
        .then(resp => {
            c = resp;
            return c;
        })
        .catch(err => console.error(err));

        return c;
    }

    public static async deleteUser(): Promise<void> {
        await RequestsManager.makeAuthGetRequest("delete", UserType.Contestant);
        localStorage.removeItem("token")
    }

    public static async leaveTeam(): Promise<void> {
        await RequestsManager.makeAuthGetRequest("leave-team", UserType.Contestant);
    }

    public static async deleteTeam(team: Team): Promise<void> {
        await RequestsManager.makeAuthPostRequest("delete-team", UserType.Contestant, team);
    }
}

export default ContestantRequests;

import Team from "@/models/Team";
import Contest from "@/models/Contest";
import Contestant from "@/models/Contestant";
import config from "@/config";
import GoogleLogin from "@/utils/requests/GoogleLogin";
import RequestsManager, {UserType} from "@/utils/requests/RequestsManager";
import Organizer from "@/models/Organizer";

class OrganizerRequests {
    public static async getParticipants(contest: Contest): Promise<string> {
        let parts = "";
        await RequestsManager.makeAuthPostRequest("get-participants", UserType.Organizer, contest)
            .then(resp => resp.text())
            .then(resp => {
                parts = <string>resp;
                return parts;
            });

        return parts;
    }

    public static async sendContestOverNotifications(contest: Contest): Promise<void> {
        await RequestsManager.makeAuthPostRequest("send-sheev-notifications", UserType.Organizer, contest)
            .then(() => {
                window.alert("done :)");
            })
            .catch(() => {
                window.alert("something went wrong!");
            });
    }

    public static async updateTeams(teams: Team[], removedContestants: Contestant[]): Promise<void> {
        await RequestsManager.makeAuthPostRequest("update-teams", UserType.Organizer,
            {"teams": teams, "removed_contestants": removedContestants})
            .catch(err => window.alert(err));
    }

    public static async saveTeams(teams: Team[]): Promise<void> {
        await RequestsManager.makeAuthPostRequest("register-generated-teams", UserType.Organizer, teams);
    }

    public static async generateTeams(contest: Contest, genType: string): Promise<[Array<Team>, Array<Contestant>]> {
        let teams = new Array<Team>();
        let leftTeamless = new Array<Contestant>();

        await fetch(`${config.backendAddress}/organizer/auto-generate-teams/?gen-type=${genType}`, {
            method: "POST",
            mode: "cors",
            headers: {
                "Authorization": <string>localStorage.getItem("org_token"),
            },
            body: JSON.stringify(contest),
        })
            .then(resp => resp.json())
            .then(jResp => {
                teams = <Team[]>jResp["teams"];
                leftTeamless = <Contestant[]>jResp["left_teamless"];

                return [teams, leftTeamless];
            })

        return [teams, leftTeamless];
    }

    public static async createOrganizer(org: Organizer): Promise<void> {
        await RequestsManager.makeAuthPostRequest("add-organizer", UserType.Organizer, org);
    }

    public static async deleteOrganizer(org: Organizer): Promise<void> {
        await RequestsManager.makeAuthPostRequest("delete-organizer", UserType.Organizer, org);
    }

    public static async getSubOrganizers(): Promise<Array<Organizer>> {
        let orgs = new Array<Organizer>();

        await RequestsManager.makeAuthGetRequest("get-sub-organizers", UserType.Organizer)
            .then(resp => {
                orgs = resp.json();
                return orgs;
            })
            .catch(() => {
                window.alert("something went wrong!");
            });

        return orgs;
    }

    public static async googleLogin(user: any): Promise<void> {
        await GoogleLogin.loginOrganizerWithGoogle(user);
    }

    public static async deleteContest(contest: Contest): Promise<void> {
        await RequestsManager.makeAuthPostRequest("delete-contest", UserType.Organizer, contest);
    }

    public static async createContest(contest: Contest): Promise<void> {
        await RequestsManager.makeAuthPostRequest("create-contest", UserType.Organizer, contest);
    }

    public static async getContest(contestID: number): Promise<Contest> {
        let contest = new Contest();
        await RequestsManager.makeAuthPostRequest("get-contest", UserType.Organizer, <Contest>{id: contestID})
            .then(resp => resp.json())
            .then(resp => {
                contest = <Contest>resp;
            })
            .catch(err => window.alert("oi mama" + err.message));

        return contest;
    }

    public static async getContests(): Promise<Array<Contest>> {
        let contests = new Array<Contest>();
        await RequestsManager.makeAuthGetRequest("get-contests", UserType.Organizer)
            .then(resp => resp.json())
            .then(resp => {
                contests = <Array<Contest>>resp;
            })
            .catch(err => window.alert("oi mama" + err.message));

        return contests;
    }

    public static async login(): Promise<Organizer> {
        let org: Organizer | null = new Organizer();

        await RequestsManager.makeAuthGetRequest("login", UserType.Organizer)
            .then(resp => resp.json())
            .then(jResp => {
                org = jResp as Organizer;
                return org;
            })
            .catch(() => {
                console.error(`Unauthorized!`);
                org = null;
            });

        return org;
    }

    public static async logout(): Promise<void> {
        await RequestsManager.makeAuthGetRequest("logout", UserType.Organizer);
        localStorage.removeItem("org_token")
    }

    public static async finishProfile(profile: Organizer): Promise<void> {
        await RequestsManager.makeAuthPostRequest("finish-profile", UserType.Organizer, profile);
    }
}

export default OrganizerRequests;

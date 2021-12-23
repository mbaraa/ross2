import Team from "@/models/Team";
import Contest from "@/models/Contest";
import Contestant from "@/models/Contestant";
import config from "@/config";
import RequestsManager, {UserType} from "@/utils/requests/RequestsManager";
import Organizer from "@/models/Organizer";

class OrganizerRequests {
    public static async generateTeamsPosts(data: any): Promise<string> {
        let zipFile = ""
        await RequestsManager.makeAuthPostRequest("generate-teams-posts", UserType.Organizer, data)
            .then(resp => resp.text())
            .then(resp => {
                zipFile = resp;
                return zipFile;
            })
            .catch(err => window.alert(err));

        return zipFile;
    }

    public static async getParticipants(contest: Contest): Promise<string> {
        let parts = "";
        await RequestsManager.makeAuthPostRequest("get-participants-csv", UserType.Organizer, contest)
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

    public static async generateTeams(contest: Contest, genType: string, names: string[]): Promise<[Array<Team>, Array<Contestant>]> {
        let teams = new Array<Team>();
        let leftTeamless = new Array<Contestant>();

        await fetch(`${config.backendAddress}/organizer/generate-teams/?gen-type=${genType}`, {
            method: "POST",
            mode: "cors",
            headers: {
                "Authorization": <string>localStorage.getItem("token"),
            },
            body: JSON.stringify({"contest": contest, "names": names}),
        })
            .then(resp => resp.json())
            .then(jResp => {
                teams = <Team[]>jResp["teams"];
                leftTeamless = <Contestant[]>jResp["left_teamless"];

                return [teams, leftTeamless];
            })
            .catch(err => window.alert(err));

        return [teams, leftTeamless];
    }

    public static async createOrganizer(org: Organizer): Promise<Response> {
        return await RequestsManager.makeAuthPostRequest("add-organizer", UserType.Organizer, org);
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
            .catch(err => window.alert("oi mama " + err.message));

        return contests;
    }

    public static async getProfile(): Promise<Organizer> {
        let o = new Organizer();
        await RequestsManager.makeAuthGetRequest("profile", UserType.Organizer)
            .then(resp => resp.json())
            .then(resp => {
                o = resp;
                return o;
            })
            .catch(err => console.error(err));

        return o;
    }

    public static async finishProfile(org: Organizer): Promise<void> {
        await RequestsManager.makeAuthPostRequest("finish-profile", UserType.Organizer, org);
    }
}

export default OrganizerRequests;

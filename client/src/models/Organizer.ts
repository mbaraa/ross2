import Contest from "@/models/Contest";
import config from "@/config";
import User, {ContactInfo} from "@/models/User";
import Team from "@/models/Team";
import Contestant from "@/models/Contestant";

class Organizer implements User {
    id: number | undefined;
    email: string | undefined;
    name: string | undefined;
    avatar_url: string | undefined;
    profile_finished: boolean | undefined;

    contact_info: ContactInfo | undefined;

    director: Organizer | undefined;
    contests: Contest[] | undefined;
    roles: number | undefined;
    roles_names: string[] | undefined;

    constructor() {
        this.contests = new Array<Contest>();
        this.contact_info = new ContactInfo();
        this.roles_names = new Array<string>();
    }

    public static async saveTeams(teams: Team[]): Promise<void> {
        await this.makeAuthPostRequest("register-generated-teams", teams);
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
        await this.makeAuthPostRequest("add-organizer", org);
    }

    public static async deleteOrganizer(org: Organizer): Promise<void> {
        await this.makeAuthPostRequest("delete-organizer", org);
    }

    public static async getSubOrganizers(): Promise<Array<Organizer>> {
        let orgs = new Array<Organizer>();

        await this.makeAuthGetRequest("get-sub-organizers")
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
        await fetch(`${config.backendAddress}/gauth/org-login/`, {
            method: "POST",
            mode: "cors",
            headers: {
                "Authorization": user.Zb.id_token,
            },
            body: JSON.stringify({ // only if Google didn't use such fucky names :)
                name: user.it.Se,
                avatar_url: user.it.SJ,
                email: user.it.Tt,
            })
        })
            .then(resp => resp.json())
            .then(data => {
                localStorage.setItem("org_token", <string>data["token"]);
            });
    }

    public static async deleteContest(contest: Contest): Promise<void> {
        await this.makeAuthPostRequest("delete-contest", contest);
    }

    public static async createContest(contest: Contest): Promise<void> {
        await this.makeAuthPostRequest("create-contest", contest);
    }

    public static async getContests(): Promise<Array<Contest>> {
        let contests = new Array<Contest>();
        await this.makeAuthGetRequest("get-contests")
            .then(resp => resp.json())
            .then(resp => {
                contests = <Array<Contest>>resp;
            })
            .catch(err => window.alert("oi mama" + err.message));

        return contests;
    }

    public static async login(): Promise<Organizer> {
        let org: Organizer | null = new Organizer();

        await this.makeAuthGetRequest("login")
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
        await this.makeAuthGetRequest("logout");
        localStorage.removeItem("org_token")
    }

    public static async finishProfile(profile: Organizer): Promise<void> {
        await this.makeAuthPostRequest("finish-profile", profile);
    }

    private static async makeAuthGetRequest(action: string): Promise<any> {
        return this.makeRequest("GET", action, null);
    }

    private static async makeAuthPostRequest(action: string, body: any): Promise<any> {
        return this.makeRequest("POST", action, body)
    }

    private static async makeRequest(method: string, action: string, body: any): Promise<any> {
        return fetch(`${config.backendAddress}/organizer/${action}/`, {
            method: method,
            mode: "cors",
            headers: {
                "Authorization": <string>localStorage.getItem("org_token"),
            },
            body: method == "POST" ? JSON.stringify(body) : null,
        })
    }
}

export default Organizer;

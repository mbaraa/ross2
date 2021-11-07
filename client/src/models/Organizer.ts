import Contest from "@/models/Contest";
import config from "@/config";
import User, {ContactInfo} from "@/models/User";

class Organizer {
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
        const _ = "lol";
        // super();
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

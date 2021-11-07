import Team from "@/models/Team";
import User, {ContactInfo} from "@/models/User";
import config from "@/config";

class Contestant implements User {
    id: number | undefined;
    email: string | undefined;
    name: string | undefined;
    avatar_url: string | undefined;
    profile_finished: boolean | undefined;

    contact_info: ContactInfo | undefined;

    university_id: string | undefined;
    team: Team | undefined;
    major_name: string | undefined;

    teamlessed_at: Date | undefined;
    teamless_contest_id: number | undefined;

    constructor() {
        const _ = "lol";
    }

    public static async googleLogin(user: any): Promise<void> {
        await fetch(`${config.backendAddress}/gauth/cont-login/`, {
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
                localStorage.setItem("token", <string>data["token"]);
            });
    }

    public static async login(): Promise<Contestant> {
        let cont: Contestant | null = new Contestant();

        await this.makeAuthGetRequest("login")
            .then(resp => resp.json())
            .then(jResp => {
                cont = jResp as Contestant;
                return cont;
            })
            .catch(() => {
                console.error(`Unauthorized!`);
                cont = null;
            });

        return cont;
    }

    public static async signup(profile: Contestant): Promise<void> {
        await this.makeAuthPostRequest("signup", profile);
    }

    public static async logout(): Promise<void> {
        await this.makeAuthGetRequest("logout");
        localStorage.removeItem("token")
    }

    public static async deleteUser(): Promise<void> {
        await this.makeAuthGetRequest("delete");
        localStorage.removeItem("token")
    }

    public static async leaveTeam(): Promise<void> {
        await this.makeAuthGetRequest("leave-team");
    }

    public static async deleteTeam(team: Team): Promise<void> {
        await this.makeAuthPostRequest("delete-team", team);
    }

    private static async makeAuthGetRequest(action: string): Promise<any> {
        return this.makeRequest("GET", action, null);
    }

    private static async makeAuthPostRequest(action: string, body: any): Promise<any> {
        return this.makeRequest("POST", action, body)
    }

    private static async makeRequest(method: string, action: string, body: any): Promise<any> {
        return fetch(`${config.backendAddress}/contestant/${action}/`, {
            method: method,
            mode: "cors",
            headers: {
                "Authorization": <string>localStorage.getItem("token"),
            },
            body: method == "POST" ? JSON.stringify(body) : null,
        })
    }
}


export default Contestant;

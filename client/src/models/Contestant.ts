import Team from "@/models/Team";
import User from "@/models/User";
import config from "@/config";

class Contestant extends User {
    university_id: string | undefined;
    team: Team | undefined;
    major_name: string | undefined;

    teamlessed_at: Date | undefined;
    teamless_contest_id: number | undefined;

    constructor() {
        super();
    }

    public static async login(): Promise<Contestant> {
        let cont: Contestant | null = new Contestant();

        await fetch(`${config.backendAddress}/contestant/login/`, {
            method: "GET",
            mode: "cors",
            headers: {
                "Authorization": <string>localStorage.getItem("token"),
            },
        })
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

    public static async logout(): Promise<void> {
        await fetch(`${config.backendAddress}/contestant/logout/`, {
            method: "GET",
            mode:"cors",
            headers: {
                "Authorization": <string>localStorage.getItem("token"),
            }
        })
        localStorage.removeItem("token")
    }

    public static async deleteUser(): Promise<void> {
        await fetch(`${config.backendAddress}/contestant/delete/`, {
            method: "GET",
            mode:"cors",
            headers: {
                "Authorization": <string>localStorage.getItem("token"),
            }
        })
        localStorage.removeItem("token")
    }
}


export default Contestant;

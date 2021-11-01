import User, {ContactInfo} from "@/models/User";
import Contest from "@/models/Contest";
import config from "@/config";

class Organizer extends User  {
    contests: Contest[] | undefined;
    roles: number | undefined;
    roles_names: string[] | undefined;

    constructor() {
        super();
    }

    public static async login(): Promise<Organizer> {
        let org: Organizer | null = new Organizer();
        const token = localStorage.getItem("token");
        await fetch(`${config.backendAddress}/gauth/org-login/`, {
            method: "GET",
            mode: "cors",
            headers: {
                "Authorization": token != null ? token : "",
            },
        })
            .then(resp => resp.json())
            .then(jResp => {
                org = jResp as Organizer;
                return org;
            })
            .catch(err => {
                console.log(`${err.status} Unauthorized!`)
                org = null;
            });

        return org;
    }
}

export default Organizer;

import Contest from "@/models/Contest";
import User, {ContactInfo} from "@/models/User";

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
}

export default Organizer;

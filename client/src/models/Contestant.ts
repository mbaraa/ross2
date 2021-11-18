import Team from "@/models/Team";
import User, {ContactInfo} from "@/models/User";

class Contestant implements User {
    id: number | undefined;
    email: string | undefined;
    name: string | undefined;
    avatar_url: string | undefined;
    profile_finished: boolean | undefined;
    contact_info: ContactInfo;

    university_id: string | undefined;
    team: Team | undefined;
    team_id: number | undefined;
    major_name: string | undefined;

    teamlessed_at: Date | undefined;
    teamless_contest_id: number | undefined;

    gender: boolean | undefined;
    participate_with_other: boolean | undefined;

    constructor() {
        this.contact_info = new ContactInfo();
        this.teamlessed_at = new Date();
    }
}


export default Contestant;

import Team from "./Team";
import User from "./User";

class Contestant {
    user: User;

    team?: Team;
    team_id?: number;
    major_name?: string;

    teamlessed_at?: Date;
    teamless_contest_id?: number;

    gender: boolean;
    participate_with_other: boolean;

    constructor() {
        this.gender = false;
        this.participate_with_other = false;
        this.user = new User();
        this.teamlessed_at = new Date();
    }
}


export default Contestant;

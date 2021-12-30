import Contest from "./Contest";
import Contestant from "./Contestant";

class Team {
    id?: number;
    name?: string;
    leader: Contestant;
    leader_id?: number;
    contests: Contest[];
    members: Contestant[];
    join_id?: number;

    inTeam?: boolean; // once you go spaghetti, you can't turn back :)

    constructor() {
        this.leader = new Contestant();
        this.contests = new Array<Contest>();
        this.members = new Array<Contestant>();
    }
}

export default Team;

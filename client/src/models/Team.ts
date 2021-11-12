import Contest from "@/models/Contest";
import Contestant from "@/models/Contestant";

class Team {
    id: number | undefined;
    name: string | undefined;
    leader: Contestant;
    contests: Contest[];
    members: Contestant[];

    constructor() {
        this.leader = new Contestant();
        this.contests = new Array<Contest>();
        this.members = new Array<Contestant>();
    }
}

export default Team;

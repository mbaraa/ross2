import Contest from "@/models/Contest";
import Contestant from "@/models/Contestant";

class Team {
    id: number;
    name: string;
    leader: Contestant;
    contests: Contest[];
    members: Contestant[];

    constructor(id: number, name: string, leader: Contestant, contests: Contest[], members: Contestant[]) {
        this.id = id;
        this.name = name;
        this.leader = leader;
        this.contests = contests;
        this.members = members;
    }
}

export default Team;

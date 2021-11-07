import Contest from "@/models/Contest";
import Contestant from "@/models/Contestant";

class Team {
    id: number | undefined;
    name: string | undefined;
    leader: Contestant | undefined;
    contests: Contest[] | undefined;
    members: Contestant[] | undefined;

    constructor() {
        const _ = "lol";
    }
}

export default Team;

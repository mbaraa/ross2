import Team from "@/models/Team";
import Contestant from "@/models/Contestant";
import Organizer from "@/models/Organizer";
import config from "@/config";

class Contest {
    id: number | undefined;
    name: string | undefined;
    starts_at: Date | undefined;
    duration: number | undefined;
    location: string | undefined;
    logo_path: string | undefined;
    description: string | undefined;
    participation_conditions: ParticipationConditions | undefined;
    teams: Team[] | undefined;
    organizers: Organizer[] | undefined;
    teamless_contestants: Contestant[] | undefined;

    constructor() {
        const _ = "lol";
    }

    public static async getContestFromServer(contestId: number): Promise<Contest> {
        let contest = new Contest();
        await fetch(`${config.backendAddress}/contest/single/${contestId}`, {
            method: "GET",
            mode: "cors"
        })
            .then(resp => resp.json())
            .then(data => {
                contest = data as Contest;

                return contest;
            })
            .catch(err => console.log(err));

        return contest;
    }

    public static async getContestsFromServer(): Promise<Contest[]> {
        let contests = new Array<Contest>();

        await fetch(`${config.backendAddress}/contest/all/`, {
            method: "GET",
            mode: "cors"
        })
            .then(resp => resp.json())
            .then(data => {
                contests = data["contests"] as Contest[];
                return contests;
            })
            .catch(err => console.log(err));

        return contests;
    }
}

export class ParticipationConditions {
    majors: number;
    majors_names: string[];
    min_team_members: number;
    max_team_members: number;

    constructor(majors: number, majors_names: string[], min_team_members: number, max_team_members: number) {
        this.majors = majors;
        this.majors_names = majors_names;
        this.min_team_members = min_team_members;
        this.max_team_members = max_team_members;
    }
}

export default Contest;

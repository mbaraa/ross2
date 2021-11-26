import Team from "@/models/Team";
import Contestant from "@/models/Contestant";
import Organizer from "@/models/Organizer";
import config from "@/config";

class Contest {
    id: number | undefined;
    name: string | undefined;
    starts_at: number | undefined; // timestamp goes brr
    registration_ends: number | undefined;
    duration: number | undefined;
    location: string | undefined;
    logo_path: string | undefined;
    description: string | undefined;
    participation_conditions: ParticipationConditions | undefined;
    teams_hidden?: boolean;
    teams: Team[] | undefined;
    organizers: Organizer[] | undefined;
    teamless_contestants: Contestant[] | undefined;

    constructor() {
        this.participation_conditions = new ParticipationConditions();
        this.organizers = new Array<Organizer>();
        this.teams = new Array<Team>();
        this.teamless_contestants = new Array<Contestant>();
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
    majors: number | undefined;
    majors_names: string[] | undefined;
    min_team_members: number | undefined;
    max_team_members: number | undefined;
}

export default Contest;

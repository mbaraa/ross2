import Team from "./Team";
import Contestant from "./Contestant";
import Organizer from "./Organizer";
import config from "../config";

class Contest {
  id: number;
  name: string | undefined;
  starts_at: number; // timestamp goes brr
  registration_ends: number | undefined;
  duration: number | undefined;
  location: string | undefined;
  logo_path: string | undefined;
  description: string | undefined;
  participation_conditions: ParticipationConditions;
  teams_hidden: boolean;
  teams: Team[];
  organizers: Organizer[] | undefined;
  teamless_contestants: Contestant[];

  constructor() {
    this.participation_conditions = new ParticipationConditions();
    this.organizers = [];
    this.teams = [];
    this.teamless_contestants = [];
    this.starts_at = 0;
    this.id = 0;
    this.starts_at = new Date().getTime();
    this.registration_ends = new Date().getTime();
    this.teams_hidden = false;
  }

  public static async getContestFromServer(
    contestId: number
  ): Promise<Contest> {
    let contest = new Contest();
    await fetch(`${config.backendAddress}/contest/single/${contestId}`, {
      method: "GET",
      mode: "cors",
    })
      .then((resp) => resp.json())
      .then((data) => {
        contest = data as Contest;
        return contest;
      })
      .catch((err) => console.log(err));

    return contest;
  }

  public static async getContestsFromServer(): Promise<Contest[]> {
    let contests = new Array<Contest>();

    await fetch(`${config.backendAddress}/contest/all/`, {
      method: "GET",
      mode: "cors",
    })
      .then((resp) => resp.json())
      .then((data) => {
        contests = data["contests"] as Contest[];
        return contests;
      })
      .catch((err) => console.log(err));

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

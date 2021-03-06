import Team from "./Team";
import Contest from "./Contest";

export default class JoinRequest {
    requested_team: Team;
    requested_team_id?: number;
    request_message: string;
    requested_contest_id?: number;
    requested_contest: Contest;
    requested_team_join_id?: string 

    constructor() {
        this.requested_team = new Team();
        this.request_message = "";
        this.requested_contest = new Contest();
    }
}

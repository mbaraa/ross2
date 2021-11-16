import Team from "@/models/Team";
import Contest from "@/models/Contest";

export default class JoinRequest {
    requested_team: Team;
    requested_team_id: number | undefined;
    request_message: string;
    requested_contest_id: number | undefined;
    requested_contest: Contest;

    constructor() {
        this.requested_team = new Team();
        this.request_message = "";
        this.requested_contest = new Contest();
    }
}

import Team from "@/models/Team";

export default class JoinRequest {
    requested_team: Team;
    requested_team_id: number | undefined;
    request_message: string;

    constructor() {
        this.requested_team = new Team();
        this.request_message = "";
    }
}

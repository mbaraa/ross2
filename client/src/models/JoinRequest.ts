import Contestant from "@/models/Contestant";
import Team from "@/models/Team";

export default class JoinRequest {
    requester: Contestant;
    requested_team: Team;
    request_message: string;

    constructor() {
        this.requester = new Contestant();
        this.requested_team = new Team();
        this.request_message = "";
    }
}

import Team from "../../models/Team";
import Contest from "../../models/Contest";
import Contestant from "../../models/Contestant";
import config from "../../config";
import RequestsManager, { UserType } from "./RequestsManager";
import Organizer, { OrganizerRole } from "../../models/Organizer";
import User from "../../models/User";

class OrganizerRequests {
  public static async getOrgRoles(
    organizerID: number,
    contestID: number
  ): Promise<any> {
    let roles: any;

    await RequestsManager.makeAuthPostRequest(
      "get-org-roles",
      UserType.Organizer,
      {
        contest_id: contestID,
        organizer_id: organizerID,
      }
    )
      .then((resp) => resp.json())
      .then((resp) => {
        roles = resp;
        return roles;
      })
      .catch((err) => console.error(err));

    return roles;
  }

  public static async checkOrgRole(
    contestID: number,
    organizerID: number,
    roles: OrganizerRole
  ): Promise<boolean> {
    let ok = false;
    await RequestsManager.makeAuthPostRequest(
      "check-role",
      UserType.Organizer,
      {
        contest_id: contestID,
        organizer_id: organizerID,
        roles: roles,
      }
    )
      .then((resp) => resp.json())
      .then((resp) => {
        ok = resp as boolean;
        return ok;
      })
      .catch((err) => {
        console.error(err);
        return false;
      });

    return ok;
  }

  public static async markParticipantAsPresent(
    user: User,
    contest: Contest
  ): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "mark-participant-as-present",
      UserType.Organizer,
      {
        user: user,
        contest: contest,
      }
    );
  }

  public static async getParticipantsList(contest: Contest): Promise<User[]> {
    let users = new Array<User>();

    await RequestsManager.makeAuthPostRequest(
      "get-participants",
      UserType.Organizer,
      contest
    )
      .then((resp) => resp.json())
      .then((resp) => {
        users = resp;
        return users;
      })
      .catch((err) => console.error(err));

    return users;
  }

  public static async generateTeamsPosts(data: any): Promise<string> {
    let zipFile = "";
    await RequestsManager.makeAuthPostRequest(
      "generate-teams-posts",
      UserType.Organizer,
      data
    )
      .then((resp) => resp.text())
      .then((resp) => {
        zipFile = resp;
        return zipFile;
      })
      .catch((err) => window.alert(err));

    return zipFile;
  }

  public static async getTeamsCSV(contest: Contest): Promise<string> {
    let parts = "";
    await RequestsManager.makeAuthPostRequest(
      "get-teams-csv",
      UserType.Organizer,
      contest
    )
      .then((resp) => resp.text())
      .then((resp) => {
        parts = resp as string;
        return parts;
      });

    return parts;
  }

  public static async getParticipants(contest: Contest): Promise<string> {
    let parts = "";
    await RequestsManager.makeAuthPostRequest(
      "get-participants-csv",
      UserType.Organizer,
      contest
    )
      .then((resp) => resp.text())
      .then((resp) => {
        parts = resp as string;
        return parts;
      });

    return parts;
  }

  public static async sendContestOverNotifications(
    contest: Contest
  ): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "send-sheev-notifications",
      UserType.Organizer,
      contest
    )
  }

  public static async updateTeams(
    teams: Team[],
    removedContestants: Contestant[]
  ): Promise<void> {
    await RequestsManager.makeAuthPostRequest(
      "update-teams",
      UserType.Organizer,
      { teams: teams, removed_contestants: removedContestants }
    ).catch((err) => window.alert(err));
  }

  public static async saveTeams(
    teams: Team[],
    teamless: Contestant[],
    contest: Contest
  ): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "save-teams",
      UserType.Organizer,
      {
        contest: contest,
        teams: teams,
        teamless: teamless,
      }
    );
  }

  public static async generateTeams(
    contest: Contest,
    genType: string,
    names: string[]
  ): Promise<[Array<Team>, Array<Contestant>]> {
    let teams = new Array<Team>();
    let leftTeamless = new Array<Contestant>();

    await fetch(
      `${config.backendAddress}/organizer/generate-teams/?gen-type=${genType}`,
      {
        method: "POST",
        mode: "cors",
        headers: {
          Authorization: localStorage.getItem("token") as string,
        },
        body: JSON.stringify({ contest: contest, names: names }),
      }
    )
      .then((resp) => resp.json())
      .then((jResp) => {
        teams = jResp["teams"] as Team[];
        leftTeamless = jResp["left_teamless"] as Contestant[];

        return [teams, leftTeamless];
      })
      .catch((err) => window.alert(err));

    return [teams, leftTeamless];
  }

  public static async createOrganizer(
    org: Organizer,
    contest: Contest,
    roles: OrganizerRole
  ): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "add-organizer",
      UserType.Organizer,
      {
        organizer: org,
        contest: contest,
        roles: roles,
      }
    );
  }

  public static async updateOrganizer(
    org: Organizer,
    contest: Contest,
    roles: OrganizerRole
  ): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "update-organizer",
      UserType.Organizer,
      {
        organizer: org,
        contest: contest,
        roles: roles,
      }
    );
  }

  public static async deleteOrganizer(
    org: Organizer,
    contest: Contest
  ): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "delete-organizer",
      UserType.Organizer,
      {
        organizer: org,
        contest: contest,
      }
    );
  }

  public static async getSubOrganizers(
    contest: Contest
  ): Promise<Array<Organizer>> {
    let orgs = new Array<Organizer>();

    await RequestsManager.makeAuthPostRequest(
      "get-sub-organizers",
      UserType.Organizer,
      contest
    )
      .then((resp) => {
        orgs = resp.json();
        return orgs;
      })
      .catch((err) => {
        console.error(err);
      });

    return orgs;
  }

  public static async deleteContest(contest: Contest): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "delete-contest",
      UserType.Organizer,
      contest
    );
  }

  public static async updateContest(contest: Contest): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "update-contest",
      UserType.Organizer,
      contest
    );
  }

  public static async createContest(contest: Contest): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "create-contest",
      UserType.Organizer,
      contest
    );
  }

  public static async getContest(contestID: number): Promise<Contest> {
    let contest = new Contest();
    await RequestsManager.makeAuthPostRequest(
      "get-contest",
      UserType.Organizer,
      { id: contestID } as Contest
    )
      .then((resp) => resp.json())
      .then((resp) => {
        contest = resp as Contest;
      })
      .catch((err) => console.error(err));

    return contest;
  }

  public static async getContests(): Promise<Contest[]> {
    let contests = new Array<Contest>();
    await RequestsManager.makeAuthGetRequest("get-contests", UserType.Organizer)
      .then((resp) => resp.json())
      .then((resp) => {
        contests = resp as Contest[];
      })
      .catch((err) => window.alert("oi mama " + err.message));

    return contests;
  }

  public static async getProfile(): Promise<Organizer> {
    let o = new Organizer();
    await RequestsManager.makeAuthGetRequest("profile", UserType.Organizer)
      .then((resp) => resp.json())
      .then((resp) => {
        o = resp;
        return o;
      })
      .catch((err) => console.error(err));

    return o;
  }

  public static async finishProfile(org: Organizer): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "finish-profile",
      UserType.Organizer,
      org
    );
  }
}

export default OrganizerRequests;

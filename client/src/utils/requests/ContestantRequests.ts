import Team from "../../models/Team";
import Notification from "../../models/Notification";
import JoinRequest from "../../models/JoinRequest";
import RequestsManager, { UserType } from "./RequestsManager";
import Contestant from "../../models/Contestant";
import Contest from "../../models/Contest";
import config from "../../config";

class ContestantRequests {
  public static async getTeamByJoinID(joinID: string): Promise<Team | null> {
    return await fetch(
      `${config.backendAddress}/contestant/get-team-by-join-id/?join-id=${joinID}`,
      {
        method: "GET",
        mode: "cors",
        headers: {
          Authorization: localStorage.getItem("token") as string,
        },
      }
    )
      .then((resp) => resp.json())
      .then((resp) => {
        return resp as Team;
      })
      .catch((err) => {
        console.error(err);
        return null;
      });
  }
  public static async registerInContest(contest: Contest): Promise<string> {
    let msg = "";
    await RequestsManager.makeAuthPostRequest(
      "register-in-contest",
      UserType.Contestant,
      contest
    )
      .then((resp) => resp.json())
      .then((resp) => {
        msg = resp as string;
        return msg;
      })
      .catch((err) => console.error(err));

    return msg;
  }
  public static async checkContestJoin(contest: Contest): Promise<boolean> {
    let joined = false;

    await RequestsManager.makeAuthPostRequest(
      "check-contest-join",
      UserType.Contestant,
      contest
    )
      .then((resp) => resp.json())
      .then((resp) => {
        joined = resp as boolean;
        return joined;
      })
      .catch((err) => console.error(err));

    return joined;
  }

  public static async getTeam(): Promise<Team> {
    let team = new Team();
    await RequestsManager.makeAuthGetRequest("get-team", UserType.Contestant)
      .then((resp) => resp.json())
      .then((jResp) => {
        team = jResp as Team;
        return team;
      })
      .catch((err) => console.error(err));

    return team;
  }

  public static async acceptJoinRequest(
    notification: Notification
  ): Promise<void> {
    try {
      await RequestsManager.makeAuthPostRequest(
        "accept-join-request",
        UserType.Contestant,
        notification
      )
        .then((resp) => resp.json())
        .then((resp) => {
          window.alert(resp["err"] as string);
        });
    } finally {
      window.location.reload();
    }
  }

  public static async rejectJoinRequest(
    notification: Notification
  ): Promise<void> {
    await RequestsManager.makeAuthPostRequest(
      "reject-join-request",
      UserType.Contestant,
      notification
    );
  }

  public static async checkJoinedTeam(team: Team): Promise<boolean> {
    let inTeam = false;
    await RequestsManager.makeAuthPostRequest(
      "check-joined-team",
      UserType.Contestant,
      team
    )
      .then((resp) => resp.json())
      .then((resp) => {
        inTeam = resp["team_status"] as boolean;

        return inTeam;
      });

    return inTeam;
  }

  public static async createTeam(team: Team): Promise<string> {
    let respMsg = "";
    await RequestsManager.makeAuthPostRequest(
      "create-team",
      UserType.Contestant,
      team
    )
      .then((resp) => resp.text())
      .then((resp) => {
        respMsg = resp as string;
        return respMsg;
      })
      .catch((err) => console.log(err));

    return respMsg;
  }

  public static async joinAsTeamless(body: any): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "register-as-teamless",
      UserType.Contestant,
      body
    );
  }

  public static async requestJoinTeam(jr: JoinRequest): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "req-join-team",
      UserType.Contestant,
      jr
    );
  }

  public static async register(profile: Contestant): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "register",
      UserType.Contestant,
      profile
    );
  }

  public static async getProfile(): Promise<Contestant> {
    let c = new Contestant();
    await RequestsManager.makeAuthGetRequest("profile", UserType.Contestant)
      .then((resp) => resp.json())
      .then((resp) => {
        c = resp;
        return c;
      })
      .catch((err) => console.error(err));

    return c;
  }

  public static async deleteUser(): Promise<void> {
    await RequestsManager.makeAuthGetRequest("delete", UserType.Contestant);
    localStorage.removeItem("token");
  }

  public static async leaveTeam(): Promise<Response> {
    return await RequestsManager.makeAuthGetRequest(
      "leave-team",
      UserType.Contestant
    );
  }

  public static async deleteTeam(team: Team): Promise<Response> {
    return await RequestsManager.makeAuthPostRequest(
      "delete-team",
      UserType.Contestant,
      team
    );
  }
}

export default ContestantRequests;

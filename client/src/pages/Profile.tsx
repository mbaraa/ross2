import * as React from "react";
import { Button, TextField } from "@mui/material";
import Login from "../components/Shared/Login";
import ContestantRequests from "../../src/utils/requests/ContestantRequests";
import MicrosoftLogin from "../../src/utils/requests/MicrosoftLogin";
import Title from "../components/Shared/Title";
import User, { UserType } from "../models/User";
import Contestant from "../models/Contestant";
import Team, { RegisterTeam } from "../models/Team";
import Organizer from "../models/Organizer";
import OrganizerRequests from "../utils/requests/OrganizerRequests";

interface Props {
  user: User;
}

const Profile = ({ user }: Props) => {
  const checkUserType = (t: UserType): boolean => {
    return user !== null && (t & user.user_type_base) !== 0;
  };

  const [teams, setTeams] = React.useState<RegisterTeam[] | null>(
    new Array<RegisterTeam>()
  );
  const [cont, setCont] = React.useState<Contestant>(new Contestant());
  const [org, setOrg] = React.useState<Organizer>(new Organizer());
  // const [admin, setAdmin] = React.useState<Admin>(new Admin());

  React.useEffect(() => {
    (async () => {
      if (checkUserType(UserType.Contestant)) {
        const _cont = await ContestantRequests.getProfile();
        setCont(_cont);

        setTeams(await ContestantRequests.getTeams());
      }

      if (checkUserType(UserType.Organizer)) {
        const _org = await OrganizerRequests.getProfile();
        setOrg(_org);
      }

      // if (checkUserType(UserType.Admin)) {
      //   const _admin = await AdminRequests.getProfile();
      //   setAdmin(_admin);
      // }
    })();
  }, [user]);

  const [editContact, setEditContact] = React.useState(false);
  const [modified, setModified] = React.useState(false);

  if (user !== null && user.id === 0) {
    return <Title className="mb-[8px]" content="Loading..." />;
  } else if (user === null) {
    return (
      <div>
        <Title className="mb-[8px]" content="You need to Login first!" />
        <Login />
      </div>
    );
  }

  const leaveTeam = (team: Team) => {
    if (window.confirm("Are you sure you want to leave your team?")) {
      (async () => {
        const resp = await ContestantRequests.leaveTeam(team);
        if (!resp.ok) {
          window.alert("Something went wrong, try again later!");
          return;
        }
        window.location.reload();
      })();
    }
  };

  const deleteTeam = (team: Team) => {
    if (window.confirm("Are you sure you want to delete your team :)")) {
      if (team !== undefined && (team.name as string).length === 0) {
        window.alert("Woah... something went wrong :(");
        return;
      }

      (async () => {
        const resp = await ContestantRequests.deleteTeam(team as Team);
        console.log(resp);
        if (!resp.ok) {
          window.alert("Something went wrong, try again later!");
          return;
        }
        window.location.reload();
      })();
    }
  };

  const updateOrgProfile = () => {
    (async () => {
      if (modified) {
        const resp = await OrganizerRequests.finishProfile(org);
        if (!resp.ok) {
          window.alert("Something went wrong, try again later!");
          return;
        }
      }
      setModified(false);
    })();
  };

  const getJoinLink = (teamJoinID: string, contestID: number): string =>
    `https://ross2.co/contest/${contestID}/join-team/${teamJoinID}`;

  return (
    <div className="flex justify-center items-center font-Ropa">
      <div className=" grid grid-cols-1">
        <div className="border-[1px] border-[#eee] p-[16px] mb-[8px] rounded-[8px] w-[348px]">
          <div className="text-[16px] text-[#425CBA] space-y-[4px]">
            <b>Your Name: </b>
            {user.name}
          </div>
          {org.id !== 0 && (
            <>
              <hr className="pb-[10px] mt-[10px]" />
              <div className="text-[16px] text-[#425CBA] space-y-[4px] w-full">
                <b>Facebook Profile: </b>
                {editContact ? (
                  <TextField
                    className="w-[100%]"
                    variant="outlined"
                    value={org.user.contact_info.facebook_url}
                    onKeyDown={(
                      event: React.KeyboardEvent<HTMLInputElement>
                    ) => {
                      if (event.key === "Enter") {
                        updateOrgProfile();
                        setEditContact(false);
                      }
                    }}
                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                      org.user.contact_info.facebook_url = event.target.value;
                      setModified(true);
                      setOrg({ ...org });
                    }}
                    label={
                      <label className="font-Ropa text-[18px] text-indigo">
                        Facebook URL
                      </label>
                    }
                  />
                ) : (
                  <a href={org.user.contact_info.facebook_url}>
                    {org.user.contact_info.facebook_url.includes("https://")
                      ? org.user.contact_info.facebook_url.substring(
                          "https://".length
                        )
                      : org.user.contact_info.facebook_url}
                  </a>
                )}
                <div>
                  <Button
                    color="success"
                    variant="outlined"
                    onClick={() => {
                      setEditContact(!editContact);
                      if (!editContact) {
                        updateOrgProfile();
                      }
                    }}
                  >
                    <label className="normal-case cursor-pointer">Edit</label>
                  </Button>
                </div>
              </div>

              {/*******/}

              <hr className="pb-[10px] mt-[10px]" />
              <div className="text-[16px] text-[#7B83EB] space-y-[4px] w-full">
                <b>University ID for MS Teams: </b>
                {editContact ? (
                  <TextField
                    className="w-[100%]"
                    variant="outlined"
                    value={org.user.contact_info.msteams_email}
                    onKeyDown={(
                      event: React.KeyboardEvent<HTMLInputElement>
                    ) => {
                      if (event.key === "Enter") {
                        updateOrgProfile();
                        setEditContact(false);
                      }
                    }}
                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                      org.user.contact_info.msteams_email = event.target.value;
                      setModified(true);
                      setOrg({ ...org });
                    }}
                    label={
                      <label className="font-Ropa text-[18px] text-indigo">
                        University ID for MS Teams
                      </label>
                    }
                  />
                ) : (
                  <a href={org.user.contact_info.msteams_email}>
                    {org.user.contact_info.msteams_email.includes("https://")
                      ? org.user.contact_info.msteams_email.substring(
                          "https://".length
                        )
                      : org.user.contact_info.msteams_email}
                  </a>
                )}
                <div>
                  <Button
                    color="success"
                    variant="outlined"
                    onClick={() => {
                      setEditContact(!editContact);
                      if (!editContact) {
                        updateOrgProfile();
                      }
                    }}
                  >
                    <label className="normal-case cursor-pointer">Edit</label>
                  </Button>
                </div>
              </div>
            </>
          )}

          {cont !== null && teams !== null && (teams.length as number) >= 1 && (
            <>
              <hr className="pb-[10px] mt-[10px]" />
              <label className="text-[20px] text-[#ab59ab] space-y-[4px]">
                Your Teams and Contests:
              </label>
              <br />
              <ul className="list-disc text-[#ab59ab] ml-[25px]">
                {teams.map((rt: RegisterTeam) => (
                  <>
                    <li>
                      <b>Contest:</b> {rt.contest_name}
                    </li>
                    <ul className="list-square ml-[25px]">
                      <li>
                        <b>Team Name: </b>
                        {rt.team?.name}
                      </li>
                      <li>
                        <b>Team Join ID: </b>
                        {rt.team?.join_id}
                      </li>
                      <li>
                        <b>Team Members:</b>
                      </li>
                      <ul className="list-disc ml-[25px]">
                        {rt.team?.members.map((c: Contestant) => (
                          <li>{c.user.name}</li>
                        ))}
                      </ul>
                      <li>
                        <b>Team Join Link:</b>
                        <br />
                        <a
                          className="text-[#425CBA]"
                          href={getJoinLink(
                            rt.team?.join_id as string,
                            rt.contest_id as number
                          )}
                        >
                          {getJoinLink(
                            rt.team?.join_id as string,
                            rt.contest_id as number
                          )}
                        </a>
                      </li>
                    </ul>
                    {/* team control stuff */}
                    <div className="relative left-[62%] translate-x-[-50%] ml-[-8px] mt-[10px] font-Ropa">
                      <Button
                        color="info"
                        variant="outlined"
                        onClick={() => leaveTeam(rt.team as Team)}
                      >
                        <label className="normal-case cursor-pointer">
                          Leave Team
                        </label>
                      </Button>
                      {cont.user.id === rt.team?.leader_id && (
                        <>
                          {"  "}
                          <Button
                            color="error"
                            variant="outlined"
                            onClick={() => deleteTeam(rt.team as Team)}
                          >
                            <label className="normal-case cursor-pointer">
                              Delete Team
                            </label>
                          </Button>
                        </>
                      )}
                    </div>
                  </>
                ))}
              </ul>
            </>
          )}
        </div>

        <Button
          onClick={() => {
            (async () => {
              await MicrosoftLogin.logout(user);
            })();
          }}
          color="error"
          variant="outlined"
          size="large"
        >
          <label className="normal-case font-Ropa cursor-pointer">Logout</label>
        </Button>
      </div>
    </div>
  );
};

export default Profile;

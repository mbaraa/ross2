import { Button } from "@mui/material";
import * as React from "react";
import Login from "../components/Shared/Login";
import ContestantRequests from "../../src/utils/requests/ContestantRequests";
import MicrosoftLogin from "../../src/utils/requests/MicrosoftLogin";
import Title from "../components/Shared/Title";
import User, { UserType } from "../models/User";
import Contestant from "../models/Contestant";
import Team from "../models/Team";

interface Props {
  user: User;
}

const Profile = ({ user }: Props) => {
  const checkUserType = (t: UserType): boolean => {
    return user !== null && (t & user.user_type_base) !== 0;
  };

  const [cont, setCont] = React.useState<Contestant>(new Contestant());
  // const [org, setOrg] = React.useState<Organizer>(new Organizer());
  // const [admin, setAdmin] = React.useState<Admin>(new Admin());

  React.useEffect(() => {
    (async () => {
      if (checkUserType(UserType.Contestant)) {
        const _cont = await ContestantRequests.getProfile();
        setCont(_cont);
      }

      // if (checkUserType(UserType.Organizer)) {
      //   const _org = await OrganizerRequests.getProfile();
      //   setOrg(_org);
      // }

      // if (checkUserType(UserType.Admin)) {
      //   const _admin = await AdminRequests.getProfile();
      //   setAdmin(_admin);
      // }

      // TODO:
      // do this after making org signup page :)
      // if ((this.profile.user_type_base & UserType.Organizer) != 0 &&
      //       (this.profile.profile_status & ProfileStatus.OrganizerFinished) == 0) {
      //       await this.$store.dispatch("setCurrentOrganizer", await this.organizerProfile);
      //       await this.$router.push("/finish-org-profile/");
      //   }
    })();
  }, [user]);

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

  const leaveTeam = () => {
    if (window.confirm("Are you sure you want to leave your team?")) {
      (async () => {
        await ContestantRequests.leaveTeam();
        window.location.reload();
      })();
    }
  };

  const deleteTeam = () => {
    if (window.confirm("Are you sure you want to delete your team :)")) {
      if (
        cont.team === null ||
        (cont.team !== undefined && (cont.team.name as string).length === 0)
      ) {
        window.alert("Woah... something went wrong :(");
        return;
      }

      (async () => {
        await ContestantRequests.deleteTeam(cont.team as Team);
        window.location.reload();
      })();
    }
  };

  return (
    <div className="flex justify-center items-center font-Ropa">
      <div className=" grid grid-cols-1">
        <div className="border-[1px] border-[#eee] p-[16px] mb-[8px] rounded-[8px] w-[348px]">
          <div className="text-[16px] text-[#425CBA] space-y-[4px]">
            <b>Your Name: </b>
            {user.name}
          </div>

          {cont !== null &&
            cont.team !== undefined &&
            (cont.team.id as number) > 1 && (
              <>
                <hr className="pb-[10px] mt-[10px]" />
                <div className="text-[16px] text-[#425CBA] space-y-[4px]">
                  <b>Team Name: </b> {cont.team.name}
                </div>
                <hr className="pb-[10px] mt-[10px]" />
                <div className="text-[16px] text-[#425CBA] space-y-[4px]">
                  <b>Team Join ID: </b> {cont.team.join_id}
                </div>
                {cont.team.members !== null && cont.team.members.length > 1 && (
                  <>
                    <hr className="pb-[10px] mt-[10px]" />
                    <label className="text-[#425CBA] text-[16px]">
                      <b>Team Members: </b>{" "}
                      <ul>
                        {cont.team.members.map((c) => (
                          <li>{c.user.name}</li>
                        ))}
                      </ul>{" "}
                    </label>
                  </>
                )}
                <hr className="pb-[10px] mt-[10px]" />
                <div className="relative left-[62%] translate-x-[-50%]">
                  <Button color="info" variant="outlined" onClick={leaveTeam}>
                    <label className="normal-case cursor-pointer">
                      Leave Team
                    </label>
                  </Button>
                  {cont.user.id === cont.team.leader_id && (
                    <>
                      {"  "}
                      <Button
                        color="error"
                        variant="outlined"
                        onClick={deleteTeam}
                      >
                        <label className="normal-case cursor-pointer">
                          Delete Team
                        </label>
                      </Button>
                    </>
                  )}
                </div>
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

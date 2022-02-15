import Button from "../../src/components/Button";
import * as React from "react";
import Login from "../../src/components/Login";
import ContestantRequests from "../../src/utils/requests/ContestantRequests";
import MicrosoftLogin from "../../src/utils/requests/MicrosoftLogin";
import Title from "../../src/components/Title";
import User, { UserType } from "../models/User";
import Contestant from "../models/Contestant";

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
  }, []);

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

  return (
    <div className="flex justify-center items-center font-Ropa">
      <div className=" grid grid-cols-1">
        <div className="border-[1px] border-[#eee] p-[16px] mb-[8px] rounded-[8px] w-[348px]">
          <div className="text-[16px] text-[#425CBA] space-y-[4px]">
            Your Name: {user.name}
          </div>
        </div>

        {cont !== null && cont.team !== undefined && cont.team.name !== "" && (
          <div className="border-[1px] border-[#eee] p-[16px] mb-[8px] rounded-[8px] w-[348px]">
            <div className="text-[16px] text-[#425CBA] space-y-[4px]">
              Team Name: {cont.team.name}
            </div>
          </div>
        )}

        <Button
          content="Logout"
          color="#FB4646"
          onClick={() => {
            (async () => {
              await MicrosoftLogin.logout(user);
            })();
          }}
        />
      </div>
    </div>
  );
};

export default Profile;

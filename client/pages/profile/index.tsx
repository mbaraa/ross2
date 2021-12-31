import * as React from "react";
import Login from "../../src/components/Login";
import Contestant from "../../src/models/Contestant";
import User, { UserType } from "../../src/models/User";
import ContestantRequests from "../../src/utils/requests/ContestantRequests";
import MicrosoftLogin from "../../src/utils/requests/MicrosoftLogin";

interface Props {
  user: any;
}

const Profile = ({ user }: Props) => {
  let u = {
    id: 48,
    name: "Akram",
    user_type_base: 2,
    team: { id: 1, name: "Team 1" },
  };

  const [cont, setCont] = React.useState<any>(0);

  React.useEffect(() => {
    (async () => {
      const c = await ContestantRequests.getProfile();
      setCont(c);
    })();
  }, []);

  if (user === 0) {
    return <div>Loading</div>;
  } else if (user === null) {
    return <Login/>;
  }

  return (
    <div className="flex justify-center items-center">
      <div className=" grid grid-cols-1">
        <div className="border-[1px] border-[#eee] p-[16px] mb-[8px] rounded-[8px] w-[348px]">
          <div className="text-[13px] text-[#425CBA] space-y-[4px]">
            Your Name: {user.name}
          </div>
        </div>

        {cont !== 0 && cont.team.name !== "" && (
          <div className="border-[1px] border-[#eee] p-[16px] mb-[8px] rounded-[8px] w-[348px]">
            <div className="text-[13px] text-[#425CBA] space-y-[4px]">
              Team Name: {cont.team.name}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Profile;
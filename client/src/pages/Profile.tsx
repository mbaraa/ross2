import Button from "../../src/components/Button";
import * as React from "react";
import Login from "../../src/components/Login";
import ContestantRequests from "../../src/utils/requests/ContestantRequests";
import MicrosoftLogin from "../../src/utils/requests/MicrosoftLogin";
import Title from "../../src/components/Title";

interface Props {
  user: any;
}

const Profile = ({ user }: Props) => {
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

        {cont !== 0 && cont.team.name !== "" && (
          <div className="border-[1px] border-[#eee] p-[16px] mb-[8px] rounded-[8px] w-[348px]">
            <div className="text-[16px] text-[#425CBA] space-y-[4px]">
              Team Name: {cont.team.name}
            </div>
          </div>
        )}

        <Button
          className="border-[#FB4646] text-[#FB4646] hover:bg-[#FB4646] text-center text-[16px]"
          content="Logout"
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

import * as React from "react";
import Login from "../../src/components/Login";
import User, { UserType } from "../../src/models/User";
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

<<<<<<< HEAD
  if (u === 0) {
    return <div>Loading</div>;
  } else if (u === null) {
    return <div>you should login first</div>;
  } else if ((u.user_type_base & UserType.Contestant) !== 0) {
    return (
      <div className="flex justify-center items-center">
        <div className="border-[1px] border-[#eee] p-[16px] mb-[8px] rounded-[8px] w-[348px]">
          <div className="text-[13px] text-[#425CBA] space-y-[4px]">
            Your Name: {u.name}
          </div>
        </div>

        <div className="border-[1px] border-[#eee] p-[16px] mb-[8px] rounded-[8px] w-[348px]">
          <div className="text-[13px] text-[#425CBA] space-y-[4px]">
            Team Name: {u.team.name}
          </div>
        </div>
      </div>
    );
  }

  return <div>hey {u.name}</div>;
=======
const index = (): React.ReactElement => {
    // const [user, setUser] = React.useState<User>(new User());
    // React.useEffect(() => {
    //     login();
    // }, []);
    // const login = async () => {
    //     const u = await MicrosoftLogin.loginWithToken();
    //     setUser(u);
    // };

    return (<div className="font-[Poppins] absolute left-[0.2rem] top-[1.2em] w-full">
        <Login/> 
    </div>);
>>>>>>> c34836f0c50c86ab867bd1378a1ccd4bff67ae42
};

export default Profile;

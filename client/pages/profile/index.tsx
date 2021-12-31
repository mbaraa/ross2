import * as React from "react";
import Login from "../../src/components/Login";
import User from "../../src/models/User";
import MicrosoftLogin from "../../src/utils/requests/MicrosoftLogin";

interface Props {
    user: User;
}

const Profile = ({ user }: Props): React.ReactElement => {
    return (!user ?
        <Login /> :
        <>
            profile
        </>
    );
};

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
};

export default index;
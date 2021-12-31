import OAuthLogin from "../../utils/requests/OAuthLogin";
import User from "../../models/User";

class MicrosoftLogin {
    public static async login(user: any): Promise<void> {
        await OAuthLogin.login(
            {
                name: user.account.name,
                email: user.account.userName,
            },
            user.idToken.rawIdToken,
            "msauth");
    }

    public static async loginWithToken(): Promise<User> {
        return await OAuthLogin.loginWithToken("msauth");
    }

    public static async logout(user: User): Promise<void> {
        await OAuthLogin.logout(user ,"msauth");
    }
}

export default MicrosoftLogin;

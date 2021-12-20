import OAuthLogin from "@/utils/requests/OAuthLogin";
import User from "@/models/User";

class GoogleLogin {
    public static async login(user: any): Promise<void> {
        await OAuthLogin.login(
            { // only if Google didn't use such fucky names :)
                name: user.su.qf,
                avatar_url: user.su.SM,
                email: user.su.ev,
            },
            user.vc.id_token,
            "gauth");
    }

    public static async loginWithToken(): Promise<User> {
        return await OAuthLogin.loginWithToken("gauth");
    }

    public static async logout(user: User): Promise<void> {
        await OAuthLogin.logout(user ,"gauth");
    }
}

export default GoogleLogin;

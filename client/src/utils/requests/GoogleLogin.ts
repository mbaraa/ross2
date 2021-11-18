import config from "@/config";

class GoogleLogin {
    public static async loginContestantWithGoogle(user: any): Promise<void> {
        await this.loginWithGoogle(user, "cont");
    }

    public static async loginOrganizerWithGoogle(user: any): Promise<void> {
        await this.loginWithGoogle(user, "org");
    }

    private static async loginWithGoogle(user: any, userType: string): Promise<void> {
        await fetch(`${config.backendAddress}/gauth/${userType}-login/`, {
            method: "POST",
            mode: "cors",
            headers: {
                "Authorization": user.wc.id_token,
            },
            body: JSON.stringify({ // only if Google didn't use such fucky names :)
                name: user.vu.jf,
                avatar_url: user.vu.RM,
                email: user.vu.jv,
            })
        })
            .then(resp => resp.json())
            .then(data => {
                localStorage.setItem((userType == "org"? "org_token": "token"), <string>data["token"]);
            });
    }
}

export default GoogleLogin;

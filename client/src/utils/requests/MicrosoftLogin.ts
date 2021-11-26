import config from "@/config";

class MicrosoftLogin {
    public static async loginContestantWithMicrosoft(user: any): Promise<void> {
        await this.loginWithMicrosoft(user, "cont");
    }

    public static async loginOrganizerWithMicrosoft(user: any): Promise<void> {
        await this.loginWithMicrosoft(user, "org");
    }

    private static async loginWithMicrosoft(user: any, userType: string): Promise<void> {
        await fetch(`${config.backendAddress}/msauth/${userType}-login/`, {
            method: "POST",
            mode: "cors",
            headers: {
                "Authorization": user.idToken,
            },
            body: JSON.stringify({
                name: user.account.name,
                avatar_url: "/logo_500.png",
                email: user.account.username,
            })
        })
            .then(resp => resp.json())
            .then(data => {
                localStorage.setItem((userType == "org"? "org_token": "token"), <string>data["token"]);
            });
    }
}

export default MicrosoftLogin;

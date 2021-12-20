import config from "@/config";
import User from "@/models/User";

class OAuthLogin {
    public static async loginWithToken(oauthAPIEndpoint: string): Promise<User> {
        let user: User = new User();
        await fetch(`${config.backendAddress}/${oauthAPIEndpoint}/login-token/`, {
                method: "POST",
                mode: "cors",
                headers: {
                    "Authorization": localStorage.getItem("token") ?? "",
                },
            }
        )
            .then(resp => resp.json())
            .then(data => {
                user = data;
                return user;
            })
            .catch(err => console.error(err));

        return user
    }

    public static async login(user: any, idToken: string, oauthAPIEndpoint: string): Promise<void> {
        await fetch(`${config.backendAddress}/${oauthAPIEndpoint}/login/`, {
            method: "POST",
            mode: "cors",
            headers: {
                "Authorization": idToken,
            },
            body: JSON.stringify(user)
        })
            .then(resp => resp.json())
            .then(data => {
                localStorage.setItem("token", <string>data["token"]);
            })
            .catch(err => console.error(err));
    }

    public static async logout(user: User, oauthAPIEndpoint: string): Promise<void> {
        await fetch(`${config.backendAddress}/${oauthAPIEndpoint}/logout/`, {
            method: "POST",
            mode: "cors",
            headers: {
                "Authorization": localStorage.getItem("token") ?? "",
            },
            body: JSON.stringify(user),
        })
            .catch(err => console.error(err));

        localStorage.removeItem("token");
    }
}

export default OAuthLogin

import config from "../../config";

class RequestsManager {
    public static async makeAuthGetRequest(action: string, userType: UserType): Promise<any> {
        return this.makeRequest("GET", action, userType, null);
    }

    public static async makeAuthPostRequest(action: string, userType: UserType, body: any): Promise<any> {
        return this.makeRequest("POST", action, userType, body)
    }

    private static async makeRequest(method: string, action: string, userType: UserType, body: any): Promise<any> {
        return fetch(`${config.backendAddress}/${getUserTypeString(userType)}/${action}/`, {
            method: method,
            mode: "cors",
            headers: {
                "Authorization": localStorage.getItem("token") as string,
            },
            body: method === "POST" ? JSON.stringify(body) : null,
        })
    }
}

export enum UserType {
    Contestant = 0,
    Organizer = 1,
    Admin = 2,
}

function getUserTypeString(userType: UserType): string {
    switch (userType) {
        case UserType.Contestant:
            return "contestant";
        case UserType.Organizer:
            return "organizer";
        case UserType.Admin:
            return "admin";
    }
    return "";
}

export default RequestsManager;

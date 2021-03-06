import ContestantRequests from "./requests/ContestantRequests";
import OrganizerRequests from "./requests/OrganizerRequests";
import {UserType} from "../models/User";

class ActionChecker {
    private static checkToken(): boolean {
        const token = localStorage.getItem("token") as string;
        return token != null && token.length === 36;
    }

    public static checkUser(fn: () => void): void {
        if (this.checkToken()) {
            fn();
        } else {
            window.open("/profile", "_self");
        }
    }

    public static async checkContestantForAction(fn: () => void): Promise<void> {
        if (await this.checkContestant()) {
            this.checkUser(fn);
        } else {
            window.alert("you must register as contestant first!");
        }
    }

    public static async checkOrganizerForAction(fn: () => void): Promise<void> {
        if (await this.checkOrganizer()) {
            this.checkUser(fn);
        } else {
            window.alert("you are not an organizer :)");
        }
    }

    public static async checkContestant(): Promise<boolean> {
        return ((await OrganizerRequests.getProfile()).user.user_type_base & UserType.Contestant) !== 0;
    }

    public static async checkOrganizer(): Promise<boolean> {
        return ((await OrganizerRequests.getProfile()).user.user_type_base & UserType.Organizer) !== 0;
    }
}

export default ActionChecker;

import ContestantRequests from "@/utils/requests/ContestantRequests";
import OrganizerRequests from "@/utils/requests/OrganizerRequests";
import {UserType} from "@/models/User";

class ActionChecker {
    private static checkToken(): boolean {
        const token = localStorage.getItem("token") as string;
        return token != null && token.length === 36;
    }

    public static checkUser(fn: () => void): void {
        if (this.checkToken()) {
            fn();
        } else {
            window.alert("you're not logged in :)")
        }
    }

    public static async checkContestant(fn: () => void): Promise<void> {
        if (((await ContestantRequests.getProfile()).user.user_type_base & UserType.Contestant) !== 0) {
            this.checkUser(fn);
        } else {
            window.alert("you must register as contestant first!");
        }
    }

    public static async checkOrganizer(fn: () => void): Promise<void> {
        if (((await OrganizerRequests.getProfile()).user.user_type_base & UserType.Organizer) !== 0) {
            this.checkUser(fn);
        } else {
            window.alert("you are not an organizer :)");
        }
    }
}

export default ActionChecker;

import RequestsManager, {UserType} from "@/utils/requests/RequestsManager";
import Admin from "@/models/Admin";
import Organizer from "@/models/Organizer";

class AdminRequests {
    public static async getProfile(): Promise<Admin> {
        let c = new Admin();
        await RequestsManager.makeAuthGetRequest("profile", UserType.Admin)
            .then(resp => resp.json())
            .then(resp => {
                c = resp;
                return c;
            })
            .catch(err => console.error(err));

        return c;
    }

    public static async createDirector(dir: Organizer): Promise<Response> {
        return await RequestsManager.makeAuthPostRequest("add-director", UserType.Admin, dir);
    }

    public static async deleteDirector(dir: Organizer): Promise<Response> {
        return await RequestsManager.makeAuthPostRequest("delete-director", UserType.Admin, dir);
    }

    public static async getDirectors(): Promise<Organizer[]> {
        let dirs = new Array<Organizer>();

        await RequestsManager.makeAuthGetRequest("get-directors", UserType.Admin)
            .then(resp => resp.json())
            .then(resp => {
                dirs = resp;
                return dirs;
            })
            .catch(err => window.alert(err));

        return dirs;
    }
}

export default AdminRequests;

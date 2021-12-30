import config from "../../config";
import Notification from "../../models/Notification";

class NotificationRequests {
    public static async clearNotifications(): Promise<void> {
        await fetch(`${config.backendAddress}/notification/clear/`, {
            method: "GET",
            mode: "cors",
            headers: {
                "Authorization": localStorage.getItem("token") as string,
            }
        })
    }

    public static async checkNotifications(): Promise<boolean> {
        let notificationsExists = false;
        await fetch(`${config.backendAddress}/notification/check/`, {
            method: "GET",
            mode: "cors",
            headers: {
                "Authorization": localStorage.getItem("token") as string,
            }
        })
            .then(resp => resp.json())
            .then(resp => {
                notificationsExists = resp["notifications_exists"];
                return notificationsExists;
            })

        return notificationsExists;
    }

    public static async getNotifications(): Promise<Array<Notification>> {
        let notifications = new Array<Notification>();
        await fetch(`${config.backendAddress}/notification/all/`, {
            method: "GET",
            mode: "cors",
            headers: {
                "Authorization": localStorage.getItem("token") as string,
            }
        })
            .then(resp => resp.json())
            .then(jResp => {
                notifications = jResp as Notification[];
                return notifications;
            })
            .catch(err => console.log(err));

        return notifications;
    }
}

export default NotificationRequests;

import config from "@/config";

class NotificationRequests {
    public static async clearNotifications(): Promise<void> {
        await fetch(`${config.backendAddress}/notification/clear/`, {
            method: "GET",
            mode: "cors",
            headers: {
                "Authorization": <string>localStorage.getItem("token"),
            }
        })
    }

    public static async checkNotifications(): Promise<boolean> {
        let notificationsExists = false;
        await fetch(`${config.backendAddress}/notification/check/`, {
            method: "GET",
            mode: "cors",
            headers: {
                "Authorization": <string>localStorage.getItem("token"),
            }
        })
            .then(resp => resp.text())
            .then(resp => {
                notificationsExists = resp == "true";
                return notificationsExists;
            })
            .catch(() => {
                window.alert("oops I did it again!")
            });

        return notificationsExists;
    }

    public static async getNotifications(): Promise<Array<Notification>> {
        let notifications = new Array<Notification>();
        await fetch(`${config.backendAddress}/notification/all/`, {
            method: "GET",
            mode: "cors",
            headers: {
                "Authorization": <string>localStorage.getItem("token"),
            }
        })
            .then(resp => resp.json())
            .then(jResp => {
                notifications = <Notification[]>jResp;
                return notifications;
            })
            .catch(() => {
                window.alert("oops I did it again!")
            });

        return notifications;
    }
}

export default NotificationRequests;

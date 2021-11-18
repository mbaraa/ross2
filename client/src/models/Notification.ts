import config from "@/config";

class Notification {
    id: number | undefined;
    user_id: number | undefined;
    content: string | undefined;
    seen: boolean | undefined;
    seen_at: Date;

    constructor() {
        this.seen_at = new Date();
    }
}

export default Notification;

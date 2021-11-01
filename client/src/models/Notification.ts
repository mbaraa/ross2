class Notification {
    id: number;
    user_id: number;
    content: string;
    seen: boolean;
    seen_at: Date;

    constructor(id: number, user_id: number, content: string, seen: boolean, seen_at: Date) {
        this.id = id;
        this.user_id = user_id;
        this.content = content;
        this.seen = seen;
        this.seen_at = seen_at;
    }
}

export default Notification;

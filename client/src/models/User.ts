class User {
    id: number | undefined;
    email: string | undefined;
    name: string | undefined;
    avatar_url: string | undefined;
    profile_finished: boolean | undefined;

    contact_info: ContactInfo | undefined;

    constructor() {
        const _ = "lol";
    }
}

export class ContactInfo {
    facebook_url: string | undefined;
    telegram_number: string | undefined;
    whatsapp_number: string | undefined;

    constructor() {
        const _ = "lol";
    }
}

export default User;

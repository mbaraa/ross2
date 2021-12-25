import User from "@/models/User";

class Admin {
    user: User;

    constructor() {
        this.user = new User();
    }
}

export default Admin;

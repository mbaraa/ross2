import Contest from "./Contest";
import User from "./User";

class Organizer {
  user: User;

  id: number;
  director?: Organizer;
  contests?: Contest[];
  roles?: number;
  roles_names?: string[];

  constructor() {
    this.contests = new Array<Contest>();
    this.user = new User();
    this.roles_names = new Array<string>();
    this.id = 0;
  }
}

export enum OrganizerRole {
  Director = 1,
  CoreOrganizer = 2,
  ChiefJudge = 4,
  Judge = 8,
  Technical = 16,
  Coordinator = 32,
  Media = 64,
  Balloons = 128,
  Food = 256,
  Receptionist = 512,
}

export default Organizer;

import Contest from "./Contest";
import User from "./User";

class Organizer {
  id: number;
  user: User;

  director?: Organizer;
  contests?: Contest[];
  roles?: number;
  roles_names?: string[];

  constructor() {
    this.id = 0;
    this.contests = new Array<Contest>();
    this.user = new User();
    this.roles_names = new Array<string>();
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

export function checkOrgType(org: Organizer, role: OrganizerRole): boolean {
  return ((org.roles as number) & role) !== 0;
}

export default Organizer;

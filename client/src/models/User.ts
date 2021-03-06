export class ContactInfo {
  facebook_url: string;
  telegram_number: string;
  msteams_email: string;

  constructor() {
    this.facebook_url = "";
    this.telegram_number = "";
    this.msteams_email = "";
  }
}

export enum UserType {
  Fresh = 1,
  Contestant = 2,
  Organizer = 4,
  Director = 8,
  Admin = 16,
}

export enum ProfileStatus {
  Fresh = 1,
  ContestantFinished = 2,
  OrganizerFinished = 4,
}

export default class User {
  id?: number;
  email?: string;
  name?: string;
  avatar_url?: string;
  profile_status: ProfileStatus;

  user_type_base: number;
  user_type?: string;

  contact_info: ContactInfo;

  constructor() {
    this.contact_info = new ContactInfo();
    this.user_type_base = UserType.Fresh;
    this.profile_status = ProfileStatus.Fresh;
    this.id = 0;
  }
}

export function checkUserType(user: User, type: UserType): boolean {
  return user !== null && (user.user_type_base & type) !== 0;
}

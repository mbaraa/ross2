import * as React from "react";
import { MdFacebook } from "react-icons/md";
import { SiMicrosoftteams } from "react-icons/si";
import Contest from "../../../models/Contest";
import Organizer from "../../../models/Organizer";

interface Props {
  organizer: Organizer;
}

const OrganizerCard = ({ organizer }: Props): React.ReactElement => {
  return (
    <div className="p-[25px] w-[300px] h-auto rounded border-[1px] border-ross2 mr-[10px] last:mr-0 mb-[10px]">
      <label className="text-ross2">
        <b>Name: </b>
        {organizer.user.name}
      </label>
      <hr className="border-ross2 pb-[10px] mt-[10px]" />
      {organizer.user.contact_info.facebook_url !== "https://" && (
        <div className="font-Ropa text-[#4267B2]">
          <MdFacebook className="inline-block" />{" "}
          <a
            target="_blank"
            rel="noreferrer"
            href={organizer.user.contact_info?.facebook_url}
          >
            {" "}
            Facebook URL
          </a>
        </div>
      )}
      {organizer.user.contact_info.msteams_email.length !== 0 && (
        <div className="font-Ropa text-[#7B83EB]">
          <SiMicrosoftteams className="inline-block" />{" "}
          <b>MS Teams Account: </b> {organizer.user.contact_info?.msteams_email}
        </div>
      )}
    </div>
  );
};

export default OrganizerCard;

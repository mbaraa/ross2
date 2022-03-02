import * as React from "react";
import { MdFacebook } from "react-icons/md";
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
      <div className="font-Ropa text-[#4267B2]">
      <MdFacebook className="inline-block"/> <a target="_blank" rel="noreferrer" href={organizer.user.contact_info?.facebook_url} > Facebook URL</a>
      </div>
    </div>
  );
};

export default OrganizerCard;

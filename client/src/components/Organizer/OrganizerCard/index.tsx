import { Button } from "@mui/material";
import * as React from "react";
import Contest from "../../../models/Contest";
import Organizer from "../../../models/Organizer";
import OrganizerRequests from "../../../utils/requests/OrganizerRequests";

interface Props {
  organizer: Organizer;
  contestID: number;
}

const OrganizerCard = ({ organizer, contestID }: Props): React.ReactElement => {
  const deleteOrganizer = () => {
    (async () => {
      if (window.confirm("are you sure you want to delete this organizer?")) {
        await OrganizerRequests.deleteOrganizer(organizer, {
          id: contestID,
        } as Contest);
        window.location.reload();
      }
    })();
  };

  const [rolesNames, setRolesName] = React.useState([""]);
  React.useEffect(() => {
    (async () => {
      setRolesName(
        (await OrganizerRequests.getOrgRoles(organizer?.id, contestID))
          .roles_names
      );
    })();
  }, [organizer]);

  return (
    <div className="p-[25px] w-[300px] h-auto rounded border-[1px] border-ross2 mr-[10px] last:mr-0 mb-[10px] font-Ropa">
      <label className="text-ross2">
        <b>Name: </b>
        {organizer.user.name}
      </label>
      <hr className="border-ross2 pb-[10px] mt-[10px]" />
      <label className="text-ross2">
        <b>Roles: </b>
        {rolesNames.join(", ")}
      </label>
      <hr className="border-ross2 pb-[10px] mt-[10px]" />
      <Button
        variant="outlined"
        color="error"
        className="w-full"
        onClick={deleteOrganizer}
      >
        <label className="font-Ropa text-[18px] text-[#d93333] normal-case cursor-pointer">
          Delete
        </label>
      </Button>
    </div>
  );
};

export default OrganizerCard;

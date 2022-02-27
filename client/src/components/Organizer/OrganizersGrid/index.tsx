import { Button } from "@mui/material";
import * as React from "react";
import { GoArrowLeft, GoPlus } from "react-icons/go";
import Contest from "../../../models/Contest";
import Organizer from "../../../models/Organizer";
import User from "../../../models/User";
import OrganizerRequests from "../../../utils/requests/OrganizerRequests";
import CreateEditOrganizer from "../CreateEditOrganizer";
import OrganizerCard from "../OrganizerCard";

interface Props {
  user: User;
  contest: Contest;
}

const OrganizersGrid = ({ user, contest }: Props): React.ReactElement => {
  const [orgs, setOrgs] = React.useState(new Array<Organizer>());

  React.useEffect(() => {
    (async () => {
      setOrgs(await OrganizerRequests.getSubOrganizers(contest));
    })();
  }, [user]);

  const [newOrg, setNewOrg] = React.useState(false);
  const [currentOrg, setCurrentOrg] = React.useState<Organizer>(new Organizer());

  return (
    <>
      <div className="pb-[10px]">
        <Button
          startIcon={newOrg ? <GoArrowLeft /> : <GoPlus size={12} />}
          color="error"
          variant="outlined"
          size="large"
          onClick={() => {
            setCurrentOrg(new Organizer());
            setNewOrg(!newOrg);
          }}
        >
          <label className="normal-case font-Ropa cursor-pointer">
            {newOrg ? "Go Back" : "Create Organizer"}
          </label>
        </Button>
        {newOrg && (
          <div>
            <CreateEditOrganizer
              user={user}
              contest={contest}
              organizer={currentOrg}
            />
          </div>
        )}
      </div>
      {/*  */}
      <div className="grid w-full grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
        {orgs.length > 0 &&
          !newOrg &&
          orgs.map((org) => (
            <div
              onClick={() => {
                setNewOrg(true);
                setCurrentOrg(org);
              }}
              className="cursor-pointer"
            >
              <OrganizerCard key={Math.random()} organizer={org} contestID={contest.id} />
            </div>
          ))}
      </div>
    </>
  );
};

export default OrganizersGrid;

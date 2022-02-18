import * as React from "react";
import Organizer from "../../../models/Organizer";
import OrganizerRequests from "../../../utils/requests/OrganizerRequests";
import OrganizerCard from "../OrganizerCard";

const OrganizersGrid = (): React.ReactElement => {
  const [orgs, setOrgs] = React.useState(new Array<Organizer>());

  React.useEffect(() => {
    (async () => {
      setOrgs(await OrganizerRequests.getSubOrganizers());
    })();
  }, []);

  return (
    <div className="grid w-full grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
      {orgs.length > 0 &&
        orgs.map((org) => (
          <OrganizerCard key={Math.random()} organizer={org} />
        ))}
    </div>
  );
};

export default OrganizersGrid;

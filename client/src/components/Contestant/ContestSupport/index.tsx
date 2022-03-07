import * as React from "react";
import Contest from "../../../models/Contest";
import Organizer from "../../../models/Organizer";
import Title from "../../Shared/Title";
import OrganizerCard from "../OrganizerCard";

interface Props {
  orgs: Organizer[];
  contest: Contest;
}

const ContestSupport = ({ orgs, contest }: Props): React.ReactElement => {
  const filterOrgs = (): Organizer[] => {
    return orgs.filter(
      (o) =>
        (o.user.contact_info?.facebook_url as string).includes(
          "https://facebook.com"
        ) || (o.user.contact_info?.msteams_email as string).includes("20")
    );
  };

  return (
    <div className="grid w-full grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
      {filterOrgs().length > 0 ? (
        <>
          <Title content="Contact us for support" className="" />
          <br />
          {filterOrgs().map((org) => (
            <OrganizerCard key={Math.random()} organizer={org} />
          ))}
        </>
      ) : (
        <Title content="Hmm..." className="" />
      )}
    </div>
  );
};

export default ContestSupport;

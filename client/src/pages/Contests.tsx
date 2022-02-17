import { useEffect, useState } from "react";
import ContestsGrid from "../../src/components/Contestant/ContestsGrid";
import OrganizerContestCards from "../components/Organizer/ContestsTable";
import Contest from "../../src/models/Contest";
import User, { checkUserType, UserType } from "../../src/models/User";
import Title from "../components/Shared/Title";

interface Props {
  user: User;
}

const Contests = ({ user }: Props) => {
  const [contests, setContests] = useState<Contest[]>([]);

  useEffect(() => {
    (async () => {
      const c = await Contest.getContestsFromServer();
      setContests(c);
    })();
  }, []);

  if (user !== null && user.id === 0) {
    return <Title className="mb-[8px]" content="Loading..." />;
  }
  if (contests.length === 0 && !checkUserType(user, UserType.Organizer)) {
    return (
      <>
        <h1 className="text-ross2 font-bold text-[2rem] font-Ropa">
          No contests are available at this time
        </h1>
      </>
    );
  }

  if (user !== null && checkUserType(user, UserType.Organizer)) {
    return <OrganizerContestCards contests={contests} />;
  }

  return <ContestsGrid contests={contests} />;
};

export default Contests;

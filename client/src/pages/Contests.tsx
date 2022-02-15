import { useEffect, useState } from "react";
import ContestContestantCards from "../../src/components/ContestContestantCards";
import OrganizerContestCards from "../../src/components/OrganizerContestCards";
import Contest from "../../src/models/Contest";
import User, { UserType } from "../../src/models/User";
import Title from "../components/Title";

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
  if (contests.length === 0) {
    return (
      <h1 className="text-ross2 font-bold text-[2rem] font-Ropa">
        No contests are available at this time
      </h1>
    );
  }

  if (user !== null && (user.user_type_base & UserType.Organizer) !== 0) {
    return <OrganizerContestCards contests={contests}/>;
  }

  return <ContestContestantCards contests={contests} />;
};

export default Contests;

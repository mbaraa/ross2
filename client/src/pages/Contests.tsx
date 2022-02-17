import { useEffect, useState } from "react";
import ContestContestantCards from "../../src/components/ContestContestantCards";
import OrganizerContestCards from "../../src/components/OrganizerContestCards";
import Contest from "../../src/models/Contest";
import User, { checkUserType, UserType } from "../../src/models/User";
import Title from "../components/Title";
import { Button } from "@mui/material";
import { Link } from "react-router-dom";
import { GoPlus } from "react-icons/go";

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
  console.log(user);

  if (user !== null && checkUserType(user, UserType.Organizer)) {
    return <OrganizerContestCards contests={contests} />;
  }

  return <ContestContestantCards contests={contests} />;
};

export default Contests;

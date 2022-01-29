import { useEffect, useState } from "react";
import ContestContestantCards from "../../src/components/ContestContestantCards";
import Table from "../../src/components/Table";
import Contest from "../../src/models/Contest";
import User, { UserType } from "../../src/models/User";

interface Props {
  user: any;
}

const Contests = ({ user }: Props) => {
  const [contests, setContests] = useState<Contest[]>([]);

  useEffect(() => {
    (async () => {
      const c = await Contest.getContestsFromServer();
      setContests(c);
    })();
  }, []);

  if (contests.length === 0) {
    return (
      <h1 className="text-ross2 font-bold text-[2rem]">
        No contests are available at this time
      </h1>
    );
  }

  // if (user === 0) {
  //   return <div>Loading</div>;
  // } else
  if (user !== null && (user.user_type_base & UserType.Organizer) !== 0) {
    return <Table contests={contests}></Table>;
  }

  return <ContestContestantCards contests={contests} />;
};

export default Contests;

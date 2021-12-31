import { useEffect, useState } from "react";
import ContestContestantCards from "../../src/components/ContestContestantCards";
import Table from "../../src/components/Table";
import Contest from "../../src/models/Contest";
import User, { UserType } from "../../src/models/User";

interface Props {
  user: any;
}

const Contests = ({ user }: Props) => {
  const [list, setList] = useState<Contest[]>([]);

  useEffect(() => {
    setConts();
  }, []);

  const setConts = async () => {
    const c = await Contest.getContestsFromServer();
    setList(c);
  };

  if (user === 0) {
    return <div>Loading</div>;
  } else if (
    user === null ||
    (user.user_type_base & UserType.Contestant) !== 0
  ) {
    return <ContestContestantCards contests={list} />;
  }

  return (
    <div>
      <Table contests={list}></Table>
    </div>
  );
};

export default Contests;

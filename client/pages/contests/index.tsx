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
    console.log(c);
    setList(c);
  };

  if (user === 0) {
    return <div>Loading</div>;
  } else if (
    user === null ||
    (user.user_type_base & UserType.Organizer) !== 0
  ) {
    return <Table contests={list}></Table>;
  }

  return (
    <div>
      <ContestContestantCards contests={list} />
    </div>
  );
};

export default Contests;

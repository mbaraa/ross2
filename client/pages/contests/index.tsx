import { useEffect, useState } from "react";
import Table from "../../src/components/Table";
import Contest from "../../src/models/Contest";

const Contests = () => {
  const [list, setList] = useState<Contest[]>([]);

  useEffect(() => {
    setConts();
  }, []);

  const setConts = async () => {
    const c = await Contest.getContestsFromServer();
    setList(c);
  };

  return (
    <div>
      <Table contests={list}></Table>
    </div>
  );
};

export default Contests;

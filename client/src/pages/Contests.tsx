import { useEffect, useState } from "react";
import ContestsGrid from "../../src/components/Contestant/ContestsGrid";
import Contest from "../../src/models/Contest";
import User, { checkUserType, UserType } from "../../src/models/User";
import Title from "../components/Shared/Title";
import { Link } from "react-router-dom";
import { Button } from "@mui/material";
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
  if (contests.length === 0) {
    return (
      <>
        <h1 className="text-ross2 font-bold text-[2rem] font-Ropa">
          No contests are available at this time!
        </h1>
      </>
    );
  }

  return (
    <>
      <div className="grid grid-cols-2 w-full">
        {contests.length > 0 ? (
          <Title className="mb-[12px]" content="Contests"></Title>
        ) : (
          <h1 className="text-ross2 font-bold text-[2rem] font-Ropa">
            No contests are available at this time
          </h1>
        )}
        {user !== null && checkUserType(user, UserType.Organizer) && (
          <div className="absolute items-right pt-[5px] right-[52px]">
            <Link to="/contests/new">
              <Button
                startIcon={<GoPlus size={12} />}
                color="error"
                variant="outlined"
                size="large"
              >
                <label className="normal-case font-Ropa cursor-pointer">
                  Create Contest
                </label>
              </Button>
            </Link>
          </div>
        )}
      </div>

      <ContestsGrid contests={contests} />
    </>
  );
};

export default Contests;

import Contest from "../../../models/Contest";
import ContestantContestCard from "../ContestCard";

interface Props {
  contests: Contest[];
}

const ContestsGrid = ({ contests }: Props) => {
  return (
    <>
      <div className="font-Ropa grid w-full grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
        {contests.map((contest) => {
          return <ContestantContestCard key={contest.id} contest={contest} />;
        })}
      </div>
    </>
  );
};

export default ContestsGrid;

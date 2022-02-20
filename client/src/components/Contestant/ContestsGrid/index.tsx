import Contest from "../../../models/Contest";
import ContestantContestCard from "../ContestCard";

interface Props {
  contests: Contest[];
}

const ContestsGrid = ({ contests }: Props) => {
  return (
    <div>
      <div className="font-Ropa">
        {contests.map((contest) => {
          return <ContestantContestCard key={contest.id} contest={contest} />;
        })}
      </div>
    </div>
  );
};

export default ContestsGrid;

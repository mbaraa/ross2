import Contest from "../../../models/Contest";
import ContestantContestCard from "../ContestCard";
import Title from "../../Shared/Title";

interface Props {
  contests: Contest[];
}

const ContestsGrid = ({ contests }: Props) => {
  return (
    <div>
      <Title className="mb-[12px]" content="Contests"></Title>

      <div className="font-Ropa">
        {contests.map((contest) => {
          return <ContestantContestCard key={contest.id} contest={contest} />;
        })}
      </div>
    </div>
  );
};

export default ContestsGrid;

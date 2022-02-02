import Contest from "../../models/Contest";
import ContestantContestCard from "../ContestantContestCard";
import Title from "../Title";

interface Props {
  contests: Contest[];
}

const ContestContestantCards = ({ contests }: Props) => {
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

export default ContestContestantCards;

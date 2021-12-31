import ContestantContestCard from "../ContestantContestCard";
import Title from "../Title";

interface Props {
  contests: any[];
}

const ContestContestantCards = ({ contests }: Props) => {
  let contests1 = [
    {
      id: 1,
      name: "Contest 1",
      start_at: 123456788,
      location: "Amman - Jordan",
    },
    {
      id: 2,
      name: "Contest 2",
      start_at: 123456788,
      location: "Amman - Jordan",
    },
  ];

  return (
    <div>
      <Title className="mb-[12px]" content="Contests"></Title>

      <div className="">
        {contests.map((contest) => {
          return <ContestantContestCard contest={contest} />;
        })}
      </div>
    </div>
  );
};

export default ContestContestantCards;

import ContestantContestCard from "../ContestantContestCard";

interface Props {
  contests: any[];
}

const ContestContestantCards = () => {
  let contests = [
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
      <div className="">
        {contests.map((contest) => {
          return <ContestantContestCard contest={contest} />;
        })}
      </div>
    </div>
  );
};

export default ContestContestantCards;

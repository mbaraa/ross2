import { getLocaleTime } from "../../utils";

interface Props {
  contest: any;
}

const ContestantContestCard = ({ contest }: Props) => {
  return (
    <div className="float-left border-[1px] border-[#425CBA] rounded h-auto inline-block w-[280px] ml-[16px] mb-[16px]">
      <div className="p-[28px] ">
        <div className="font-[700] text-[16px] text-ross2 ">{contest.name}</div>
        <div className="font-[500] text-[13px] text-[#425CBA] ">
          {getLocaleTime(contest.starts_at)}
        </div>
        <div className="font-[500] text-[13px] text-[#425CBA] ">
          {contest.location}
        </div>
      </div>

      <div className="border-t-[1px] border-[#425CBA] py-[12px] text-[13px] font-[600] text-[#425CBA] text-center">
        Create Team
      </div>
      <div className="border-t-[1px] border-[#425CBA] py-[12px] text-[13px] font-[600] text-[#425CBA] text-center">
        Join Team
      </div>
      <div className="border-t-[1px] border-[#425CBA] py-[12px] text-[13px] font-[600] text-[#425CBA] text-center">
        Joain as Teamless
      </div>
    </div>
  );
};

export default ContestantContestCard;

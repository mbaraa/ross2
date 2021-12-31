import * as React from "react";
import config from "../../config";
import { formatDuration, getLocaleTime } from "../../utils";

interface Props {
  contest: any;
}

const ContestantContestCard = ({ contest }: Props): React.ReactElement => {
  return (
    <div className="float-left border-[1px] border-[#425CBA] rounded h-auto inline-block w-[300px] ml-[16px] mb-[16px]">
      <div className="p-[25px] ">
        <div>
          <img
            src={`${config.backendAddress}${contest.logo_path}`}
            className="rounded-[100%] w-[75px] h-[75px] inline-block"
            alt={`${contest.name} logo`}
          />

          <span className="pl-[10px] font-[700] text-[20px] text-ross2 mb-[20px]">
            {contest.name}
          </span>
          <br />

          <span>{getLocaleTime(contest.starts_at)}</span>
        </div>

        {/* <div className="">
          {team.members.map((member) => {
            return (
              <div
                className="border-[1px] border-[#eee] p-[16px] mb-[8px] rounded-[8px] "
                key={member.user.id}
              >
                <div className="text-[13px] text-[#425CBA]">
                  {member.user.name}
                </div>
              </div>
            );
          })}
        </div> */}
      </div>

      <div className="border-t-[1px] border-[#425CBA] py-[12px] text-[13px] font-[600] text-[#425CBA] text-center">
        Add Contestant
      </div>
    </div>
  );
};

export default ContestantContestCard;

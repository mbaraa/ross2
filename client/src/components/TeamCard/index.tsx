import { CardMembership } from "@mui/icons-material";
import { PropsWithRef } from "react";

interface Props {
  team: any;
}

const TeamCard = ({ team }: Props) => {
  return (
    <div className="float-left border-[1px] border-[#425CBA] rounded h-auto inline-block w-[280px] ml-[16px] mb-[16px]">
      <div className="p-[28px] ">
        <div className="font-[700] text-[16px] text-ross2 mb-[20px]">
          {team.name}
        </div>

        <div className="">
          {team.members.map((member: any) => {
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
        </div>
      </div>

      <div className="border-t-[1px] border-[#425CBA] py-[12px] text-[13px] font-[600] text-[#425CBA] text-center">
        Add Contestant
      </div>
    </div>
  );
};

export default TeamCard;

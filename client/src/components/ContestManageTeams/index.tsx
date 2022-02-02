import TeamCard from "../TeamCard";
import Button from "../Button";
import Dialog from "@mui/material/Dialog";
import Title from "../Title";
import * as React from "react";

interface Props {
  teams: any;
  teamless: any;
}

const ContestManageTeams = ({ teams, teamless }: Props) => {
  const [open, setOpen] = React.useState(false);

  return (
    <div className="font-Ropa">
      <Dialog open={open} onClose={() => setOpen(false)}>
        <div className="min-w-[348px] max-w-[348px] p-[28px]">
          <div className="mb-[28px]">
            <Title
              className="text-[18px] font-[400] mb-[16px]"
              content="Teamless Contestants"
            />

            <Button
              className="border-[#FB4646] text-[#FB4646] hover:bg-[#FB4646]"
              content="Close"
              onClick={() => setOpen(false)}
            />
          </div>

          <div>
            {teamless.map((contestant: any) => {
              return (
                <div
                  className="border-[1px] border-[#eee] p-[16px] mb-[8px] rounded-[8px] "
                  key={contestant.user.id}
                >
                  <div className="text-[13px] text-[#425CBA] space-y-[4px]">
                    <div>id: {contestant.user.id}</div>
                    <div>name: {contestant.user.name}</div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </Dialog>

      <div className="w-[100%] inline-block mb-[12px] pt-1">
        <Button
          className="float-left"
          onClick={() => setOpen(true)}
          content="View Teamless"
        />
      </div>

      <div className="">
        {teams.map((team: any) => {
          return <TeamCard key={team.id} team={team} />;
        })}
      </div>
    </div>
  );
};

export default ContestManageTeams;

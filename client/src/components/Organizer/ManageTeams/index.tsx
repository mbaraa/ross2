import * as React from "react";
import TeamCard from "../TeamCard";
import { Button } from "@mui/material";
import Dialog from "@mui/material/Dialog";
import Title from "../../Shared/Title";
import Contestant from "../../../models/Contestant";
import Team from "../../../models/Team";

interface Props {
  teams: Team[];
  teamless: Contestant[];
}

const ManageTeams = ({ teams, teamless }: Props): React.ReactElement => {
  const [open, setOpen] = React.useState(false);

  if (teams === null || teams.length === 0) {
    return <Title content="No teams are available!" className=""/>
  }

  teams.sort((ti: Team, tj: Team) => ((ti.members.length as number) < (tj.members.length as number)? -1: 1));

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
              color="error"
              variant="outlined"
              onClick={() => setOpen(false)}
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Close
              </label>
            </Button>
          </div>

          <div>
            {teamless.map((contestant: Contestant) => {
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

      <div className="w-[100%] inline-block mb-[12px]">
        <Button
          color="info"
          variant="outlined"
          onClick={() => setOpen(true)}
          size="large"
        >
          <label className="normal-case font-Ropa cursor-pointer">
            View Teamless
          </label>
        </Button>
      </div>

      <div className="grid w-full grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
        {teams.map((team: any) => {
          return <TeamCard key={Math.random()} team={team} teamless={teamless}/>;
        })}
      </div>
    </div>
  );
};

export default ManageTeams;

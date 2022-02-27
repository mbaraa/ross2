import * as React from "react";
import TeamCard from "../TeamCard";
import { Button } from "@mui/material";
import Dialog from "@mui/material/Dialog";
import Title from "../../Shared/Title";
import Contestant from "../../../models/Contestant";
import Team from "../../../models/Team";
import { GiGears } from "react-icons/gi";
import OrganizerRequests from "../../../utils/requests/OrganizerRequests";
import Contest from "../../../models/Contest";

interface Props {
  teams: Team[];
  teamless: Contestant[];
  showGender: boolean;
  contest: Contest;
}

const ManageTeams = ({
  teams,
  teamless,
  showGender,
  contest,
}: Props): React.ReactElement => {
  const [open, setOpen] = React.useState(false);

  if (teams === null || teams.length === 0) {
    return <Title content="No teams are available!" className="" />;
  }

  teams.sort((ti: Team, tj: Team) =>
    (ti.members.length as number) < (tj.members.length as number) ? -1 : 1
  );

  const saveTeams = () => {
    (async () => {
      if (
        window.confirm("are you sure of the teams you are about to register?")
      ) {
        for (const team of teams) {
          if (team.contests !== undefined) {
            team.contests.push(contest);
          } else {
            team.contests = [contest];
          }
        }
        await OrganizerRequests.saveTeams(teams, contest);
        window.location.reload();
      }
    })();
  };

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
                  <div className="text-[14px] text-[#425CBA] space-y-[4px] font-Ropa">
                    <div>
                      <b>ID: </b> {contestant.user.id}
                    </div>
                    <div>
                      <b>Name: </b> {contestant.user.name}
                    </div>
                    <div>
                      <b>Gender</b> {contestant.gender ? "Male" : "Female"}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </Dialog>

      <div className="inline-block my-[12px]">
        <div className="grid grid-cols-2 w-auto">
          <div>
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

          {teams !== null && teams.length > 0 && (
            <div className="mr-[10px]">
              <Button
                variant="outlined"
                color="error"
                size="large"
                startIcon={<GiGears size={15} />}
                onClick={saveTeams}
              >
                <label className="normal-case font-Ropa text-[15px] cursor-pointer">
                  Save Teams
                </label>
              </Button>
            </div>
          )}
        </div>
      </div>

      <div className="grid w-full grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
        {teams.map((team: any) => {
          return (
            <TeamCard
              key={Math.random()}
              team={team}
              teamless={teamless}
              showGender={showGender}
            />
          );
        })}
      </div>
    </div>
  );
};

export default ManageTeams;

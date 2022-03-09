import * as React from "react";
import Team from "../../../models/Team";
import { GoPlus } from "react-icons/go";
import { BiTrash } from "react-icons/bi";
import { Button, TextField } from "@mui/material";
import Contestant from "../../../models/Contestant";
import Title from "../../Shared/Title";

interface LabelProps {
  text: string;
}

const FieldLabel = ({ text }: LabelProps): React.ReactElement => {
  return (
    <label className="font-Ropa text-[18px] text-ross2 normal-case">
      {text}
    </label>
  );
};

interface Props {
  team: Team;
  setTeam: React.Dispatch<React.SetStateAction<Team>>;
  teamless: Contestant[];
}

const CreateEditTeam = ({
  team,
  setTeam,
  teamless,
}: Props): React.ReactElement => {
  const isEdit = team !== undefined;

  const [name, setName] = React.useState(team.name ?? "");

  const [members, setMembers] = React.useState(
    team.members ?? new Array<Contestant>()
  );

  const saveTeam = () => {
    if (team.name === "") {
      window.alert("Enter team name!");
      return;
    }

    team.name = name;
    team.members = members;

    setTeam({ ...team });
  };

  const addContestant = (cont: Contestant) => {
    const contIndex = teamless.findIndex((c) => c.user.id === cont.user.id);
    if (contIndex >= 0) {
      members.push(teamless[contIndex] as Contestant);
      teamless.splice(contIndex, 1);
    } else {
      window.alert("Contestant doesn't exist!");
    }
    setMembers(members.flat());
  };

  const removeCont = (cont: Contestant) => {
    const mi = members.findIndex((c) => c.user.id === cont.user.id);
    teamless.push(members[mi]);
    members.splice(mi, 1);
    setMembers(members.flat());
  };

  return (
    <div>
      <div className="grid md:grid-cols-2 grid-cols-1">
        {/* left side */}
        <div className="mr-[20px]">
          <div className="mb-[10px]">
            <TextField
              className="w-[100%]"
              variant="outlined"
              value={name}
              onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                setName(event.target.value);
              }}
              label={<FieldLabel text="Team Name" />}
            />
          </div>

          {members.length === 0 ? (
            <Title content="This team has no members!" className="" />
          ) : (
            <div className="">
              {members.map((contestant: Contestant) => {
                return (
                  <div
                    className="border-[1px] border-[#eee] p-[16px] mr-[8px] mb-[8px] rounded-[8px] "
                    key={contestant.user.id}
                  >
                    <div className="text-[14px] text-[#425CBA] space-y-[4px] font-Ropa">
                      <div>
                        <b>Name: </b> {contestant.user.name}
                      </div>
                      <div>
                        <b>Gender: </b> {contestant.gender ? "Male" : "Female"}
                      </div>
                      <div>
                        <Button
                          variant="outlined"
                          color={"error"}
                          startIcon={<BiTrash size={18} />}
                          onClick={() => removeCont(contestant)}
                        >
                          <label className="normal-case font-Ropa text-[20px] cursor-pointer">
                            Remove from Team
                          </label>
                        </Button>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}

          <div>
            <Button
              variant="outlined"
              color={"error"}
              disabled={members.length === 0}
              // className="w-full"
              startIcon={<GoPlus size={18} />}
              onClick={saveTeam}
            >
              <label className="normal-case font-Ropa text-[20px] cursor-pointer">
                Save
              </label>
            </Button>
          </div>
        </div>

        {/* right side */}
        <div>
          {teamless !== undefined && teamless.length > 0 ? (
            <>
              <Title content="Current Team-less contestants:" className="" />
              <div className="overflow-y-scroll h-[50vh]">
                {teamless.map((contestant: Contestant) => {
                  return (
                    <div
                      className="border-[1px] border-[#eee] p-[16px] mr-[8px] mb-[8px] rounded-[8px] "
                      key={contestant.user.id}
                    >
                      <div className="text-[14px] text-[#425CBA] space-y-[4px] font-Ropa">
                        <div>
                          <b>Name: </b> {contestant.user.name}
                        </div>
                        <div>
                          <b>Gender: </b>{" "}
                          {contestant.gender ? "Male" : "Female"}
                        </div>
                        <div>
                          <b>Can Participate With the Other Gender: </b>{" "}
                          {contestant.participate_with_other ? "Yes" : "No"}
                        </div>
                        <div>
                          <Button
                            variant="outlined"
                            color={"error"}
                            startIcon={<GoPlus size={18} />}
                            onClick={() => addContestant(contestant)}
                          >
                            <label className="normal-case font-Ropa text-[20px] cursor-pointer">
                              Add to Team
                            </label>
                          </Button>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </>
          ) : (
            <Title
              content="No Team-less contestants are available at this time!"
              className=""
            />
          )}
        </div>
      </div>
    </div>
  );
};

export default CreateEditTeam;

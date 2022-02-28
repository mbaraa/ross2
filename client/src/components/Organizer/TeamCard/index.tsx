import * as React from "react";
import { Button, TextField, Dialog } from "@mui/material";
import Team from "../../../models/Team";
import Title from "../../Shared/Title";
import Contestant from "../../../models/Contestant";

interface Props {
  team: Team;
  teamless: Contestant[];
  showGender: boolean;
  updateTeams(): void;
}

const TeamCard = ({
  team,
  teamless,
  showGender,
  updateTeams,
}: Props): React.ReactElement => {
  const setColorFromGender = (): string => {
    if (team.members.length > 0) {
      const firstMember = team.members[0];
      return firstMember.participate_with_other
        ? "#6a1b9a"
        : firstMember.gender
        ? "#00796B"
        : "#C2185B";
    }
    return "#6a1b9a";
  };

  const [color, _] = React.useState(setColorFromGender());

  const gender =
    color === "#00796B" ? "Males" : color === "#C2185B" ? "Females" : "Mixed";

  // add contestant stuff
  interface State {
    newContID: string;
  }
  const [state, setState] = React.useState<State>({
    newContID: "",
  });

  const handleChange =
    (prop: keyof State) => (event: React.ChangeEvent<HTMLInputElement>) => {
      setState({ ...state, [prop]: event.target.value });
    };

  const [openNewCont, setOpenNewCont] = React.useState(false);
  const openNewContHandler = () => {
    setOpenNewCont(true);
  };

  const closeNewContHandler = () => {
    setOpenNewCont(false);
  };

  const addContestant = () => {
    const contIndex = teamless.findIndex(
      (c) => c.user.id === Number(state.newContID)
    );
    if (contIndex >= 0) {
      team.members.push(teamless[contIndex] as Contestant);
      teamless.splice(contIndex, 1);
    } else {
      window.alert("Contestant doesn't exist!");
    }

    setState({ newContID: "" });
    closeNewContHandler();
  };

  const removeCont = (cont: Contestant) => {
    const mi = team.members.findIndex((c) => c.user.id === cont.user.id);
    teamless.push(team.members[mi]);
    team.members.splice(mi, 1);
    updateTeams();
  };

  return (
    <>
      <div
        className={`float-left border-[1px] border-[${color}] rounded h-auto inline-block w-[280px] mr-[16px] mb-[56px] p-0 font-Ropa`}
      >
        <div className="p-[28px]">
          <div className="inline-grid grid-cols-2 w-full">
            <label
              className={`font-[700] text-[20px] text-[${color}] mb-[20px]`}
            >
              {team.name}
            </label>
            <div className="relative top-0 right-[-3px] ">
              <Button
                variant="outlined"
                color="info"
                onClick={openNewContHandler}
              >
                <label className="text-[10px] normal-case cursor-pointer text-center">
                  Add Contestant
                </label>
              </Button>
            </div>
          </div>

          {showGender && (
            <label className={`text-[${color}] font-[15px]`}>
              <b>Gender:</b> {gender}
            </label>
          )}

          <div className="pt-[10px]">
            {team.members.map((member: any) => {
              return (
                <div
                  className="border-[1px] border-[#eee] p-[16px] mb-[8px] last:mb-0 rounded-[8px] "
                  key={member.user.id}
                >
                  <div className={`text-[13px] text-[${color}]`}>
                    {member.user.name}
                  </div>
                  <Button
                    color="error"
                    variant="outlined"
                    onClick={() => removeCont(member)}
                    size="large"
                  >
                    <label className="normal-case font-Ropa cursor-pointer">
                      Remove
                    </label>
                  </Button>
                </div>
              );
            })}
          </div>
        </div>
      </div>
      {/* add contestant */}
      <Dialog open={openNewCont} onClose={closeNewContHandler}>
        <div className="min-w-[348px] max-w-[348px] p-[28px]">
          <div className="mb-[28px]">
            <Title
              className="text-[18px] font-[400] mb-[16px]"
              content="Add New Contestant"
            />
          </div>
          <TextField
            label="New Contestant ID"
            value={state.newContID}
            onChange={handleChange("newContID")}
          />
          <br />
          <br />
          <div className=" space-x-[4px] float-right">
            <Button
              color="error"
              variant="outlined"
              onClick={closeNewContHandler}
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Cancel
              </label>
            </Button>
            <Button
              color="info"
              variant="outlined"
              onClick={addContestant}
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Add Contestant
              </label>
            </Button>
          </div>
        </div>
      </Dialog>
    </>
  );
};

export default TeamCard;

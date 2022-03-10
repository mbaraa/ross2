import * as React from "react";
import { Button, TextField } from "@mui/material";
import Dialog from "@mui/material/Dialog";
import Title from "../../Shared/Title";
import Contestant from "../../../models/Contestant";
import Team from "../../../models/Team";
import { GiGears } from "react-icons/gi";
import { GoPlus } from "react-icons/go";
import { BsArrowLeft } from "react-icons/bs";
import { BiTrash } from "react-icons/bi";
import OrganizerRequests from "../../../utils/requests/OrganizerRequests";
import Contest from "../../../models/Contest";

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
  teams: Team[];
  teamless: Contestant[];
  showGender: boolean;
  contest: Contest;
  updateTeams(): void;
}

const ManageTeams = ({
  teams,
  teamless,
  showGender,
  contest,
  updateTeams,
}: Props): React.ReactElement => {
  const [open, setOpen] = React.useState(false);

  teams.sort((ti: Team, tj: Team) =>
    (ti.members.length as number) < (tj.members.length as number) ? -1 : 1
  );

  const [modTeams, setModTeams] = React.useState(new Array<Team>());
  const saveTeams = () => {
    (async () => {
      if (
        window.confirm("are you sure of the teams you are about to register?")
      ) {
        await OrganizerRequests.saveTeams(modTeams, teamless, contest);
      }
    })();
    // window.location.reload();
  };

  const filterTeams = (): Team[] => {
    return teams.filter((t) => t.members.length > 0);
  };

  const getCompleteTeams = (): Team[] => {
    return filterTeams().filter(
      (t) =>
        t.members.length === contest.participation_conditions.max_team_members
    );
  };

  const getInCompleteTeams = (): Team[] => {
    return filterTeams().filter(
      (t) =>
        t.members.length !== contest.participation_conditions.max_team_members
    );
  };

  const getTeamGender = (_team: Team): string => {
    const firstMember = _team.members[0];
    return firstMember.participate_with_other
      ? "Any"
      : firstMember.gender
      ? "Males"
      : "Females";
  };

  const [isModTeam, setIsModTeam] = React.useState(false);
  const [modTeam, setModTeam] = React.useState(new Team());

  const addContestant = (cont: Contestant) => {
    const contIndex = teamless.findIndex((c) => c.user.id === cont.user.id);
    if (contIndex >= 0) {
      modTeam.members.push(teamless[contIndex] as Contestant);
      teamless.splice(contIndex, 1);
    } else {
      window.alert("Contestant doesn't exist!");
    }
    setModTeam({ ...modTeam });
  };

  const removeCont = (cont: Contestant) => {
    const mi = modTeam.members.findIndex((c) => c.user.id === cont.user.id);
    teamless.push(modTeam.members[mi]);
    modTeam.members.splice(mi, 1);
    setModTeam({ ...modTeam });
  };

  const saveTeam = () => {
    setModTeam({ ...modTeam });
    modTeams.push(modTeam);
    setModTeams(modTeams.flat());
    updateTeams();
  };

  const deleteTeam = () => {
    const ti = teams.findIndex((t) => t.id === modTeam.id);
    if (
      ti > -1 &&
      window.confirm("Are you sure that you want to delete this team?")
    ) {
      for (const m of modTeam.members) {
        const mi = modTeam.members.findIndex((c) => c.user.id === m.user.id);
        teamless.push(modTeam.members[mi]);
      }
      modTeam.members = new Array<Contestant>();
      setModTeam({ ...modTeam });

      teams.splice(ti, 1);
      updateTeams();
      setModTeam(new Team());
      setIsModTeam(false);
    }
  };

  if (teams === null || teams.length === 0) {
    return <Title content="No teams are available!" className="" />;
  }

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
            {teamless !== null &&
              teamless.map((contestant: Contestant) => {
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
        <div className="grid grid-cols-3 w-auto">
          {!isModTeam && (
            <>
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
            </>
          )}

          {isModTeam ? (
            <div className="mr-[10px] inline-block">
              <Button
                variant="outlined"
                color="error"
                size="large"
                startIcon={<BsArrowLeft size={15} />}
                onClick={() => {
                  setIsModTeam(false);
                  setModTeam(new Team());
                }}
              >
                <label className="normal-case font-Ropa text-[15px] cursor-pointer">
                  Go Back
                </label>
              </Button>
            </div>
          ) : (
            <div className="mr-[10px]">
              <Button
                variant="outlined"
                color="success"
                size="large"
                startIcon={<GoPlus size={15} />}
                onClick={() => {
                  setModTeam(new Team());
                  // modTeam.id = getNewID();
                  setModTeam({ ...modTeam });
                  setIsModTeam(true);
                }}
              >
                <label className="normal-case font-Ropa text-[15px] cursor-pointer">
                  Create New Team
                </label>
              </Button>
            </div>
          )}
        </div>
      </div>

      {/*************/}
      {/* all teams */}
      {/*************/}
      {!isModTeam && (
        <>
          <Title
            content={`${filterTeams().length} Teams, ${
              getCompleteTeams().length
            } Complete Teams`}
            className="text-indigo text-[22px] mb-[10px]"
          />

          {/* complete teams */}
          <div className="border-t-2 border-t-gray-300 mt-[10px] pb-[10px]" />
          <Title
            content="Completed Teams"
            className="text-indigo text-[22px] mb-[10px]"
          />
          <div className="font-Ropa flex flex-row flex-wrap justify-center sm:grid sm:w-full sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
            {getCompleteTeams().map((team: any) => (
              <div
                className={`float-left border-[1px] border-ross2 rounded h-auto inline-block w-[280px] mr-[16px] mb-[56px] p-0 font-Ropa hover:cursor-pointer`}
                onClick={() => {
                  setModTeam(team);
                  setIsModTeam(true);
                }}
                title="Click to edit this team"
              >
                <div className="p-[28px]">
                  <label
                    className={`font-[700] text-[20px] text-indigo mb-[20px]`}
                  >
                    {team.name}
                  </label>

                  {showGender && (
                    <label className={`text-black font-[15px]`}>
                      <b>Gender:</b> {getTeamGender(team)}
                    </label>
                  )}

                  <div className="pt-[10px]">
                    {team.members.map((member: any) => {
                      return (
                        <div
                          className="border-[1px] border-[#eee] p-[16px] mb-[8px] last:mb-0 rounded-[8px] "
                          key={member.user.id}
                        >
                          <div className={`text-[13px] text-indigo`}>
                            {member.user.name}
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* incomplete teams */}
          <div className="border-t-2 border-t-gray-300 pb-[10px]" />
          <Title
            content="Incompleted Teams"
            className="text-indigo text-[22px] mb-[10px]"
          />
          <div className="font-Ropa flex flex-row flex-wrap justify-center sm:grid sm:w-full sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
            {getInCompleteTeams().map((team: any) => (
              <div
                className={`float-left border-[1px] border-ross2 rounded h-auto inline-block w-[280px] mr-[16px] mb-[56px] p-0 font-Ropa hover:cursor-pointer`}
                onClick={() => {
                  setModTeam(team);
                  setIsModTeam(true);
                }}
                title="Click to edit this team"
              >
                <div className="p-[28px]">
                  <label
                    className={`font-[700] text-[20px] text-indigo mb-[20px]`}
                  >
                    {team.name}
                  </label>

                  {showGender && (
                    <label className={`text-black font-[15px]`}>
                      <b>Gender:</b> {getTeamGender(team)}
                    </label>
                  )}

                  <div className="pt-[10px]">
                    {team.members.map((member: any) => {
                      return (
                        <div
                          className="border-[1px] border-[#eee] p-[16px] mb-[8px] last:mb-0 rounded-[8px] "
                          key={member.user.id}
                        >
                          <div className={`text-[13px] text-indigo`}>
                            {member.user.name}
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </>
      )}

      {/********************/}
      {/* create|edit team */}
      {/********************/}
      {isModTeam && (
        <div>
          <div>
            <div className="grid md:grid-cols-2 grid-cols-1">
              {/* left side */}
              <div className="mr-[20px]">
                <div className="my-[15px]">
                  <TextField
                    className="w-[100%]"
                    variant="outlined"
                    value={modTeam.name}
                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                      modTeam.name = event.target.value;
                      setModTeam({ ...modTeam });
                    }}
                    label={<FieldLabel text="Team Name" />}
                  />
                </div>

                {modTeam.members.length === 0 ? (
                  <Title content="This team has no members!" className="" />
                ) : (
                  <div className="">
                    {modTeam.members.map((contestant: Contestant) => {
                      return (
                        <div
                          className="border-[1px] border-[#eee] p-[16px] mr-[8px] mb-[8px] rounded-[8px] "
                          key={contestant.user.id}
                        >
                          <div className="text-[15px] text-[#425CBA] space-y-[4px] font-Ropa">
                            <div>
                              <b>Name: </b> {contestant.user.name}
                            </div>
                            <div>
                              <b>Email: </b> {contestant.user.email}
                            </div>
                            {showGender && (
                              <div>
                                <b>Gender: </b>{" "}
                                {contestant.gender ? "Male" : "Female"}
                              </div>
                            )}
                            <div>
                              <Button
                                variant="outlined"
                                color={"error"}
                                startIcon={<BiTrash size={13} />}
                                onClick={() => removeCont(contestant)}
                              >
                                <label className="normal-case font-Ropa text-[13px] cursor-pointer">
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
                    color={"secondary"}
                    disabled={modTeam.members.length === 0}
                    // className="w-full"
                    startIcon={<GoPlus size={18} />}
                    onClick={saveTeam}
                  >
                    <label className="normal-case font-Ropa text-[20px] cursor-pointer">
                      Save
                    </label>
                  </Button>

                  <div className="mx-[5px] inline-grid" />
                  <Button
                    variant="outlined"
                    color={"error"}
                    disabled={modTeam.members.length === 0}
                    // className="w-full"
                    startIcon={<BiTrash size={18} />}
                    onClick={deleteTeam}
                  >
                    <label className="normal-case font-Ropa text-[20px] cursor-pointer">
                      Delete
                    </label>
                  </Button>
                </div>
              </div>

              {/* right side */}
              <div className="mt-[10px] sm:mt-0">
                {teamless !== undefined && teamless.length > 0 ? (
                  <>
                    <Title
                      content="Current Team-less contestants:"
                      className=""
                    />
                    <div className="overflow-y-scroll h-[50vh]">
                      {teamless.map((contestant: Contestant) => {
                        return (
                          <div
                            className="border-[1px] border-[#eee] p-[16px] mr-[8px] mb-[8px] rounded-[8px] "
                            key={contestant.user.id}
                          >
                            <div className="text-[15px] text-[#425CBA] space-y-[4px] font-Ropa">
                              <div>
                                <b>Name: </b> {contestant.user.name}
                              </div>
                              <div>
                                <b>Email: </b> {contestant.user.email}
                              </div>
                              <div>
                                <b>Gender: </b>{" "}
                                {contestant.gender ? "Male" : "Female"}
                              </div>
                              <div>
                                <b>Can Participate With the Other Gender: </b>{" "}
                                {contestant.participate_with_other
                                  ? "Yes"
                                  : "No"}
                              </div>
                              <div>
                                <Button
                                  variant="outlined"
                                  color={"info"}
                                  startIcon={<GoPlus size={18} />}
                                  onClick={() => addContestant(contestant)}
                                >
                                  <label className="normal-case font-Ropa text-[13px] cursor-pointer">
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
        </div>
      )}
    </div>
  );
};

export default ManageTeams;

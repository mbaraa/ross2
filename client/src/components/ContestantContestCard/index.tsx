import { getLocaleTime } from "../../utils";
import * as React from "react";
import {
  Dialog,
  FormControlLabel,
  Radio,
  RadioGroup,
  TextField,
} from "@mui/material";
import Button from "../Button";
import Title from "../Title";
import Team from "../../models/Team";
import ContestantRequests from "../../utils/requests/ContestantRequests";
import JoinRequest from "../../models/JoinRequest";
import Contestant from "../../models/Contestant";
import config from "../../config";

interface Props {
  contest: any;
}

const ContestantContestCard = ({ contest }: Props) => {
  const [team, ] = React.useState<Team>(new Team());

  interface State {
    teamName: string;
    teamID: string;
    gender: string;
    partWithOther: string;
  }
  const [state, setState] = React.useState<State>({
    teamName: "",
    teamID: "",
    gender: "",
    partWithOther: "",
  });

  const handleChange =
    (prop: keyof State) => (event: React.ChangeEvent<HTMLInputElement>) => {
      setState({ ...state, [prop]: event.target.value });
    };

  const [openJT, setOpenJT] = React.useState(false);
  const openJTHandler = () => {
    setOpenJT(true);
  };

  const closeJTHandler = () => {
    setOpenJT(false);
  };

  const [openCT, setOpenCT] = React.useState(false);
  const openCTHandler = () => {
    setOpenCT(true);
  };

  const closeCTHandler = () => {
    setOpenCT(false);
  };

  const [openJTL, setOpenJTL] = React.useState(false);
  const openJTLHandler = () => {
    setOpenJTL(true);
  };

  const closeJTLHandler = () => {
    setOpenJTL(false);
  };

  const checkRegisterEnds = (): boolean => {
    const regOver = new Date().getTime() > contest.registration_ends;
    if (regOver) {
      window.alert("sorry, the registration for this contest is over!");
    }
    return regOver;
  };

  const createTeam = async () => {
    if (checkRegisterEnds()) {
      return;
    }
    team.name = state.teamName;
    team.contests.push(contest);
    await ContestantRequests.createTeam(team);
    closeCTHandler();
    window.alert(`your team "${team.name}" was created successfully ☺️`);
    window.location.reload();
  };

  const joinTeam = async () => {
    const resp = await ContestantRequests.requestJoinTeam({
      requested_team: team,
      requested_team_id: team.id as number,
      requested_team_join_id: state.teamID,
      request_message: "",
      requested_contest_id: +contest.id,
      requested_contest: contest,
    } as JoinRequest);
    if (resp.ok) {
      window.alert(
        "request sent successfully, now wait for the team's leader to accept your request!"
      );
    } else {
      window.alert(await resp.text());
    }
    closeJTHandler();
    window.location.reload();
  };

  const [contestantProfile, setCont] = React.useState<Contestant>(
    new Contestant()
  );
  React.useEffect(() => {
    (async () => {
      const c = await ContestantRequests.getProfile();
      setCont(c);
    })();
  }, []);

  const [hasTeam, setHasTeam] = React.useState<boolean>(false);

  React.useEffect(() => {
    setHasTeam(
      contestantProfile != null && (contestantProfile.team_id as number) > 1
    );
  }, []);

  const joinAsTeamless = async () => {
    if (
      window.confirm(
        `are you sure you want to join the contest "${contest.name}" as teamless?`
      )
    ) {
      contestantProfile.gender = state.gender === "true";
      contestantProfile.participate_with_other = state.partWithOther === "true";

      await ContestantRequests.joinAsTeamless({
        contest: contest,
        contestant: contestantProfile,
      });
      window.alert(`you have registered as teamless in "${contest.name}"`);
      closeJTLHandler();
    }
  };

  return (
    <div className="float-left border-[1px] font-Ropa border-[#425CBA] rounded h-auto inline-block w-[280px] ml-[16px] mb-[16px]">
      <div className="p-[28px] ">
        <img alt={contest.name} src={`${config.backendAddress}${contest.logo_path}`} 
          className="rounded-full w-[75px] h-[75px] inline"/>

        <label className="font-[700] text-[20px] text-ross2 pl-[15px]">{contest.name}</label>
        <br/>
        <label className="font-[500] text-[13px] text-[#425CBA] ">
          {getLocaleTime(contest.starts_at)}
        </label>
        <br/>
        <label className="font-[500] text-[13px] text-[#425CBA] ">
          {contest.location}
        </label>
      </div>

      {!hasTeam && (
        <>
          <div
            className="border-t-[1px] border-[#425CBA] py-[12px] text-[13px] font-[600] text-[#425CBA] text-center cursor-pointer"
            onClick={openCTHandler}
          >
            Create Team
          </div>
          <div
            className="border-t-[1px] border-[#425CBA] py-[12px] text-[13px] font-[600] text-[#425CBA] text-center cursor-pointer"
            onClick={openJTHandler}
          >
            Join Team
          </div>
          <div
            className="border-t-[1px] border-[#425CBA] py-[12px] text-[13px] font-[600] text-[#425CBA] text-center cursor-pointer"
            onClick={openJTLHandler}
          >
            Join as Teamless
          </div>
        </>
      )}
      <Dialog open={openCT} onClose={closeCTHandler}>
        <div className="min-w-[348px] max-w-[348px] p-[28px]">
          <div className="mb-[28px]">
            <Title
              className="text-[18px] font-[400] mb-[16px]"
              content="Create Team"
            />
          </div>
          <TextField
            label="Team Name"
            value={state.teamName}
            onChange={handleChange("teamName")}
          />
          <br />
          <br />
          <div className=" space-x-[4px] float-right">
            <Button
              className="border-[#FB4646] text-[#FB4646] hover:bg-[#FB4646]"
              content="Cancel"
              onClick={closeCTHandler}
            />
            <Button className="" content="Generate" onClick={createTeam} />
          </div>
        </div>
      </Dialog>

      {/*********************/}

      <Dialog open={openJT} onClose={closeJTHandler}>
        <div className="min-w-[348px] max-w-[348px] p-[28px]">
          <div className="mb-[28px]">
            <Title
              className="text-[18px] font-[400] mb-[16px]"
              content="Join Team"
            />
          </div>
          <TextField
            label="Team ID"
            value={state.teamID}
            onChange={handleChange("teamID")}
          />
          <br />
          <br />
          <div className=" space-x-[4px] float-right">
            <Button
              className="border-[#FB4646] text-[#FB4646] hover:bg-[#FB4646]"
              content="Cancel"
              onClick={closeJTHandler}
            />
            <Button className="" content="Join" onClick={joinTeam} />
          </div>
        </div>
      </Dialog>

      {/*********************/}
      {/*********************/}
      <Dialog open={openJTL} onClose={closeJTLHandler}>
        <div className="min-w-[348px] max-w-[348px] p-[28px]">
          <div className="mb-[28px]">
            <Title
              className="text-[18px] font-[400] mb-[16px]"
              content="Register as Teamless"
            />
          </div>

          <label>Select your gender: </label>
          <RadioGroup
            aria-label="gender"
            name="controlled-radio-buttons-group"
            value={state.gender}
            onChange={handleChange("gender")}
          >
            <FormControlLabel value="true" control={<Radio />} label="Male" />
            <FormControlLabel
              value="false"
              control={<Radio />}
              label="Female"
            />
          </RadioGroup>

          <label>Do you mind participating with the other gender? </label>
          <RadioGroup
            aria-label="gender"
            name="controlled-radio-buttons-group"
            value={state.partWithOther}
            onChange={handleChange("partWithOther")}
          >
            <FormControlLabel
              value="true"
              control={<Radio />}
              label="Yes, I mind"
            />
            <FormControlLabel
              value="false"
              control={<Radio />}
              label="No, I don't mind"
            />
          </RadioGroup>

          <br />
          <br />
          <div className=" space-x-[4px] float-right">
            <Button
              className="border-[#FB4646] text-[#FB4646] hover:bg-[#FB4646]"
              content="Cancel"
              onClick={closeJTLHandler}
            />
            <Button className="" content="Register" onClick={joinAsTeamless} />
          </div>
        </div>
      </Dialog>
      {/*********************/}
    </div>
  );
};

export default ContestantContestCard;

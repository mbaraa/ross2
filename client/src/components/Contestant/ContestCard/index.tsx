import { getLocaleTime } from "../../../utils";
import * as React from "react";
import {
  Dialog,
  FormControlLabel,
  Radio,
  RadioGroup,
  TextField,
} from "@mui/material";
import { GoLocation } from "react-icons/go";
import { BiTimeFive } from "react-icons/bi";
import { Button } from "@mui/material";
import Title from "../../Shared/Title";
import Team from "../../../models/Team";
import ContestantRequests from "../../../utils/requests/ContestantRequests";
import JoinRequest from "../../../models/JoinRequest";
import Contestant from "../../../models/Contestant";
import Contest from "../../../models/Contest";
import config from "../../../config";
import ActionChecker from "../../../utils/ActionChecker";
import { checkUserType, UserType } from "../../../models/User";
import { useHistory } from "react-router-dom";

interface Props {
  contest: Contest;
}

const ContestCard = ({ contest }: Props) => {
  const [team] = React.useState<Team>(new Team());

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
    ActionChecker.checkUser(() => {
      setOpenJT(true);
    });
  };

  const closeJTHandler = () => {
    setOpenJT(false);
  };

  const [openCT, setOpenCT] = React.useState(false);
  const openCTHandler = () => {
    ActionChecker.checkUser(() => {
      setOpenCT(true);
    });
  };

  const closeCTHandler = () => {
    setOpenCT(false);
  };

  const [openJTL, setOpenJTL] = React.useState(false);
  const openJTLHandler = () => {
    ActionChecker.checkUser(() => {
      setOpenJTL(true);
    });
  };

  const closeJTLHandler = () => {
    setOpenJTL(false);
  };

  const checkRegisterEnds = (): boolean => {
    const regOver =
      new Date().getTime() > (contest.registration_ends as number);
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
      requested_contest_id: contest.id as number,
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
      setCont(await ContestantRequests.getProfile());
    })();
  }, []);

  const [joinedContest, setJoinedContest] = React.useState(false);

  React.useEffect(() => {
    (async () => {
      setJoinedContest(await ContestantRequests.checkContestJoin(contest));
    })();
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

  const [hasTeam, setHasTeam] = React.useState<boolean>(false);

  React.useEffect(() => {
    setHasTeam(
      contestantProfile !== null && (contestantProfile.team_id as number) > 1
    );
  }, [contestantProfile]);

  const registerInContest = () => {
    (async () => {
      if (
        window.confirm(
          `Are you sure that you want to register in the contest "${contest.name}"`
        )
      ) {
        window.alert(await ContestantRequests.registerInContest(contest));
      }
    })();
  };

  const router = useHistory();

  return (
    <div className="float-left border-[1px] font-Ropa border-ross2 rounded h-auto inline-block w-[300px] mr-[16px] last:mr-0 mb-[16px]">
      <div
        className=" cursor-pointer"
        title="Go to contest's page"
        onClick={() => {
          router.push(`/contest/${contest.id}`, "_blank");
        }}
      >
        <div className="grid grid-cols-2 p-[20px]">
          <div>
            <img
              alt={contest.name}
              src={`${config.backendAddress}${contest.logo_path}`}
              className="rounded-full w-[75px] h-[75px] inline"
            />
          </div>

          <div className="text-left">
            <label className="font-[700] text-[22px] text-ross2">
              {contest.name}
            </label>
            <br />
            <label className="font-[500] text-[15px] text-ross2 ">
              <BiTimeFive className="inline-block" />{" "}
              {getLocaleTime(contest.starts_at)}
            </label>
            <br />
            <label className="font-[500] text-[15px] text-ross2 ">
              <GoLocation size="15px" className="inline-block" />{" "}
              {contest.location}
            </label>
          </div>
        </div>
      </div>

      {joinedContest ? (
        <div
          className="border-t-[1px] border-ross2 py-[12px] text-[13px] font-[600] text-ross2 text-center"
          title="You are already registered in this contest!"
        >
          Already Registerd ✅
        </div>
      ) : hasTeam ? (
        <div
          className="border-t-[1px] border-ross2 py-[12px] text-[13px] font-[600] text-ross2 text-center cursor-pointer hover:bg-ross2 hover:text-white"
          title="Join using the current team"
          onClick={registerInContest}
        >
          Register in Contest
        </div>
      ) : (
        <>
          <div
            className="border-t-[1px] border-ross2 py-[12px] text-[13px] font-[600] text-ross2 text-center cursor-pointer hover:bg-ross2 hover:text-white"
            onClick={openCTHandler}
          >
            Create Team
          </div>
          <div
            className="border-t-[1px] border-ross2 py-[12px] text-[13px] font-[600] text-ross2 text-center cursor-pointer hover:bg-ross2 hover:text-white"
            onClick={openJTHandler}
          >
            Join Team
          </div>
          <div
            className="border-t-[1px] border-ross2 py-[12px] text-[13px] font-[600] text-ross2 text-center cursor-pointer hover:bg-ross2 hover:text-white"
            onClick={openJTLHandler}
          >
            Join as Teamless
          </div>
        </>
      )}

      {/*********************/}

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
              color="error"
              variant="outlined"
              onClick={closeCTHandler}
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Cancel
              </label>
            </Button>
            <Button
              color="info"
              variant="outlined"
              onClick={createTeam}
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Create Team
              </label>
            </Button>
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
              color="error"
              variant="outlined"
              onClick={closeJTHandler}
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Cancel
              </label>
            </Button>
            <Button
              color="info"
              variant="outlined"
              onClick={joinTeam}
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Join Team
              </label>
            </Button>
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
              color="error"
              variant="outlined"
              onClick={closeJTLHandler}
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Cancel
              </label>
            </Button>
            <Button
              color="info"
              variant="outlined"
              onClick={joinAsTeamless}
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Join as Teamless
              </label>
            </Button>
          </div>
        </div>
      </Dialog>
      {/*********************/}
    </div>
  );
};

export default ContestCard;

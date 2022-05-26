import { Button } from "@mui/material";
import * as React from "react";
import { ImCross } from "react-icons/im";
import { BsCheckLg } from "react-icons/bs";
import { useParams, useHistory } from "react-router-dom";
import JoinRequest from "../../../models/JoinRequest";
import Team from "../../../models/Team";
import ContestantRequests from "../../../utils/requests/ContestantRequests";
import BaseError from "../../Shared/Errors/BaseError";
import Title from "../../Shared/Title";
import Contestant from "../../../models/Contestant";

interface LabelProps {
  text: string;
}

const FieldLabel = ({ text }: LabelProps): React.ReactElement => {
  return <label className="font-Ropa text-[18px] text-indigo">{text}</label>;
};

const JoinTeam = (): React.ReactElement => {
  const router = useHistory();
  const { id, join }: any = useParams();
  const [team, setTeam] = React.useState<Team>(new Team());

  const [err, setErr] = React.useState(false);

  // fetch team
  React.useEffect(() => {
    (async () => {
      const team = await ContestantRequests.getTeamByJoinID(join);
      if (team === null) {
        setErr(true);
        return;
      }
      setTeam(team as unknown as Team);
      setErr(false);
    })();
  }, []);

  const joinTeam = async () => {
    const resp = await ContestantRequests.requestJoinTeam({
      requested_team_join_id: join,
      request_message: "",
      requested_contest_id: parseInt(id as string),
    } as JoinRequest);
    if (resp.ok) {
      window.alert(
        "request sent successfully, now wait for the team's creator to accept your request!"
      );
    } else {
      window.alert(await resp.text());
    }
    // window.location.reload();
  };

  const [contestantProfile, setCont] = React.useState(new Contestant());
  React.useEffect(() => {
    (async () => {
      setCont(await ContestantRequests.getProfile());
    })();
  }, []);

  const [hasTeam, setHasTeam] = React.useState<boolean>(false);

  React.useEffect(() => {
    setHasTeam(
      contestantProfile !== null && (contestantProfile.team_id as number) > 1
    );
  }, [contestantProfile]);

  return (
    <div className="absolute left-[50%] translate-x-[-50%] font-Ropa">
      {err ? (
        <BaseError errMsg="Join ID is incorrect!" />
      ) : hasTeam ? (
        <BaseError errMsg="You are already in a Team!" />
      ) : (
        <div>
          <Title
            className=""
            content={`Are you sure that you want to join the team "${team.name}"?`}
          />
          <div className="absolute left-[50%] translate-x-[-50%] font-Ropa">
            <Button
              startIcon={<BsCheckLg size={12} />}
              color="info"
              variant="outlined"
              size="large"
              onClick={joinTeam}
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Yes
              </label>
            </Button>
            &nbsp;
            <Button
              startIcon={<ImCross size={12} />}
              color="error"
              variant="outlined"
              size="large"
              onClick={() => router.push(`/contest/${id}`)}
            >
              <label className="normal-case font-Ropa cursor-pointer">No</label>
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};

export default JoinTeam;

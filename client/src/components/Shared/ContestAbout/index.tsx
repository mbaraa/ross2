import * as React from "react";
import { GiSandsOfTime } from "react-icons/gi";
import { GoClock, GoLocation } from "react-icons/go";
import { MdGroup } from "react-icons/md";
import { BsTextLeft } from "react-icons/bs";
import Contest from "../../../models/Contest";
import { formatDuration, getLocaleTime } from "../../../utils";
import TimerCountdown from "../TimerCountdown";

interface Props {
  contest: Contest;
}

const ContestAbout = ({ contest }: Props): React.ReactElement => {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 p-[25px]">
      <div>
        <div className="text-ross2 font-[25px] font-Ropa pb-[25px]">
          <b>
            <GoClock className="inline-grid" />
            {"  "}Starts At:
          </b>
          <br />
          {getLocaleTime(contest.starts_at)}
        </div>
        <div className="text-ross2 font-[20px] font-Ropa pb-[25px]">
          <b>
            <MdGroup className="inline-grid" />
            {"  "}Members Per Team Limit:
          </b>
          <br />
          Min {contest.participation_conditions.min_team_members}, Max{" "}
          {contest.participation_conditions.max_team_members}
        </div>
        <div className="text-ross2 font-[20px] font-Ropa pb-[25px]">
          <b>
            <GoLocation className="inline-grid" />
            {"  "}Location:
          </b>
          <br />
          {contest.location}
        </div>
      </div>
      <div>
        <div className="text-ross2 font-[25px] font-Ropa pb-[25px]">
          <b>
            <GoClock className="inline-grid" />
            {"  "}Registration Ends In:
          </b>
          <br />
          <TimerCountdown endTimestamp={contest.registration_ends as number} />
        </div>
        <div className="text-ross2 font-[20px] font-Ropa pb-[25px]">
          <b>
            <GiSandsOfTime className="inline-grid" />
            {"  "}Duration:
          </b>
          <br />
          {formatDuration(contest.duration as number)}
        </div>
      </div>
      <div className="lg:w-[150%]">
        <div className="text-ross2 font-[20px] font-Ropa pb-[25px]">
          <b>
            <BsTextLeft className="inline-grid" />
            {"  "}Description:
          </b>
          <br />
          {contest.description}
        </div>
      </div>
    </div>
  );
};

export default ContestAbout;

import * as React from "react";
import { RadialChart } from "react-vis";
import Contest from "../../../models/Contest";
import User, { checkUserType, UserType } from "../../../models/User";
import OrganizerRequests from "../../../utils/requests/OrganizerRequests";
import Title from "../../Shared/Title";

interface Props {
  contest: Contest;
}

const UsersChart = ({ contest }: Props): React.ReactElement => {
  const [participants, setParticipants] = React.useState(new Array<User>());

  React.useEffect(() => {
    (async () => {
      setParticipants(await OrganizerRequests.getParticipantsList(contest));
    })();
  }, []);

  const defineUserYear = (u: User): number => {
    const now = new Date().getFullYear();
    const year = now - parseInt((u.email ?? "1970").substring(0, 4));

    return year;
  };

  interface PiePart {
    angle: number;
    label: string;
    subLabel: string;
  }

  const getnthYear = (n: number, userType: UserType): User[] => {
    return participants.filter(
      (u) => defineUserYear(u) === n && checkUserType(u, userType)
    );
  };

  const getParts = (userType: UserType): PiePart[] => {
    const _1st = getnthYear(1, userType);
    const _2nd = getnthYear(2, userType);
    const _3rd = getnthYear(3, userType);
    const _4th = getnthYear(4, userType);
    const _5th = getnthYear(5, userType);
    const _6th = getnthYear(6, userType);

    let parts = new Array<PiePart>();
    if (_1st.length > 0) {
      parts.push({
        label: "1st Year",
        subLabel: _1st.length.toString(),
        angle: _1st.length,
      });
    }
    if (_2nd.length > 0) {
      parts.push({
        label: "2nd Year",
        subLabel: _2nd.length.toString(),
        angle: _2nd.length,
      });
    }
    if (_3rd.length > 0) {
      parts.push({
        label: "3rd Year",
        subLabel: _3rd.length.toString(),
        angle: _3rd.length,
      });
    }
    if (_4th.length > 0) {
      parts.push({
        label: "4th Year",
        subLabel: _4th.length.toString(),
        angle: _4th.length,
      });
    }
    if (_5th.length > 0) {
      parts.push({
        label: "5th Year",
        subLabel: _5th.length.toString(),
        angle: _5th.length,
      });
    }
    if (_6th.length > 0) {
      parts.push({
        label: "6th Year",
        subLabel: _6th.length.toString(),
        angle: _6th.length,
      });
    }

    return parts;
  };

  const getMaxID = (): string => {
    let max = 0;
    participants.forEach((u) => {
      const uniID = parseInt((u.email ?? "").split("@")[0]);
      if (uniID > max) {
        max = uniID;
      }
    });

    return `${max}`;
    //, ${      (participants.find((u) => u.email?.includes(max.toString())) as User).name}`;
  };

  const getMinID = (): string => {
    let min = 2 << 61;
    participants.forEach((u) => {
      const uniID = parseInt((u.email ?? "").split("@")[0]);
      if (uniID < min) {
        min = uniID;
      }
    });

    return `${min}`;
    //${(participants.find((u) => u.email?.includes(min.toString())) as User).name
    //}`;
  };

  return (
    <>
      <div className="grid grid-cols-1 2xl:grid-cols-2">
        {getParts(UserType.Contestant).length > 0 && (
          <div>
            <Title content="Contestants" className="w-full" />
            <RadialChart
              height={350}
              width={350}
              data={getParts(UserType.Contestant)}
              showLabels
              className="font-Ropa"
            />
          </div>
        )}
        {getParts(UserType.Organizer).length > 0 && (
          <div>
            <Title content="Organizers" className="w-full" />
            <RadialChart
              height={350}
              width={350}
              data={getParts(UserType.Organizer)}
              showLabels
              className="font-Ropa"
            />
          </div>
        )}
      </div>

      <br />

      <div className="py-[20px]">
        <Title
          content={`Max University ID: ${getMaxID()}`}
          className="text-indigo text-[20px]"
        />
        <Title
          content={`Min University ID: ${getMinID()}`}
          className="text-indigo text-[20px]"
        />
      </div>
    </>
  );
};

export default UsersChart;

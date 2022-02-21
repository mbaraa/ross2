import { Button, TextField } from "@mui/material";
import * as React from "react";
import { BiUser } from "react-icons/bi";
import { MdDone } from "react-icons/md";
import Contest from "../../../models/Contest";
import User from "../../../models/User";
import OrganizerRequests from "../../../utils/requests/OrganizerRequests";

interface Props {
  contest: Contest;
}

const AttentanceList = ({ contest }: Props): React.ReactElement => {
  const [participants, setParticipants] = React.useState(new Array<User>());

  React.useEffect(() => {
    (async () => {
      setParticipants(await OrganizerRequests.getParticipantsList(contest));
    })();
  }, [contest]);

  const [searchQuery, setSearchQuery] = React.useState("");

  const filterUsers = (): User[] => {
    return participants.filter((user) => user.email?.includes(searchQuery));
    // .slice(0, participants.length >= 10? 10: participants.length);
  };

  const checkAttended = (user: User) => {
    (async () => {
      await OrganizerRequests.markParticipantAsPresent(user, contest);
      setParticipants(
        participants.splice(
          participants.findIndex((u) => u.email === user.email),
          1
        )
      );
    })();
  };

  return (
    <>
      {participants !== null && participants.length === 0 && (
        <>
          <div className="pb-[20px] px-[14px]">
            <TextField
              className="w-[100%]"
              variant="outlined"
              value={searchQuery}
              onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                setSearchQuery(event.target.value);
              }}
              label={
                <label className="font-Ropa text-[18px] text-indigo">
                  University ID (email)
                </label>
              }
            />
          </div>
          <div className="grid w-full grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
            {filterUsers().map((user) => (
              <>
                <div
                  key={user.id}
                  className="float-left border-[1px] border-indigo rounded h-auto inline-block w-[280px] ml-[16px] mb-[16px] last:mb-0 font-Ropa p-[25px]"
                >
                  <BiUser className="inline" />
                  <label>&nbsp;{user.name}</label>
                  <hr className="pb-[10px] mt-[10px]" />
                  <label>{user.email}</label>
                  <hr className="pb-[10px] mt-[10px]" />
                  <Button
                    endIcon={<MdDone size={12} />}
                    color="info"
                    variant="outlined"
                    size="large"
                    onClick={() => checkAttended(user)}
                  >
                    <label className="normal-case font-Ropa cursor-pointer">
                      Attended
                    </label>
                  </Button>
                </div>
              </>
            ))}
          </div>{" "}
          {/*}*/}
        </>
      )}
    </>
  );
};

export default AttentanceList;

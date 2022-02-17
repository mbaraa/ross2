import { ReactElement } from "react";
import Title from "../Title";
import MenuButton from "../MenuButton";
import { useHistory } from "react-router-dom";
import Contest from "../../models/Contest";
import { getLocaleTime } from "../../utils";
import OrganizerRequests from "../../utils/requests/OrganizerRequests";
import { Button } from "@mui/material";
import { GoPlus } from "react-icons/go";
import { Link } from "react-router-dom";

interface Props {
  contests: any[];
}

const OrganizerContestCards = ({ contests }: Props): ReactElement => {
  const router = useHistory();

  let lables = [
    { id: 1, value: "Name" },
    { id: 2, value: "Date" },
    { id: 3, value: "Location" },
  ];

  let rowDesign =
    "float-left text-ellipsis text-[16px] w-[25%] font-[400] text-ross2 px-[18px] py-[16px] font-Ropa";

  const sendFormNotificationHandler = async (contest: Contest) => {
    await OrganizerRequests.sendContestOverNotifications(contest);
  };

  const downloadCSVHandler = async (contest: Contest) => {
    const parts = await OrganizerRequests.getParticipants(contest);
    const f = document.createElement("a");
    const blob = new Blob(["\ufeff", parts]);
    f.href = URL.createObjectURL(blob);
    f.download = `${contest.name} participants.csv`;
    f.click();
  };

  const deleteHandler = async (contest: Contest) => {
    if (
      window.confirm(
        `Are you sure you want to delete the contest ${contest.name}?`
      )
    ) {
      await OrganizerRequests.deleteContest(contest);
      window.location.reload();
    }
  };

  return (
    <div>
      <div className="grid grid-cols-2 w-full">
        {contests.length > 0 ? (
          <Title className="mb-[12px]" content="Contests"></Title>
        ) : (
          <h1 className="text-ross2 font-bold text-[2rem] font-Ropa">
            No contests are available at this time
          </h1>
        )}
        <div className="absolute items-right pt-[5px] right-[52px]">
          <Link to="/contests/new">
            <Button
              startIcon={<GoPlus size={12} />}
              color="error"
              variant="outlined"
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Create Contest
              </label>
            </Button>
          </Link>
        </div>
      </div>
      {contests.length > 0 && (
        <div className="rounded-[10px] border-[1px] border-ross2">
          <table className="w-[100%]">
            <tbody>
              <tr className="border-b-[1px] border-ross2">
                {lables.map((lable) => {
                  return (
                    <th
                      key={lable.id}
                      className="float-left w-[25%] text-ellipsis text-left text-[18px] text-ross2 px-[20px] py-[16px]"
                    >
                      {lable.value}
                    </th>
                  );
                })}
              </tr>

              {contests.map((contest) => {
                return (
                  <tr
                    key={contest.id}
                    className="border-b-[1px] border-ross2 last:border-b-[0px] font-Ropa"
                  >
                    <td className={rowDesign}>{contest.name}</td>
                    <td className={rowDesign}>
                      {getLocaleTime(contest.starts_at)}
                    </td>
                    <td className={rowDesign}>{contest.location}</td>
                    <td className="float-right text-[13px] px-[20px] py-[16px] font-Ropa">
                      <MenuButton
                        menuItems={[
                          {
                            id: 1,
                            title: "Contest Page",
                            action: () => router.push(`contest/${contest.id}`),
                          },
                          {
                            id: 2,
                            title: "Send Form Notification",
                            action: () => sendFormNotificationHandler(contest),
                          },
                          {
                            id: 3,
                            title: "Download CSV",
                            action: () => downloadCSVHandler(contest),
                          },
                          {
                            id: 4,
                            title: "Delete",
                            action: () => deleteHandler(contest),
                          },
                        ]}
                      />
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default OrganizerContestCards;

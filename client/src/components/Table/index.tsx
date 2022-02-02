import * as React from "react";
import Title from "../Title";
import MenuButton from "../MenuButton";
import { useHistory } from "react-router-dom";
import Contest from "../../models/Contest";
import { getLocaleTime } from "../../utils";
import OrganizerRequests from "../../utils/requests/OrganizerRequests";

interface Props {
  contests: any[];
}

const Table = ({ contests }: Props) => {
  const router = useHistory();

  let lables = [
    { id: 1, value: "Name" },
    { id: 2, value: "Date" },
    { id: 3, value: "Location" },
  ];

  let rowDesign =
    "float-left text-ellipsis text-[13px] w-[25%] font-[400] text-[#425CBA] px-[20px] py-[16px] font-Ropa";

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
      <Title className="mb-[12px]" content="Contests"></Title>
      <div className="rounded-[10px] border-[1px] border-[#425CBA]">
        <table className="w-[100%]">
          <tbody>
            <tr className="border-b-[1px] border-[#425CBA]">
              {lables.map((lable) => {
                return (
                  <th
                    key={lable.id}
                    className="float-left w-[25%] text-ellipsis text-left text-[14px] text-[#425CBA] px-[20px] py-[16px]"
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
                  className="border-b-[1px] border-indigo last:border-b-[0px] font-Ropa"
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
    </div>
  );
};

export default Table;

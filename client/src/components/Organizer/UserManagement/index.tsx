import { Button } from "@mui/material";
import * as React from "react";
import { BsList } from "react-icons/bs";
import Contest from "../../../models/Contest";
import OrganizerRequests from "../../../utils/requests/OrganizerRequests";
import AttentanceList from "../AttendanceList";
import UsersChart from "../UsersChart";

interface Props {
  contest: Contest;
  isDirector: boolean;
  isReceptionist: boolean;
}

const UserManagerment = ({
  contest,
  isDirector,
  isReceptionist,
}: Props): React.ReactElement => {
  const sendFormNotificationHandler = () => {
    (async () => {
      const resp = await OrganizerRequests.sendContestOverNotifications(contest);
      if (!resp.ok) {
        window.alert("Something went wrong, try again later!");
      }
    })();
  };

  const downloadCSVHandler = () => {
    (async () => {
      const parts = await OrganizerRequests.getParticipants(contest);
      const f = document.createElement("a");
      const blob = new Blob(["\ufeff", parts]);
      f.href = URL.createObjectURL(blob);
      f.download = `${contest.name} Tarticipants.csv`;
      f.click();
    })();
  };

  const downloadTeamsCSVHandler = () => {
    (async () => {
      const parts = await OrganizerRequests.getTeamsCSV(contest);
      const f = document.createElement("a");
      const blob = new Blob(["\ufeff", parts]);
      f.href = URL.createObjectURL(blob);
      f.download = `${contest.name} Teams.csv`;
      f.click();
    })();
  };

  return (
    <>
      {/* csv & notifications */}
      {isDirector && (
        <>
          <div className="grid grid-cols-1 sm:grid-cols-3 pb-[20px]">
            <div>
              <Button
                startIcon={<BsList size={12} />}
                color="success"
                variant="outlined"
                size="large"
                onClick={downloadCSVHandler}
              >
                <label className="normal-case font-Ropa cursor-pointer">
                  Export Attendance to CSV
                </label>
              </Button>
            </div>
            <div>
              <Button
                startIcon={<BsList size={12} />}
                color="success"
                variant="outlined"
                size="large"
                onClick={downloadTeamsCSVHandler}
              >
                <label className="normal-case font-Ropa cursor-pointer">
                  Export Teams to CSV
                </label>
              </Button>
            </div>
            <div>
              <Button
                startIcon={<BsList size={12} />}
                color="info"
                variant="outlined"
                size="large"
                onClick={sendFormNotificationHandler}
              >
                <label className="normal-case font-Ropa cursor-pointer">
                  Send Socity Service Form Notifications
                </label>
              </Button>
            </div>
          </div>
          {/* charts */}
          <UsersChart contest={contest} />
        </>
      )}
      {/* attendance */}
      {isReceptionist && (
        <div>
          <AttentanceList contest={contest} />
        </div>
      )}
    </>
  );
};

export default UserManagerment;

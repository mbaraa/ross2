import * as React from "react";
import Notification from "../../models/Notification";
import { Button } from "@mui/material";
import ContestantRequests from "../../utils/requests/ContestantRequests";

interface Props {
  notification: Notification;
}

const NotificationTitle = ({ notification }: Props): React.ReactElement => {
  const isRequest = (): boolean => {
    return (notification !== null &&
      (notification.content as string).length > 0 &&
      notification.content?.startsWith("_REQ")) as boolean;
  };

  const filterContent = (): string => {
    const lastUnderscore = (notification.content as string).lastIndexOf("_");
    const content = notification.content;
    return (
      isRequest() ? (content as string).substring(4, lastUnderscore) : content
    ) as string;
  };

  const acReq = async () => {
    await ContestantRequests.acceptJoinRequest(notification);
  };

  const waReq = async () => {
    await ContestantRequests.rejectJoinRequest(notification);
    window.location.reload();
  };

  return (
    <div className="text-center rounded-[5px] border-ross2 p-[10px] m-[10px] text-ross2 font-Ropa">
      {/* <!-- non-request --> */}
      <div className="content" v-html="filterContent()">
        {" "}
        {filterContent()}
      </div>
      <div dangerouslySetInnerHTML={{ __html: filterContent() }}></div>
      {/* <!-- request --> */}
      <div>
        <Button color="info" variant="outlined" onClick={acReq} size="large">
          <label className="normal-case font-Ropa cursor-pointer">Accept</label>
        </Button>
        &nbsp;
        <Button color="info" variant="outlined" onClick={waReq} size="large">
          <label className="normal-case font-Ropa cursor-pointer">Reject</label>
        </Button>
      </div>
    </div>
  );
};

export default NotificationTitle;

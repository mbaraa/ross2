import * as React from "react";
import Notification from "../../models/Notification";
import Button from "../Button";
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
      <div className="content" v-html="filterContent()"> {filterContent()}</div>
      <div dangerouslySetInnerHTML={{ __html: filterContent()}}>

      </div>
      {/* <!-- request --> */}
      <div>
        <Button
          content="Accept"
          className="border-sky-800 text-sky-800"
          onClick={acReq}
        />
        &nbsp;
        <Button
          content="Reject"
          className="border-red-800 text-red-800 hover:bg-red-800"
          onClick={waReq}
        />
      </div>
    </div>
  );
};

export default NotificationTitle;

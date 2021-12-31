import { Badge, Dialog, IconButton } from "@mui/material";
import * as React from "react";
import NotificationsIcon from "@mui/icons-material/Notifications";
import NotificationRequests from "../../utils/requests/NotificationRequests";
import Title from "../Title";
import Button from "../Button";
import Notification from "../../models/Notification";
import NotificationTitle from "../NotificationTile";

const Notifications = (): React.ReactElement => {
  const [nots, setNots] = React.useState<any[]>([]);
  React.useEffect(() => {
    setNotifications();
  }, []);

  const setNotifications = async () => {
    const n = await NotificationRequests.getNotifications();
    setNots(n);
  };

  const [open, setOpen] = React.useState(false);

  const openHandler = () => {
    setOpen(true);
  };

  const closeHandler = () => {
    setOpen(false);
  };

  return (
    <>
      <IconButton size="large" aria-label="notifications" onClick={openHandler}>
        <Badge badgeContent={nots.length} color="error">
          <NotificationsIcon className="text-ross2" />
        </Badge>
      </IconButton>

      <Dialog open={open} onClose={closeHandler}>
        <div className="min-w-[348px] max-w-[348px] p-[28px]">
          <div className="mb-[28px]">
            <Title
              className="text-[18px] font-[400] mb-[16px]"
              content="Notifications"
            />
          </div>

          {nots.map((n: Notification) => (
            <NotificationTitle key={n.id} notification={n} />
          ))}

          <div className=" space-x-[4px] float-right">
            <Button
              className="border-[#FB4646] text-[#FB4646] hover:bg-[#FB4646]"
              content="Cancel"
              onClick={closeHandler}
            />
          </div>
        </div>
      </Dialog>
    </>
  );
};

export default Notifications;

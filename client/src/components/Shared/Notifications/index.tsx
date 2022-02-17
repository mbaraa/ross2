import { Badge, Dialog, IconButton } from "@mui/material";
import * as React from "react";
import NotificationsIcon from "@mui/icons-material/Notifications";
import NotificationRequests from "../../../utils/requests/NotificationRequests";
import Title from "../Title";
import Notification from "../../../models/Notification";
import NotificationTitle from "../NotificationTile";
import {Button} from "@mui/material";

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
              color="error"
              variant="contained"
              onClick={closeHandler}
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">Cancel</label>
            </Button>
          </div>
        </div>
      </Dialog>
    </>
  );
};

export default Notifications;

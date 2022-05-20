import { Badge, Dialog, IconButton } from "@mui/material";
import * as React from "react";
import NotificationsIcon from "@mui/icons-material/Notifications";
import NotificationRequests from "../../../utils/requests/NotificationRequests";
import Title from "../Title";
import Notification from "../../../models/Notification";
import NotificationTitle from "../NotificationTile";
import { Button } from "@mui/material";
import { MdClose } from "react-icons/md";
import { BsTrash } from "react-icons/bs";

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

  const clearNotifications = () => {
    (async () => {
      if (
        window.confirm("Are you sure you want to clear your notifications?")
      ) {
        const resp = await NotificationRequests.clearNotifications();
        if (!resp.ok) {
          window.alert("Something went wrong, try again later!");
          return;
        }
      }
    })();
    setNots(new Array<Notification>());
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
          <div className="mb-[10px]">
            <Title className="text-[20px] font-[400]" content="Notifications" />
          </div>

          {nots.length > 0 ? (
            nots.map((n: Notification) => (
              <NotificationTitle key={n.id} notification={n} />
            ))
          ) : (
            <div className="text-red-600 font-Ropa text-[20px] pt-0 pb-[20px]">
              No New Notifications!
            </div>
          )}

          <div className=" space-x-[4px] float-right">
            <Button
              startIcon={<BsTrash />}
              color="warning"
              variant="contained"
              onClick={clearNotifications}
              size="large"
              disabled={nots.length === 0}
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Clear
              </label>
            </Button>
            <Button
              startIcon={<MdClose />}
              color="error"
              variant="contained"
              onClick={closeHandler}
              size="large"
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Close
              </label>
            </Button>
          </div>
        </div>
      </Dialog>
    </>
  );
};

export default Notifications;

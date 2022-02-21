import { Button, TextField } from "@mui/material";
import * as React from "react";
import { BiPlus } from "react-icons/bi";
import Organizer from "../../../models/Organizer";
import User from "../../../models/User";
import AdminRequests from "../../../utils/requests/AdminRequests";

const CreateDirector = (): React.ReactElement => {
  const [newDirEmail, setNewDirEmail] = React.useState("");

  const createDirector = () => {
    (async () => {
      const resp = await AdminRequests.createDirector({
        user: {
          email: newDirEmail,
        } as User,
      } as Organizer);
      if (!resp.ok) {
        window.alert(await resp.text());
        return;
      }
      window.alert("director was created successfully!");
      window.location.reload();
    })();
  };
  return (
    <>
    <div className="mr-[20px] inline">
      <TextField
        className=" w-[40%]"
        variant="outlined"
        value={newDirEmail}
        onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
          setNewDirEmail(event.target.value);
        }}
        label={
          <label className="font-Ropa text-[18px] text-red-600">
            User Email
          </label>
        }
      />
    </div>
      <Button
              startIcon={<BiPlus size={12} />}
              color="error"
              variant="outlined"
              size="large"
              className="h-[56px]"
              onClick={createDirector}
            >
              <label className="normal-case font-Ropa cursor-pointer">
                K
              </label>
              </Button>
    </>
  );
};

export default CreateDirector;

import { Button } from "@mui/material";
import * as React from "react";
import { GoPlus } from "react-icons/go";
import CreateDirector from "../components/Admin/CreateDirector";
import DirectorsGrid from "../components/Admin/DirectorsGrid";
import Title from "../components/Shared/Title";
import User, { checkUserType, UserType } from "../models/User";

interface Props {
  user: User;
}

const Admin = ({ user }: Props): React.ReactElement => {
  const isAdmin =
    user !== null &&
    user !== undefined &&
    user.id !== 0 &&
    checkUserType(user, UserType.Admin);

  const [newDir, setNewDir] = React.useState(false);

  return (
    <>
      {!isAdmin ? (
        <>
          <Title content="Ha ha very funny 🙂" className="" />
          <Title
            content="You've had your fun, now get out!"
            className="font-[15px]"
          />
        </>
      ) : (
        <>
          <Button
            startIcon={<GoPlus size={12} />}
            color="error"
            variant="outlined"
            size="large"
            onClick={() => {
              setNewDir(!newDir);
            }}
          >
            <label className="normal-case font-Ropa cursor-pointer">
              Create Director
            </label>
          </Button>
          {newDir && <div className="py-[10px]">
              <CreateDirector/>
          </div>}
          <div className="pt-[10px]">
              <DirectorsGrid/>
          </div>
        </>
      )}
    </>
  );
};

export default Admin;

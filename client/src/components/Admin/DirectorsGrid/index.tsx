import { Button } from "@mui/material";
import * as React from "react";
import { BiTrash } from "react-icons/bi";
import Organizer from "../../../models/Organizer";
import AdminRequests from "../../../utils/requests/AdminRequests";

const DirectorsGrid = (): React.ReactElement => {
  const [dirs, setDirs] = React.useState(new Array<Organizer>());
  React.useEffect(() => {
    (async () => {
      setDirs(await AdminRequests.getDirectors());
    })();
  }, []);

  const deleteDirector = (director: Organizer) => {
    (async () => {
      if (window.confirm("are you sure you want to delete this director?")) {
        await AdminRequests.deleteDirector(director);
        window.location.reload();
      }
    })();
  };

  return (
    <div className="grid w-full grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
      {dirs !== null && dirs.map((dir) => (
        <>
          <div
            className={`float-left border-[1px] border-red-600 rounded h-auto inline-block w-[280px] mr-[16px] mb-[56px] p-[25px] font-Ropa`}
          >
            <div className="p-[10px] pl-0">
              <b>Name: </b>
              {dir.user.name}
            </div>
            <Button
              startIcon={<BiTrash size={12} />}
              color="error"
              variant="outlined"
              size="large"
              onClick={() => deleteDirector(dir)}
            >
              <label className="normal-case font-Ropa cursor-pointer">
                Delete Director
              </label>
            </Button>
          </div>
        </>
      ))}
    </div>
  );
};

export default DirectorsGrid;

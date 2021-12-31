import Button from "../Button";
import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import Title from "../Title";
import Dropdown from "../Dropdown";

interface Props {
  id: number;
}

const ContestGenerateTeams = ({ id }: Props) => {
  const [open, setOpen] = React.useState(false);
  const [selectedType, setSelectedType] = React.useState("ordered");

  let typesOfGenerate = [
    { id: 1, value: "ordered", name: "Numbered" },
    { id: 2, value: "random", name: "Random teams names" },
    { id: 3, value: "given", name: "Given teams names" },
  ];

  const openHandler = () => {
    setOpen(true);
  };

  const closeHandler = () => {
    setOpen(false);
  };

  return (
    <div>
      <Button
        className=""
        content="Generate Teams"
        onClick={() => openHandler()}
      />

      <Dialog open={open} onClose={closeHandler}>
        <div className="min-w-[348px] max-w-[348px] p-[28px]">
          <div className="mb-[28px]">
            <Title
              className="text-[18px] font-[400] mb-[16px]"
              content="Select The Way Of Generate"
            />

            <Dropdown
              lable=""
              value={selectedType}
              setValue={(value: string) => {
                setSelectedType(value);
                console.log(value);
              }}
              items={typesOfGenerate}
            />
          </div>

          <div className=" space-x-[4px] float-right">
            <Button
              className="border-[#FB4646] text-[#FB4646] hover:bg-[#FB4646]"
              content="Cancel"
              onClick={closeHandler}
            />
            <Button className="" content="Generate" onClick={closeHandler} />
          </div>
        </div>
      </Dialog>
    </div>
  );
};

export default ContestGenerateTeams;

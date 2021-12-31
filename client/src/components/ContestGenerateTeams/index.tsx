import Button from "../Button";
import * as React from "react";
import { Dialog, Input, TextField } from "@mui/material";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import Title from "../Title";
import Dropdown from "../Dropdown";
import OrganizerRequests from "../../utils/requests/OrganizerRequests";
import Contest from "../../models/Contest";
import Team from "../../models/Team";
import Contestant from "../../models/Contestant";

interface Props {
  id: number;
}

const ContestGenerateTeams = ({ id }: Props) => {
  const [contest, setContest] = React.useState<Contest>(new Contest());

  React.useEffect(() => {
    (async () => {
      const c = await Contest.getContestFromServer(id);
      setContest(c);
    })();
  }, []);

  const [loading, setLoadin] = React.useState<boolean>(false);

  const [open, setOpen] = React.useState(false);
  const [selectedType, setSelectedType] = React.useState("ordered");

  const [useGivenNames, setUseGivenNames] = React.useState<boolean>(false);

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

  const [noTeamless, setNoTeamless] = React.useState<boolean>(false);
  const [generated, setGenerated] = React.useState<boolean>(false);

  const [namesFile, setNamesFile] = React.useState<any>(undefined);

  // const upload = () => {
  //   if (document != null && document != undefined) {
  //     setNamesFile(document?.getElementById("names").files[0]);
  //     if (namesFile.type != "text/plain") {
  //       window.alert("file must be of text type!");
  //       setNamesFile(undefined);
  //     }
  //   }
  // };
  interface State {
    names: string;
  }
  const [state, setState] = React.useState<any>({
    names: "",
  });

  const readNamesFile = (): string[] => {
    return state.names.replace("\n", "").split(",");
    // return (await namesFile.text()).replace("\n", "").split(",");
  };

  const checkNamesFile = (): boolean => {
    return (
      selectedType != "given" ||
      (selectedType == "given" && namesFile !== undefined)
    );
  };

  const [genTeams, setGenTeams] = React.useState<Team[]>([]);
  const [teamless, setTeamless] = React.useState<Contestant[]>([]);

  const generateTeams = async () => {
    setLoadin(true);

    const [generatedTeams, leftTeamless] =
      await OrganizerRequests.generateTeams(
        contest,
        selectedType,
        readNamesFile()
      );

    if (generateTeams.length == 0 && leftTeamless == null) {
      setNoTeamless(true);
      return;
    }

    setGenerated(true);
    setNoTeamless(false);
    setLoadin(false);
    setGenTeams(generatedTeams);
    setTeamless(leftTeamless);
    setOpen(false);

    // console.log("gen", generatedTeams);
    // console.log("tl", leftTeamless);
  };

  const handleChange =
    (prop: keyof State) => (event: React.ChangeEvent<HTMLInputElement>) => {
      setState({ ...state, [prop]: event.target.value });
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
              content="Select The Way Of Generation"
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

          {selectedType === "given" && (
            <>
              {/* <Input type="file" onClick={upload} id="names"/> */}
              past a list of the wanted teams names{" "}
              <b>separated by a comma(,)</b>
              <br />
              eg: name1,name2,name3...
              <br />
              <TextField
                label="Teams names list"
                required
                value={state.names}
                onChange={handleChange("names")}
              />
              <br />
              <br />
            </>
          )}

          <div className=" space-x-[4px] float-right">
            <Button
              className="border-[#FB4646] text-[#FB4646] hover:bg-[#FB4646]"
              content="Cancel"
              onClick={closeHandler}
            />
            <Button className="" content="Generate" onClick={generateTeams} />
          </div>
        </div>
      </Dialog>
    </div>
  );
};

export default ContestGenerateTeams;

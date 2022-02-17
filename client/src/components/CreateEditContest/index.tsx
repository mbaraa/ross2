import {
  Button,
  TextField,
  FormControlLabel,
  Radio,
  RadioGroup,
} from "@mui/material";
import * as React from "react";
import Contest, { ParticipationConditions } from "../../models/Contest";
import User from "../../models/User";
import { GoPlus } from "react-icons/go";
import { MdSave } from "react-icons/md";
import AdapterDateFns from "@mui/lab/AdapterDateFns";
import LocalizationProvider from "@mui/lab/LocalizationProvider";
import DateTimePicker from "@mui/lab/DateTimePicker";
import ImageUploader from "../ImageUploader";
import config from "../../config";
import OrganizerRequests from "../../utils/requests/OrganizerRequests";

interface LabelProps {
  text: string;
}

const FieldLabel = ({ text }: LabelProps): React.ReactElement => {
  return <label className="font-Ropa text-[18px] text-indigo">{text}</label>;
};

interface Props {
  user: User;
  contest?: Contest;
}

const CreateEditContest = ({ user, contest }: Props): React.ReactElement => {
  const isEdit = contest !== undefined && contest.id !== 0;

  const [isModified, setIsModified] = React.useState(false);

  const [contest2, setContest] = React.useState<Contest>({
    ...contest,
  } as Contest);
  React.useEffect(() => {
    if (contest === undefined) {
      setContest(new Contest());
    }
  }, []);

  const handleChange =
    (prop: keyof Contest) => (event: React.ChangeEvent<HTMLInputElement>) => {
      setContest({ ...contest2, [prop]: event.target.value });
      setIsModified(true);
    };

  const [registrationEnds, setRegistrationEnds] = React.useState<Date>(
    new Date()
  );
  const [startsAt, setStartsAt] = React.useState<Date>(new Date());

  const checkRegAndStartDate = (): boolean => {
    return (
      (contest2.registration_ends as number) < (contest2.starts_at as number)
    );
  };

  const [partConds, setPartConds] = React.useState<ParticipationConditions>(
    JSON.parse(JSON.stringify({ ...contest2.participation_conditions }))
  );
  const handlePartCondsChange =
    (prop: keyof ParticipationConditions) =>
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setPartConds({ ...partConds, [prop]: event.target.value });
      setIsModified(true);
    };

  const [logoFile, setLogoFile] = React.useState<File | null>(null);

  const uploadLogo = async (): Promise<string> => {
    if (!logoFile) {
      window.alert("select a logo file!");
    }

    const formData = new FormData();
    formData.append("file", logoFile as File);

    let res = "";
    await fetch(
      `${config.backendAddress}/organizer/upload-contest-logo-file/`,
      {
        method: "POST",
        mode: "cors",
        headers: {
          // "Content-Type": `multipart/form-data", content type and boundary is calculated by the browser
          Authorization: localStorage.getItem("token") as string,
        },
        body: formData,
      }
    )
      .then((resp) => resp.text())
      .then((resp) => (res = resp as string))
      .catch((err) => {
        res = err.message;
      });

    return res;
  };

  const createContest = () => {
    (async () => {
      contest2.starts_at = startsAt.getTime();
      contest2.registration_ends = registrationEnds.getTime();

      contest2.participation_conditions.min_team_members = Number(
        partConds.min_team_members
      );
      contest2.participation_conditions.max_team_members = Number(
        partConds.max_team_members
      );
      contest2.duration = Number(contest2.duration);
      contest2.teams_hidden = Boolean(contest2.teams_hidden);

      const errMsg = await uploadLogo();
      if (errMsg.length > 0) {
        window.alert(errMsg);
        return;
      }

      if (!checkRegAndStartDate()) {
        window.alert(
          "Woah... 'Start Date' should be after 'End of Registration Date'!"
        );
        // this.contest = new Contest();
        return;
      }

      contest2.logo_path = "/" + logoFile?.name;
      await OrganizerRequests.createContest(contest2);
      window.alert("Contest was created successfully!");

      window.open("/", "_self");
    })();
  };

  const updateContest = () => {
    (async () => {})();
  };

  return (
    <div className="grid md:grid-cols-2 grid-cols-1">
      {/* left side */}
      <div>
        {!isEdit && (
          <h1 className="font-Ropa text-[30px] text-ross2">New Contest</h1>
        )}
        <div className="grid sm:grid-cols-2 grid-cols-1 pt-[20px]">
          {/* inner right side */}
          <div className="pr-[25px]">
            <TextField
              className="w-[100%]"
              variant="outlined"
              value={contest2.name}
              onChange={handleChange("name")}
              label={<FieldLabel text="Contest Title" />}
            />
            <div className="pt-[25px]" />
            <TextField
              className="w-[100%]"
              variant="outlined"
              value={contest2.location}
              onChange={handleChange("location")}
              label={<FieldLabel text="Contest Location" />}
            />
            <div className="pt-[15px]" />
            <label className="font-Ropa text-[20px] text-indigo">
              Teams Visibility
            </label>
            <RadioGroup
              row
              aria-label="gender"
              value={contest2.teams_hidden}
              onChange={handleChange("teams_hidden")}
            >
              <FormControlLabel
                value="false"
                control={<Radio />}
                label={<FieldLabel text="Visible" />}
              />
              <FormControlLabel
                value="true"
                control={<Radio />}
                label={<FieldLabel text="Hidden" />}
              />
            </RadioGroup>
            <div className="pt-[19px]" />

            <TextField
              className="w-[100%]"
              variant="outlined"
              value={partConds.min_team_members}
              onChange={handlePartCondsChange("min_team_members")}
              label={<FieldLabel text="Min Members Per Team" />}
              type="number"
            />

            <div className="pt-[25px]" />

            <TextField
              className="sm:w-[209%] w-[100%] h-[130px] row-span-4 col-span-4"
              multiline
              rows={4}
              value={contest2.description}
              onChange={handleChange("description")}
              label={<FieldLabel text="Description" />}
            />

            <div className="pt-[25px]" />
          </div>

          {/* inner left side */}

          <div className="pr-[25px]">
            <LocalizationProvider dateAdapter={AdapterDateFns}>
              <DateTimePicker
                renderInput={(props) => (
                  <TextField
                    className="w-[100%]"
                    variant="outlined"
                    {...props}
                  />
                )}
                label={<FieldLabel text="Rigistration End Date" />}
                value={registrationEnds}
                onChange={(newValue) => {
                  setRegistrationEnds(newValue as Date);
                }}
              />
            </LocalizationProvider>

            <div className="pt-[25px]" />

            <LocalizationProvider dateAdapter={AdapterDateFns}>
              <DateTimePicker
                renderInput={(props) => (
                  <TextField
                    className="w-[100%]"
                    variant="outlined"
                    {...props}
                  />
                )}
                label={<FieldLabel text="Start Date" />}
                value={startsAt}
                onChange={(newValue) => {
                  setStartsAt(newValue as Date);
                }}
              />
            </LocalizationProvider>

            <div className="pt-[25px]" />

            <TextField
              className="w-[100%]"
              variant="outlined"
              value={contest2.duration}
              onChange={handleChange("duration")}
              label={<FieldLabel text="Duration (in minutes)" />}
              type="number"
            />
            <div className="pt-[25px]" />

            <TextField
              className="w-[100%]"
              variant="outlined"
              value={partConds.max_team_members}
              onChange={handlePartCondsChange("max_team_members")}
              label={<FieldLabel text="Max Members Per Team" />}
              type="number"
            />

            <div className="pt-[25px]" />
          </div>
        </div>
      </div>

      {/* right side */}
      <div className="sm:pt-[60px]">
        <div onClick={() => setIsModified(true)}>
          <ImageUploader
            maxSize={2560}
            imageFile={logoFile}
            setImageFile={setLogoFile}
          />
        </div>
        <Button
          variant="outlined"
          color={isEdit ? "secondary" : "error"}
          className="float-right"
          startIcon={isEdit ? <MdSave size={18} /> : <GoPlus size={18} />}
          disabled={isEdit && !isModified}
          onClick={() => {
            isEdit ? updateContest() : createContest();
          }}
        >
          <label className="normal-case font-Ropa text-[20px] cursor-pointer">
            {isEdit ? "Save" : "Create"}
          </label>
        </Button>
      </div>
    </div>
  );
};

export default CreateEditContest;

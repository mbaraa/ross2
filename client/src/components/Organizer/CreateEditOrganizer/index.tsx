import * as React from "react";
import {
  Button,
  TextField,
  FormControlLabel,
  Radio,
  RadioGroup,
  Checkbox,
} from "@mui/material";
import User from "../../../models/User";
import { GoPlus } from "react-icons/go";
import { MdSave } from "react-icons/md";
import config from "../../../config";
import OrganizerRequests from "../../../utils/requests/OrganizerRequests";
import Contest from "../../../models/Contest";
import Organizer, { OrganizerRole } from "../../../models/Organizer";

interface LabelProps {
  text: string;
}

const FieldLabel = ({ text }: LabelProps): React.ReactElement => {
  return (
    <label className="font-Ropa text-[18px] text-ross2 normal-case">
      {text}
    </label>
  );
};

interface Props {
  user: User;
  contest: Contest;
  organizer?: Organizer;
}

// TODO
// fix this shit :)
const CreateEditOrganizer = ({
  user,
  contest,
  organizer,
}: Props): React.ReactElement => {
  const isEdit = organizer !== undefined && (organizer.user as User).id !== 0;

  const [isModified, setIsModified] = React.useState(false);

  const [organizer2, setOrganizer] = React.useState<Organizer>({
    ...organizer,
  } as Organizer);
  React.useEffect(() => {
    if (organizer === undefined) {
      setOrganizer(new Organizer());
    }
  }, [organizer]);

  const [email, setEmail] = React.useState({ email: "" });
  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setEmail({
      ...email,
      [event.target.name]: event.target.value,
    });
  };

  const [roles, _setRoles] = React.useState([
    { stateName: "isCoreOrganizer", name: "Core Organizer", checked: false },
    { stateName: "isChiefJudge", name: "Chief Judge", checked: false },
    { stateName: "isJudge", name: "Judge", checked: false },
    { stateName: "isTechnical", name: "Technical", checked: false },
    { stateName: "isCoordinator", name: "Coordinator", checked: false },
    { stateName: "isMedia", name: "Media", checked: false },
    { stateName: "isBalloons", name: "Balloons", checked: false },
    { stateName: "isFood", name: "Food", checked: false },
    { stateName: "isReceptionist", name: "Receptionist", checked: false },
  ]);

  const [selectedRoles, setSelectedRoles] = React.useState({
    isCoreOrganizer: false,
    isChiefJudge: false,
    isJudge: false,
    isTechnical: false,
    isCoordinator: false,
    isMedia: false,
    isBalloons: false,
    isFood: false,
    isReceptionist: false,
  });

  const handleRoleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedRoles({
      ...selectedRoles,
      [event.target.name]: event.target.checked,
    });
  };

  const setRoles = () => {
    organizer2.roles = 0;
    if (selectedRoles.isCoreOrganizer) {
      organizer2.roles |= OrganizerRole.CoreOrganizer;
    }
    if (selectedRoles.isChiefJudge) {
      organizer2.roles |= OrganizerRole.ChiefJudge;
    }
    if (selectedRoles.isJudge) {
      organizer2.roles |= OrganizerRole.Judge;
    }
    if (selectedRoles.isTechnical) {
      organizer2.roles |= OrganizerRole.Technical;
    }
    if (selectedRoles.isCoordinator) {
      organizer2.roles |= OrganizerRole.Coordinator;
    }
    if (selectedRoles.isMedia) {
      organizer2.roles |= OrganizerRole.Media;
    }
    if (selectedRoles.isBalloons) {
      organizer2.roles |= OrganizerRole.Balloons;
    }
    if (selectedRoles.isFood) {
      organizer2.roles |= OrganizerRole.Food;
    }
    if (selectedRoles.isReceptionist) {
      organizer2.roles |= OrganizerRole.Receptionist;
    }
    // for (
    //   let role = OrganizerRole.CoreOrganizer;
    //   role <= OrganizerRole.Receptionist;
    //   role <<= 1
    // ) {
    //   if (roles[Math.log2(role) - 1].checked) {
    //     // -1 since the roles array has only 9 elements
    //     organizer2.roles |= role;
    //   }
    // }
  };

  const createOrganizer = () => {
    setRoles();
    organizer2.user.email = email.email;
    organizer2.contests?.push(contest);

    (async () => {
      const resp = await OrganizerRequests.createOrganizer(organizer2);
      if (!resp.ok) {
        window.alert(await resp.text());
        return;
      }
      window.alert("Organizer was created successfully!");
      window.location.reload();
    })();
  };

  return (
    <div>
      {!isEdit && (
        <h1 className="font-Ropa text-[30px] text-ross2">New Organizer</h1>
      )}
      <div className="grid md:grid-cols-2 grid-cols-1">
        {/* left side */}
        <div>
          <div className="text-ross2 text-[20px] font-Ropa my-[10px]">
            Select Organizer's Roles:
          </div>
          <div className="grid grid-cols-2">
            {/* {roles.map(role => <> */}
            <div>
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectedRoles.isCoreOrganizer}
                    onChange={handleRoleChange}
                    inputProps={{ "aria-label": "controlled" }}
                    color="secondary"
                    name={roles[0].stateName}
                  />
                }
                label={
                  <label className="text-ross2 text-[18px] font-Ropa">
                    {roles[0].name}
                  </label>
                }
              />
              <br />
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectedRoles.isChiefJudge}
                    onChange={handleRoleChange}
                    inputProps={{ "aria-label": "controlled" }}
                    color="secondary"
                    name={roles[1].stateName}
                  />
                }
                label={
                  <label className="text-ross2 text-[18px] font-Ropa">
                    {roles[1].name}
                  </label>
                }
              />
              <br />
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectedRoles.isJudge}
                    onChange={handleRoleChange}
                    inputProps={{ "aria-label": "controlled" }}
                    color="secondary"
                    name={roles[2].stateName}
                  />
                }
                label={
                  <label className="text-ross2 text-[18px] font-Ropa">
                    {roles[2].name}
                  </label>
                }
              />
              <br />
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectedRoles.isTechnical}
                    onChange={handleRoleChange}
                    inputProps={{ "aria-label": "controlled" }}
                    color="secondary"
                    name={roles[3].stateName}
                  />
                }
                label={
                  <label className="text-ross2 text-[18px] font-Ropa">
                    {roles[3].name}
                  </label>
                }
              />
              <br />
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectedRoles.isCoordinator}
                    onChange={handleRoleChange}
                    inputProps={{ "aria-label": "controlled" }}
                    color="secondary"
                    name={roles[4].stateName}
                  />
                }
                label={
                  <label className="text-ross2 text-[18px] font-Ropa">
                    {roles[4].name}
                  </label>
                }
              />
            </div>
            <div>
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectedRoles.isMedia}
                    onChange={handleRoleChange}
                    inputProps={{ "aria-label": "controlled" }}
                    color="secondary"
                    name={roles[5].stateName}
                  />
                }
                label={
                  <label className="text-ross2 text-[18px] font-Ropa">
                    {roles[5].name}
                  </label>
                }
              />
              <br />
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectedRoles.isBalloons}
                    onChange={handleRoleChange}
                    inputProps={{ "aria-label": "controlled" }}
                    color="secondary"
                    name={roles[6].stateName}
                  />
                }
                label={
                  <label className="text-ross2 text-[18px] font-Ropa">
                    {roles[6].name}
                  </label>
                }
              />
              <br />
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectedRoles.isFood}
                    onChange={handleRoleChange}
                    inputProps={{ "aria-label": "controlled" }}
                    color="secondary"
                    name={roles[7].stateName}
                  />
                }
                label={
                  <label className="text-ross2 text-[18px] font-Ropa">
                    {roles[7].name}
                  </label>
                }
              />
              <br />
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectedRoles.isReceptionist}
                    onChange={handleRoleChange}
                    inputProps={{ "aria-label": "controlled" }}
                    color="secondary"
                    name={roles[8].stateName}
                  />
                }
                label={
                  <label className="text-ross2 text-[18px] font-Ropa">
                    {roles[8].name}
                  </label>
                }
              />
            </div>
          </div>
        </div>

        {/* right side */}
        <div>
          <TextField
            className="w-[100%]"
            variant="outlined"
            value={email.email}
            onChange={handleChange}
            label={<FieldLabel text="Email" />}
            name="email"
            type="email"
          />
          <br />
          <br />
          <Button
            variant="outlined"
            color="secondary"
            className="w-full"
            onClick={createOrganizer}
          >
            <FieldLabel text="Create" />
          </Button>
        </div>
      </div>
    </div>
  );
};

export default CreateEditOrganizer;

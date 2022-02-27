import * as React from "react";
import { Button, TextField, FormControlLabel, Checkbox } from "@mui/material";
import User from "../../../models/User";
import { GoPlus } from "react-icons/go";
import { MdSave } from "react-icons/md";
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

const CreateEditOrganizer = ({
  user,
  contest,
  organizer,
}: Props): React.ReactElement => {
  const isEdit = organizer !== undefined && (organizer.user as User).id !== 0;

  const [errMsg, _setErrMsg] = React.useState("");

  const setErrMsg = (msg: string) => {
    _setErrMsg(msg);
    setTimeout(() => {
      _setErrMsg("");
    }, 10000);
  };

  const [isModified, setIsModified] = React.useState(false);

  const [organizer2, setOrganizer] = React.useState<Organizer>({
    ...organizer,
  } as Organizer);

  React.useEffect(() => {
    if (organizer === undefined) {
      setOrganizer(new Organizer());
    }
  }, [organizer]);

  const [email, setEmail] = React.useState({ email: organizer2.user.email });
  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setEmail({
      ...email,
      [event.target.name]: event.target.value,
    });
  };

  const [roles, _setRoles] = React.useState([
    { i: 0, name: "Core Organizer", checked: false },
    { i: 1, name: "Chief Judge", checked: false },
    { i: 2, name: "Judge", checked: false },
    { i: 3, name: "Technical", checked: false },
    { i: 4, name: "Coordinator", checked: false },
    { i: 5, name: "Media", checked: false },
    { i: 6, name: "Balloons", checked: false },
    { i: 7, name: "Food", checked: false },
    { i: 8, name: "Receptionist", checked: false },
  ]);

  const setRoles = () => {
    organizer2.roles = 0;

    for (
      let role = OrganizerRole.CoreOrganizer;
      role <= OrganizerRole.Receptionist;
      role <<= 1
    ) {
      if (roles[Math.log2(role) - 1].checked) {
        // -1 since the roles array has only 9 elements
        (organizer2.roles as number) |= role;
      }
    }
  };

  React.useEffect(() => {
    if (isEdit) {
      for (
        let role = OrganizerRole.CoreOrganizer;
        role <= OrganizerRole.Receptionist;
        role <<= 1
      ) {
        if ((role & (organizer2.roles as number)) !== 0) {
          roles[Math.log2(role) - 1].checked = true;
          _setRoles(roles.flat());
        } else {
          roles[Math.log2(role) - 1].checked = false;
        }
      }
    }
  }, []);

  const checkRoles = (): boolean => (organizer2.roles as number) !== 0;

  const createOrganizer = () => {
    setRoles();
    if (!checkRoles()) {
      setErrMsg("Select at least one role for the organizer!");
      return;
    }
    organizer2.user.email = email.email;
    organizer2.contests?.push(contest);

    (async () => {
      const resp = await OrganizerRequests.createOrganizer(
        organizer2,
        contest,
        organizer2.roles as number
      );
      if (!resp.ok) {
        setErrMsg(await resp.text());
        return;
      }
      window.alert("Organizer was created successfully!");
      window.location.reload();
    })();
  };

  const updateOrganizer = () => {
    setRoles();
    if (!checkRoles()) {
      setErrMsg("Select at least one role for the organizer!");
      return;
    }
    organizer2.user.email = email.email;

    (async () => {
      const resp = await OrganizerRequests.updateOrganizer(
        organizer2,
        contest,
        organizer?.roles as number
      );
      if (!resp.ok) {
        setErrMsg(await resp.text());
        return;
      }
      window.alert("Organizer was updated successfully!");
    })();
  };

  return (
    <div>
      {!isEdit ? (
        <h1 className="font-Ropa text-[30px] text-ross2">New Organizer</h1>
      ) : (
        <h1 className="font-Ropa text-[30px] text-ross2">
          {organizer2.user.name}
        </h1>
      )}
      <div className="grid md:grid-cols-2 grid-cols-1">
        {/* left side */}
        <div>
          <div className="text-ross2 text-[20px] font-Ropa my-[10px]">
            Select Organizer's Roles:
          </div>
          <div className="grid grid-cols-2">
            {roles.map((role) => (
              <div key={role.i}>
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={roles[role.i].checked}
                      onChange={(
                        event: React.ChangeEvent<HTMLInputElement>
                      ) => {
                        roles[role.i].checked = Boolean(event.target.checked);
                        _setRoles(roles.flat());
                        setIsModified(true);
                      }}
                      inputProps={{ "aria-label": "controlled" }}
                      color="secondary"
                    />
                  }
                  label={
                    <label className="text-ross2 text-[18px] font-Ropa">
                      {roles[role.i].name}
                    </label>
                  }
                />
                <br />
              </div>
            ))}
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
            disabled={isEdit}
          />

          <div className="text-[#d63333] text-[20px] font-Ropa py-[10px]">
            {errMsg}
          </div>

          <Button
            variant="outlined"
            color={isEdit ? "secondary" : "error"}
            className="w-full"
            startIcon={isEdit ? <MdSave size={18} /> : <GoPlus size={18} />}
            disabled={isEdit && !isModified}
            onClick={() => {
              isEdit ? updateOrganizer() : createOrganizer();
            }}
          >
            <label className="normal-case font-Ropa text-[20px] cursor-pointer">
              {isEdit ? "Save" : "Create"}
            </label>
          </Button>
        </div>
      </div>
    </div>
  );
};

export default CreateEditOrganizer;

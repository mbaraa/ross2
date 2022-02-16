import { ReactElement } from "react";
import Contest from "../../models/Contest";
import User from "../../models/User";

interface Props {
  user: User;
  contest?: Contest;
}

const CreateEditContest = ({user, contest}: Props): ReactElement => {
  return <>
    {contest !== undefined? "edit": "create"}
  </>;
};

export default CreateEditContest;

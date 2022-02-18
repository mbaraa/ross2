import { useParams } from "react-router-dom";
import { useEffect, useState, ReactElement } from "react";
import { default as Contest2 } from "../models/Contest";
import Title from "../components/Shared/Title";
import * as React from "react";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import ContestGenerateTeams from "../components/Organizer/GenerateTeams";
import ContestManageTeams from "../components/Organizer/ManageTeams";
import User, { checkUserType, UserType } from "../models/User";
import Organizer, { checkOrgType, OrganizerRole } from "../models/Organizer";
import OrganizerRequests from "../utils/requests/OrganizerRequests";
import CreateEditContest from "../components/Organizer/CreateEditContest";
import OrganizersGrid from "../components/Organizer/OrganizersGrid";
import CreateEditOrganizer from "../components/Organizer/CreateEditOrganizer";

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps): ReactElement {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
}

interface LabelProps {
  text: string;
}

const TabLabel = ({ text }: LabelProps): ReactElement => {
  return (
    <label className="font-Ropa text-indigo text-[17px] normal-case">
      {text}
    </label>
  );
};

interface Props {
  user: User;
}

const Contest = ({ user }: Props): ReactElement => {
  const [value, setValue] = React.useState(0);

  const handleChange = (event: React.SyntheticEvent, newValue: number) => {
    setValue(newValue);
  };

  const [contest, setContest] = useState<Contest2>(new Contest2());

  const { id }: any = useParams();

  useEffect(() => {
    if (id !== undefined) {
      (async () => {
        setContest(await Contest2.getContestFromServer(parseInt(id as string)));
      })();
    }
  }, [id]);

  const [org, setOrg] = useState<Organizer>(new Organizer());

  useEffect(() => {
    (async () => {
      setOrg(await OrganizerRequests.getProfile());
    })();
  }, [user]);

  const defUser = user !== null && user !== undefined && user.id !== 0;

  const isAdmin = defUser && checkUserType(user, UserType.Admin);

  const isDirector =
    defUser &&
    checkUserType(user, UserType.Organizer) &&
    checkOrgType(org, OrganizerRole.Director);

  const isCoreOrg =
    defUser &&
    checkUserType(user, UserType.Organizer) &&
    checkOrgType(org, OrganizerRole.CoreOrganizer);

  const isReseptionist =
    defUser &&
    checkUserType(user, UserType.Organizer) &&
    checkOrgType(org, OrganizerRole.Receptionist);

  if (contest.id) {
    return (
      <div className="font-Ropa">
        <Title content={contest.name} className="pl-[20px]" />

        <Box className="font-Ropa" sx={{ width: "100%" }}>
          <Box
            className="border-indigo border-opacity-20"
            sx={{ borderBottom: 1, borderColor: "divider" }}
          >
            <Tabs
              value={value}
              onChange={handleChange}
              className="border-indigo"
              aria-label="basic tabs example"
            >
              <Tab label={<TabLabel text="About" />} value={0} />
              {isDirector && (
                <Tab label={<TabLabel text="Generate Posts" />} value={1} />
              )}
              {isDirector && (
                <Tab label={<TabLabel text="Generate Teams" />} value={2} />
              )}
              {(isDirector || isCoreOrg) && (
                <Tab label={<TabLabel text="Manage Teams" />} value={3} />
              )}
              {(isDirector || isAdmin) && (
                <Tab label={<TabLabel text="Manage Organizers" />} value={4} />
              )}
              {(isDirector || isAdmin) && (
                <Tab label={<TabLabel text="Create Organizer" />} value={5} />
              )}
              {isReseptionist && (
                <Tab label={<TabLabel text="Attendance" />} value={6} />
              )}
              {(isDirector || isCoreOrg) && (
                <Tab label={<TabLabel text="Edit" />} value={7} />
              )}
            </Tabs>
          </Box>

          <TabPanel value={value} index={0}>
            Building...
          </TabPanel>
          {isDirector && (
            <TabPanel value={value} index={1}>
              Building...
            </TabPanel>
          )}
          {isDirector && (
            <TabPanel value={value} index={2}>
              <ContestGenerateTeams id={contest.id} />
            </TabPanel>
          )}
          {(isDirector || isCoreOrg) && (
            <TabPanel value={value} index={3}>
              <ContestManageTeams
                teams={contest.teams}
                teamless={contest.teamless_contestants}
              />
            </TabPanel>
          )}
          {(isDirector || isAdmin) && (
            <TabPanel value={value} index={4}>
              <OrganizersGrid />
            </TabPanel>
          )}
          {(isDirector || isAdmin) && (
            <TabPanel value={value} index={5}>
              <CreateEditOrganizer user={user} contest={contest}/>
            </TabPanel>
          )}
          {isReseptionist && (
            <TabPanel value={value} index={6}>
              Building...
            </TabPanel>
          )}
          {(isDirector || isCoreOrg) && (
            <TabPanel value={value} index={7}>
              <CreateEditContest user={user} contest={contest} />
            </TabPanel>
          )}
        </Box>
      </div>
    );
  }

  return <div>Lodding..</div>;
};

export default Contest;

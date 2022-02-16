import { useParams } from "react-router-dom";
import { useEffect, useState, ReactElement } from "react";
import { default as Contest2 } from "../models/Contest";
import Title from "../components/Title";
import * as React from "react";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import ContestGenerateTeams from "../components/ContestGenerateTeams";
import ContestManageTeams from "../components/ContestManageTeams";
import User, { checkUserType, UserType } from "../models/User";
import Organizer, { checkOrgType, OrganizerRole } from "../models/Organizer";
import OrganizerRequests from "../utils/requests/OrganizerRequests";

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

function a11yProps(index: number): any {
  return {
    id: `simple-tab-${index}`,
    "aria-controls": `simple-tabpanel-${index}`,
  };
}

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

        <Box className="font-[Poppins]" sx={{ width: "100%" }}>
          <Box
            className="border-[#425CBA] border-opacity-20"
            sx={{ borderBottom: 1, borderColor: "divider" }}
          >
            <Tabs
              value={value}
              onChange={handleChange}
              aria-label="basic tabs example"
            >
              <Tab label="About" {...a11yProps(0)} />
              {isDirector && <Tab label="Generate Posts" {...a11yProps(1)} />}
              {isDirector && <Tab label="Generate Teams" {...a11yProps(2)} />}
              {(isDirector || isCoreOrg) && <Tab label="Manage Teams" {...a11yProps(3)} />}
              {(isDirector || isAdmin) && (
                <Tab label="Manage Organizers" {...a11yProps(4)} />
              )}
              {isReseptionist && <Tab label="Attendance" {...a11yProps(5)} />}
              {(isDirector || isCoreOrg) && <Tab label="Edit" {...a11yProps(6)} />}
            </Tabs>
          </Box>

          <TabPanel value={value} index={0}>
            Building...
          </TabPanel>
          <TabPanel value={value} index={1}>
            Building...
          </TabPanel>
          <TabPanel value={value} index={2}>
            <ContestGenerateTeams id={contest.id} />
          </TabPanel>
          <TabPanel value={value} index={3}>
            <ContestManageTeams
              teams={contest.teams}
              teamless={contest.teamless_contestants}
            />
          </TabPanel>
          <TabPanel value={value} index={4}>
            Bulding...
          </TabPanel>
          <TabPanel value={value} index={5}>
            Bulding...
          </TabPanel>
          <TabPanel value={value} index={6}>
            Bulding...
          </TabPanel>
        </Box>
      </div>
    );
  }

  return <div>Lodding..</div>;
};

export default Contest;

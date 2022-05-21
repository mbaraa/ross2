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
import Organizer, { OrganizerRole } from "../models/Organizer";
import OrganizerRequests from "../utils/requests/OrganizerRequests";
import CreateEditContest from "../components/Organizer/CreateEditContest";
import OrganizersGrid from "../components/Organizer/OrganizersGrid";
import GeneratePosts from "../components/Organizer/GeneratePosts";
import ContestAbout from "../components/Shared/ContestAbout";
import UserManagerment from "../components/Organizer/UserManagement";
import {
  Switch,
  Route,
  useHistory,
  useParams,
  useLocation,
} from "react-router-dom";
import ContestSupport from "../components/Contestant/ContestSupport";
import YouCantDoThat from "../components/Shared/Errors/YouCantDoThat";
import NotFound from "../components/Shared/Errors/NotFound";

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
  const router = useHistory();
  const { pathname } = useLocation();
  const getPath = (): string => {
    return pathname.substring(pathname.lastIndexOf("/") + 1);
  };

  const [value, setValue] = React.useState(0);
  React.useEffect(() => {
    let index = 0;
    switch (getPath()) {
      case "about":
        index = 0;
        break;
      case "generate-posts":
        index = 1;
        break;
      case "generate-teams":
        index = 2;
        break;
      case "manage-teams":
        (async () => {
          setContest(
            await OrganizerRequests.getContest(parseInt(id as string))
          );
        })();
        index = 3;
        break;
      case "manage-organizers":
        index = 4;
        break;
      case "manage-users":
        index = 5;
        break;
      case "edit":
        index = 6;
        break;
      case "support":
        index = 7;
        break;
    }
    setValue(index);
  }, []);

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

  const [isDirector, setIsDirector] = React.useState(false);
  React.useEffect(() => {
    (async () => {
      setIsDirector(
        await OrganizerRequests.checkOrgRole(
          parseInt(id as string),
          org.id,
          OrganizerRole.Director
        )
      );
    })();
  }, [org]);

  const [isCoreOrg, setIsCoreOrg] = React.useState(false);
  React.useEffect(() => {
    (async () => {
      setIsCoreOrg(
        await OrganizerRequests.checkOrgRole(
          parseInt(id as string),
          org.id,
          OrganizerRole.CoreOrganizer
        )
      );
    })();
  }, [org]);

  const [isReceptionist, setIsReceptionist] = React.useState(false);
  React.useEffect(() => {
    (async () => {
      setIsReceptionist(
        await OrganizerRequests.checkOrgRole(
          parseInt(id as string),
          org.id,
          OrganizerRole.Receptionist
        )
      );
    })();
  }, [org]);

  React.useEffect(() => {
    if (isDirector) {
      (async () => {
        setContest(await OrganizerRequests.getContest(parseInt(id as string)));
      })();
    }
  }, [isDirector]);

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
              variant="scrollable"
              scrollButtons="auto"
              allowScrollButtonsMobile
              onChange={handleChange}
              className="border-indigo"
              aria-label="scrollable auto tabs example"
            >
              <Tab
                label={<TabLabel text="About" />}
                value={0}
                onClick={() => router.push(`/contest/${contest.id}/about`)}
              />
              {isDirector && (
                <Tab
                  label={<TabLabel text="Generate Posts" />}
                  value={1}
                  onClick={() =>
                    router.push(`/contest/${contest.id}/generate-posts`)
                  }
                />
              )}
              {isDirector && (
                <Tab
                  label={<TabLabel text="Generate Teams" />}
                  value={2}
                  onClick={() =>
                    router.push(`/contest/${contest.id}/generate-teams`)
                  }
                />
              )}
              {(isDirector || isCoreOrg) && (
                <Tab
                  label={<TabLabel text="Manage Teams" />}
                  value={3}
                  onClick={() => {
                    router.push(`/contest/${contest.id}/manage-teams`);
                  }}
                />
              )}
              {isDirector && (
                <Tab
                  label={<TabLabel text="Manage Organizers" />}
                  value={4}
                  onClick={() =>
                    router.push(`/contest/${contest.id}/manage-organizers`)
                  }
                />
              )}
              {(isDirector || isReceptionist) && (
                <Tab
                  label={<TabLabel text="Attendance & Other User Management" />}
                  wrapped
                  value={5}
                  onClick={() =>
                    router.push(`/contest/${contest.id}/manage-users`)
                  }
                />
              )}
              {(isDirector || isCoreOrg) && (
                <Tab
                  label={<TabLabel text="Edit" />}
                  value={6}
                  onClick={() => router.push(`/contest/${contest.id}/edit`)}
                />
              )}
              {!isDirector && (
                <Tab
                  label={<TabLabel text={"Support"} />}
                  value={7}
                  onClick={() => router.push(`/contest/${contest.id}/support`)}
                />
              )}
            </Tabs>
          </Box>

          <Switch>
            <Route path="/contest/:id/about">
              <TabPanel value={value} index={0}>
                <ContestAbout contest={contest} />
              </TabPanel>
            </Route>

            <Route path="/contest/:id/generate-posts">
              {isDirector ? (
                <TabPanel value={value} index={1}>
                  <GeneratePosts contest={contest} />
                </TabPanel>
              ) : (
                <YouCantDoThat />
              )}
            </Route>

            <Route path="/contest/:id/generate-teams">
              {isDirector ? (
                <TabPanel value={value} index={2}>
                  <ContestGenerateTeams id={contest.id} />
                </TabPanel>
              ) : (
                <YouCantDoThat />
              )}
            </Route>

            <Route path="/contest/:id/manage-teams">
              {isDirector || isCoreOrg ? (
                <TabPanel value={value} index={3}>
                  {contest.teams === null ? (
                    <Title className="" content="Loading..." />
                  ) : (
                    <ContestManageTeams
                      teams={contest.teams}
                      teamless={contest.teamless_contestants}
                      showGender={false}
                      contest={contest}
                      updateTeams={() => {
                        contest.teams = contest.teams.flat();
                        setContest({ ...contest });
                      }}
                    />
                  )}
                </TabPanel>
              ) : (
                <YouCantDoThat />
              )}
            </Route>

            <Route path="/contest/:id/manage-organizers">
              {isDirector ? (
                <TabPanel value={value} index={4}>
                  <OrganizersGrid user={user} contest={contest} />
                </TabPanel>
              ) : (
                <YouCantDoThat />
              )}
            </Route>

            <Route path="/contest/:id/manage-users">
              {isDirector || isReceptionist ? (
                <TabPanel value={value} index={5}>
                  <UserManagerment
                    contest={contest}
                    isDirector={isDirector}
                    isReceptionist={isReceptionist}
                  />
                </TabPanel>
              ) : (
                <YouCantDoThat />
              )}
            </Route>

            <Route path="/contest/:id/edit">
              {isDirector || isCoreOrg ? (
                <TabPanel value={value} index={6}>
                  <CreateEditContest user={user} contest={contest} />
                </TabPanel>
              ) : (
                <YouCantDoThat />
              )}
            </Route>

            <Route path="/contest/:id/support">
              <ContestSupport
                orgs={contest.organizers as Organizer[]}
                contest={contest}
              />
            </Route>

            <Route component={NotFound} />
          </Switch>
        </Box>
      </div>
    );
  }

  return <Title content="Loading..." className="" />;
};

export default Contest;

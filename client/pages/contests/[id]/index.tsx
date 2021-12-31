import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { default as Contest2 } from "../../../src/models/Contest";
import Title from "../../../src/components/Title";
import * as React from "react";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import ContestGenerateTeams from "../../../src/components/ContestGenerateTeams";

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
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

function a11yProps(index: number) {
  return {
    id: `simple-tab-${index}`,
    "aria-controls": `simple-tabpanel-${index}`,
  };
}

const Contest = () => {
  const router = useRouter();

  const [value, setValue] = React.useState(0);

  const handleChange = (event: React.SyntheticEvent, newValue: number) => {
    setValue(newValue);
  };

  const [contest, setContest] = useState<Contest2>(new Contest2());

  useEffect(() => {
    if (router.query["id"] !== undefined) {
      setCont();
    }
  }, [router.query["id"]]);

  const setCont = async () => {
    const c = await Contest2.getContestFromServer(
      parseInt(router.query["id"] as string)
    );

    console.log(router.query["id"]);

    setContest(c);
  };

  if (contest.id) {
    return (
      <div>
        <Title content={contest.name} />

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
              <Tab className="capitalize " label="About" {...a11yProps(0)} />
              <Tab
                className="capitalize"
                label="Generate Teams"
                {...a11yProps(1)}
              />
              <Tab
                className="capitalize"
                label="Manage Teams"
                {...a11yProps(2)}
              />
            </Tabs>
          </Box>
          <TabPanel value={value} index={0}>
            Under Building
          </TabPanel>
          <TabPanel value={value} index={1}>
            <ContestGenerateTeams id={contest.id} />
          </TabPanel>
          <TabPanel value={value} index={2}>
            Item Three
          </TabPanel>
        </Box>
      </div>
    );
  }

  return <div>Lodding..</div>;
};

export default Contest;

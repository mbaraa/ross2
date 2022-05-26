import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Toolbar,
  Typography,
} from "@mui/material";
import * as React from "react";

const About = (): React.ReactElement => {
  const [techStack] = React.useState([
    {
      id: 0,
      name: "Go",
      purpose: "Backend Programming Language",
      version: "1.18",
      link: "https://golang.org",
    },
    {
      id: 1,
      name: "TypeScript",
      purpose: "Frontend Programming Language",
      version: "4.5",
      link: "https://www.typescriptlang.org",
    },
    {
      id: 2,
      name: "React.js",
      purpose: "UI Library",
      version: "17.0",
      link: "https://reactjs.org",
    },
    {
      id: 3,
      name: "MariaDB",
      purpose: "Database",
      version: "10.6",
      link: "https://mariadb.org/",
    },
    {
      id: 4,
      name: "GORM",
      purpose: "ORM Database",
      version: "1.21",
      link: "https://gorm.io",
    },
    {
      id: 5,
      name: "MUI",
      purpose: "UI Material Components",
      version: "5.2",
      link: "https://mui.com",
    },
    {
      id: 6,
      name: "Tailwind CSS",
      purpose: "CSS Framework",
      version: "3.0",
      link: "https://tailwindcss.com",
    },
    {
      id: 7,
      name: "Cairo Graphics",
      purpose: "Image Manipulation Library",
      version: "1.17",
      link: "https://www.cairographics.org/",
    },
  ]);

  const [titles] = React.useState([
    { id: 0, name: "Name" },
    { id: 1, name: "Purpose" },
    { id: 2, name: "Version" },
  ]);

  return (
    <div className="absolute left-[50%] translate-x-[-50%] font-Ropa w-[90%]">
      <div className="grid md:grid-cols-2 grid-cols-1">
        <div className="pt-[20px]">
          <a href="https://github.com/mbaraa/ross2" target="_blank">
            <img
              src="/logo512.png"
              className="w-[200px] h-[200px] border-[2px] border-grey-200 rounded-[100%] hover:opacity-[80%] relative left-[50%] translate-x-[-50%]"
            />
          </a>
          <h1 className="text-center mt-[20px] text-[40px] text-ross2 font-black">
            Ross 2
          </h1>
          <h1 className="text-center text-[30px] text-ross2">
            Contest Management Thingy
          </h1>

          <h2 className="text-center text-[20px] ">
            Support this project by staring it on GitHub 🥰
            <br />
            <iframe
              src="https://ghbtns.com/github-btn.html?user=mbaraa&repo=ross2&type=star&count=true"
              scrolling="0"
              title="GitHub"
              className="relative left-[50%] translate-x-[-50%] w-[75px] h-[25px]"
            ></iframe>
          </h2>
        </div>

        <div>
          <h1 className="text-[22px] text-[#343434] m-[10px] mx-0">
            Technologies Used:
          </h1>
          <TableContainer className="w-[80vw] mb-[30px]">
            <Table
              stickyHeader={true}
              className="border-[1px] border-grey-[300]"
            >
              <TableHead>
                <TableRow>
                  {titles.map((col) => (
                    <TableCell key={col.id}>{col.name}</TableCell>
                  ))}
                </TableRow>
              </TableHead>

              <TableBody>
                {techStack.map((row) => (
                  <TableRow key={row.id} hover={true}>
                    <TableCell>
                      <a href={row.link} target="_blank">
                        {row.name}
                      </a>
                    </TableCell>
                    <TableCell>{row.purpose}</TableCell>
                    <TableCell>{row.version}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </div>
      </div>
    </div>
  );
};

export default About;

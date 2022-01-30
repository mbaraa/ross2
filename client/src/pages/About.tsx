import * as React from "react";

const About = (): React.ReactElement => {
  let list = [
    {
      id: 1,
      text: "Baraa Al-Masri - Author & Backend & Frontend Developer",
      link: "https://github.com/mbaraa",
    },
    {
      id: 2,
      text: "Akram Khalil - UI Designer & Frontend Developer",
      link: "https://github.com/2kram",
    },
    {
      id: 3,
      text: "Go - Backend",
      link: "https://golang.org/",
    },
    {
      id: 4,
      text: "ReactJS - Frontend",
      link: "https://reactjs.org",
    },
    {
      id: 5,
      text: "GORM DB - ORM",
      link: "https://gorm.io/",
    },

    {
      id: 6,
      text: "GoJWT - jwt validator",
      link: "https://github.com/golang-jwt/jwt",
    },
    {
      id: 7,
      text: "Google UUID - uuid generator",
      link: "https://github.com/google/uuid",
    },
    {
      id: 8,
      text: "MUI - Material Components",
      link: "https://mui.com",
    },
    {
      id: 9,
      text: "Tailwind - CSS Classes",
      link: "https://tailwindcss.com",
    },
  ];

  return (
    <div>
      <ul className="space-y-[12px] list-disc">
        {list.map((i) => {
          return (
            <li className="font-[12px] text-[#000] font-[400]" key={i.id}>
              <div className="mr-[8px] float-left">{i.text}</div>

              <a href={`${i.link}`} className="text-ross2" target="_blank" rel="noreferrer">
                {i.link}
              </a>
            </li>
          );
        })}
      </ul>
    </div>
  );
};

export default About;

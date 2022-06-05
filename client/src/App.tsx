import React from "react";
import { Switch, Route } from "react-router-dom";
import About from "./pages/About";
import Contests from "./pages/Contests";
import Profile from "./pages/Profile";
import MicrosoftLogin from "./utils/requests/MicrosoftLogin";
import Header from "./components/Shared/Header";
import Contest from "./pages/Contest";
import User from "./models/User";
import CreateEditContest from "./components/Organizer/CreateEditContest";
import Admin from "./pages/Admin";
import NotFound from "./components/Shared/Errors/NotFound";
import JoinTeam from "./components/Contestant/JoinTeam";

function App() {
  const [user, setUser] = React.useState<User>(new User());
  React.useEffect(() => {
    (async () => {
      const u = await MicrosoftLogin.loginWithToken();
      setUser(u);
    })();
  }, []);
  return (
    <>
      <Header />

      <h1 className="text-[30px] font-Ropa">
        Due to a technical error registration has been moved to Google forms
        using{" "}
        <a href="https://bit.ly/jpc5reg" className="underline text-blue-600">
          this link
        </a>
      </h1>
    </>
  );
}

export default App;

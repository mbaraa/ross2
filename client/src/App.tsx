import React from "react";
import { Switch, Route, useHistory } from "react-router-dom";
import About from "./pages/About";
import Contests from "./pages/Contests";
import Profile from "./pages/Profile";
import MicrosoftLogin from "./utils/requests/MicrosoftLogin";
import Header from "./components/Shared/Header";
import Contest from "./pages/Contest";
import User from "./models/User";
import CreateEditContest from "./components/Organizer/CreateEditContest";
import Admin from "./pages/Admin";

function App() {
  const router = useHistory();
  const [user, setUser] = React.useState<User>(new User());

  return (
    <>
      <Header />

      <Switch>
        <Route exact path="/">
          <Contests user={user} />
        </Route>

        <Route path="/admin">
          <Admin user={user} />
        </Route>

        <Route path="/contest/:id">
          <Contest user={user} />
        </Route>

        <Route path="/contests/new">
          <CreateEditContest user={user} />
        </Route>

        <Route path="/profile">
          <Profile user={user} />
        </Route>

        <Route path="/about">
          <About />
        </Route>
      </Switch>
    </>
  );
}

export default App;

import React from "react";
import { Switch, Route } from "react-router-dom";
import About from "./pages/About";
import Contests from "./pages/Contests";
import Profile from "./pages/Profile";
import MicrosoftLogin from "./utils/requests/MicrosoftLogin";
import Header from "./components/Header";
import Contest from "./pages/Contest";
import User from "./models/User";
import CreateEditContest from "./components/CreateEditContest";

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

      <Switch>
        <Route exact path="/">
          <Contests user={user} />
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

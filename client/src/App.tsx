import React from "react";
import { Router, Switch, Route } from "react-router-dom";
import About from "./pages/About";
import Contests from "./pages/Contests";
import Profile from "./pages/Profile";
import MicrosoftLogin from "./utils/requests/MicrosoftLogin";
import Header from "./components/Header";

function App() {
  const [user, setUser] = React.useState<any>(0);

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

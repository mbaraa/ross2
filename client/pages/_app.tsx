import "../styles/globals.css";
import type { AppProps } from "next/app";
import Header from "../src/components/Header";
import * as React from "react";
import dynamic from "next/dynamic";
import User, { UserType } from "../src/models/User";
import MicrosoftLogin from "../src/utils/requests/MicrosoftLogin";

function MyApp({ Component, pageProps }: AppProps) {
  const [user, setUser] = React.useState<any>(0);

  React.useEffect(() => {
    (async () => {
      const u = await MicrosoftLogin.loginWithToken();
      setUser(u);
    })();
  }, []);

  console.log(user);

  return (
    <>
      <Header></Header>
      <div className="font-[Poppins] absolute left-[0.2rem] top-[4.2em] w-full p-[52px]">
        <Component user={user} {...pageProps} />
      </div>
    </>
  );
}

export default MyApp;

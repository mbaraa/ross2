import MicrosoftLogin from "react-microsoft-login";
import { default as MSLogin } from "../../utils/requests/MicrosoftLogin";
import config from "../../config";
import * as React from "react";
import Button from "@mui/material/Button";

const Login = (): React.ReactElement => {
  const [browser, setBrowser] = React.useState<boolean>(false);
  React.useEffect(() => {
    setBrowser(true);
  }, [process.browser]);

  const loginHandler = (error: any, authData: any, msalInstance: any) => {
    // console.log("data", error);
    // console.log("data", authData);
    // console.log("data", msalInstance);
    // if (!error) {
    //   console.error(error);
    //   return;
    // }

    (async () => {
      await MSLogin.login(authData);
    })();
  };

  return (
    <>
      {browser && (
        <MicrosoftLogin
          clientId={config.msalConfig.auth.clientId}
          authCallback={loginHandler}
          graphScopes={config.msalConfig.scopes}
          redirectUri={config.msalConfig.auth.redirectUri}
          buttonTheme="light_short"
          tenantUrl={config.msalConfig.auth.authority}
        />
      )}
    </>
  );
};

export default Login;
